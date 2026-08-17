package turns

import (
	"context"
	"diplomacy-api/internal/game"
	h "diplomacy-api/internal/http"
	"diplomacy-api/internal/models"
	"diplomacy-api/internal/phases"
	"diplomacy-api/internal/platform/mongo"
	"diplomacy-api/internal/utils"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

func GetAll(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	gameID := event.PathParameters["gid"]

	g, _, err := getGameAndRepo(ctx, gameID)
	if err != nil {
		return h.InternalServerError(err), err
	}

	return h.OK(g.Turns), nil
}

func GetByID(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	gameID := event.PathParameters["gid"]
	turnID := event.PathParameters["tid"]

	g, _, err := getGameAndRepo(ctx, gameID)
	if err != nil {
		return h.InternalServerError(err), err
	}

	turn, ok := findTurn(g, turnID)
	if !ok {
		return invalidTurnID(turnID)
	}

	return h.OK(turn), nil
}

func Create(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	gameID := event.PathParameters["gid"]

	db, err := mongo.NewMongoDB(ctx)
	if err != nil {
		return h.InternalServerError(err), err
	}

	gameRepo := game.NewRepository(db)
	g, err := gameRepo.Get(ctx, gameID)
	if err != nil {
		return h.InternalServerError(err), err
	}

	phaseRepo := phases.NewRepository(db)
	phase, err := phaseRepo.GetByOrder(ctx, 0)
	if err != nil {
		return h.InternalServerError(err), err
	}

	id, err := utils.RandomID()
	if err != nil {
		return h.InternalServerError(err), err
	}

	startDate := nextTurnStartDate(g)
	endDate := int(time.Unix(int64(startDate), 0).AddDate(0, 0, g.DaysPerTurn).Unix())
	turn := models.Turn{
		ID:         id,
		PhaseID:    phase.ID,
		Orders:     []models.Order{},
		TurnNumber: len(g.Turns) + 1,
		StartDate:  startDate,
		EndDate:    endDate,
	}
	g.Turns = append(g.Turns, turn)

	if err := g.Valid(); err != nil {
		return h.BadRequest(&h.Error{
			Message: err.Error(),
		}), err
	}

	if err := gameRepo.Update(ctx, g); err != nil {
		return h.InternalServerError(err), err
	}

	return h.Created(turn), nil
}

func Update(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	gameID := event.PathParameters["gid"]
	turnID := event.PathParameters["tid"]

	var turn models.Turn
	if err := json.Unmarshal([]byte(event.Body), &turn); err != nil {
		return h.BadRequest(&h.Error{
			Message: err.Error(),
		}), err
	}
	turn.ID = turnID

	g, repo, err := getGameAndRepo(ctx, gameID)
	if err != nil {
		return h.InternalServerError(err), err
	}

	existing, ok := findTurn(g, turnID)
	if !ok {
		return invalidTurnID(turnID)
	}
	*existing = turn

	if err := g.Valid(); err != nil {
		return h.BadRequest(&h.Error{
			Message: err.Error(),
		}), err
	}

	if err := repo.Update(ctx, g); err != nil {
		return h.InternalServerError(err), err
	}

	return h.OK(turn), nil
}

func Delete(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	gameID := event.PathParameters["gid"]
	turnID := event.PathParameters["tid"]

	g, repo, err := getGameAndRepo(ctx, gameID)
	if err != nil {
		return h.InternalServerError(err), err
	}

	index := -1
	for i := range g.Turns {
		if g.Turns[i].ID == turnID {
			index = i
			break
		}
	}
	if index == -1 {
		return invalidTurnID(turnID)
	}
	g.Turns = append(g.Turns[:index], g.Turns[index+1:]...)

	if err := g.Valid(); err != nil {
		return h.BadRequest(&h.Error{
			Message: err.Error(),
		}), err
	}

	if err := repo.Update(ctx, g); err != nil {
		return h.InternalServerError(err), err
	}

	return h.NoContent(), nil
}

func getGameAndRepo(ctx context.Context, gameID string) (*models.Game, *game.Repository, error) {
	db, err := mongo.NewMongoDB(ctx)
	if err != nil {
		return nil, nil, err
	}

	repo := game.NewRepository(db)
	g, err := repo.Get(ctx, gameID)
	if err != nil {
		return nil, nil, err
	}

	return g, repo, nil
}

func findTurn(g *models.Game, turnID string) (*models.Turn, bool) {
	return utils.Find(g.Turns, func(t *models.Turn) bool {
		return t.ID == turnID
	})
}

func invalidTurnID(turnID string) (events.APIGatewayV2HTTPResponse, error) {
	err := fmt.Errorf("invalid turn id '%s'", turnID)
	return h.BadRequest(&h.Error{
		Message: err.Error(),
	}), err
}

func nextTurnStartDate(g *models.Game) int {
	sourceDate := g.StartDate
	if len(g.Turns) > 0 {
		sourceDate = g.Turns[len(g.Turns)-1].EndDate
	}

	location := time.FixedZone("game", g.Timezone*60*60)
	t := time.Unix(int64(sourceDate), 0).In(location)
	start := time.Date(t.Year(), t.Month(), t.Day(), g.TurnStartHour, 0, 0, 0, location)

	return int(start.Unix())
}
