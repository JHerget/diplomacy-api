package tests

import (
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"

	"diplomacy-api/internal/game"
	"diplomacy-api/internal/models"
	"diplomacy-api/internal/turns"

	"github.com/aws/aws-lambda-go/events"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

// Player assignments live on the game. Close the turn through the real handler
// and inspect the board it saves. Mongo's in-memory mock supplies only the
// repository read and write acknowledgments.
func checkGameVariant(t *testing.T, id, number int) {
	t.Helper()
	f, players := loadFixture(t, id)
	t.Log(f.ExpectedOutcome)
	g := fixtureGame(f, players, number)
	switch id {
	case 808, 809:
		g.PlayerCount = 6
	case 810:
		g.PlayerCount = 5
	case 811:
		g.PlayerCount = 4
	case 812, 813:
		g.PlayerCount = 2
	}
	for i := range g.Players {
		g.Players[i].IsPlaying = false
	}
	for i, combo := range models.PlayerCombos[g.PlayerCount] {
		powers := strings.Split(combo, "/")
		if id == 813 && slices.Contains(powers, "France") {
			// The fixed Fall 1914 coin result assigns Italy to the Allied player.
			powers = append(powers, "Italy")
		}
		if _, err := g.AssignSpecificPlayers(fmt.Sprintf("player-%d", i+1), powers); err != nil {
			t.Fatal(err)
		}
	}

	turn := g.CurrentTurn()
	turn.Orders = slices.Clone(f.Orders)
	disorder := ""
	switch id {
	case 801, 802, 803:
		disorder = "Italy"
	case 804, 805, 806:
		disorder = "France"
	case 807:
		disorder = "England"
	}
	for _, p := range g.Players {
		if !p.IsPlaying || p.Name == disorder {
			continue
		}
		if _, ok := turn.FindPlayerOrder(p.Name); !ok {
			turn.Orders = append(turn.Orders, models.Order{PlayerName: p.Name, Value: ""})
		}
	}
	// An active power can deliberately submit a blank sheet. A power in civil
	// disorder supplies no sheet, so the deadline is what closes its turn.
	turn.EndDate = int(time.Now().Add(time.Hour).Unix())
	if turn.IsFinished(g.Players) != (disorder == "") {
		t.Fatal("player assignments and submitted order sheets do not match the scenario")
	}
	turn.EndDate = 1
	closeGameTurn(t, &g)
	checkTurnBoard(t, id, g.Board, f)
}

func fixtureGame(f fixture, players []models.Player, number int) models.Game {
	g := models.Game{
		ID: "0123456789abcdef01234567", OwnerID: "owner", Map: models.MapSummary{ID: "original"},
		InProgress: true, DaysPerTurn: 1, PlayerCount: 7,
		Players: players, Board: f.InputBoard,
	}
	// Fill the turn history so CurrentTurn selects the intended season. The
	// two-player game starts in Spring 1914, the other variants in Spring 1901.
	for n := 1; n <= number; n++ {
		g.Turns = append(g.Turns, models.Turn{TurnNumber: n})
	}
	return g
}

func closeGameTurn(t *testing.T, g *models.Game) {
	t.Helper()
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock).ShareClient(true))
	doc, err := bson.Marshal(g)
	if err != nil {
		t.Fatal(err)
	}
	var stored bson.D
	if err := bson.Unmarshal(doc, &stored); err != nil {
		t.Fatal(err)
	}
	mt.AddMockResponses(mtest.CreateCursorResponse(0, "diplomacy.games", mtest.FirstBatch, stored),
		mtest.CreateSuccessResponse(bson.E{Key: "n", Value: 1}, bson.E{Key: "nModified", Value: 1}))
	handler := turns.NewHandler(game.NewRepository(mt.Client), nil)
	resp, err := handler.Create(t.Context(), events.APIGatewayV2HTTPRequest{PathParameters: map[string]string{"gid": g.ID}})
	if err != nil || resp.StatusCode != 201 {
		t.Fatalf("close turn: HTTP %d, %v (%s)", resp.StatusCode, err, resp.Body)
	}
	for _, ev := range mt.GetAllStartedEvents() {
		if ev.CommandName == "update" {
			// Apply the game fields the repository actually persisted.
			value, err := ev.Command.LookupErr("updates", "0", "u", "$set")
			if err != nil {
				t.Fatal(err)
			}
			saved := models.Game{ID: g.ID}
			if err := value.Unmarshal(&saved); err != nil {
				t.Fatal(err)
			}
			*g = saved
			return
		}
	}
	t.Fatal("turn handler did not save an updated game")
}
