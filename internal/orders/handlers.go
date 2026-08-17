package orders

import (
	"context"
	"diplomacy-api/internal/game"
	h "diplomacy-api/internal/http"
	"diplomacy-api/internal/models"
	"diplomacy-api/internal/platform/mongo"
	"diplomacy-api/internal/utils"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAll(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
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

	return h.OK(turn.Orders), nil
}

func GetByID(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	gameID := event.PathParameters["gid"]
	turnID := event.PathParameters["tid"]
	orderID := event.PathParameters["oid"]

	g, _, err := getGameAndRepo(ctx, gameID)
	if err != nil {
		return h.InternalServerError(err), err
	}

	turn, ok := findTurn(g, turnID)
	if !ok {
		return invalidTurnID(turnID)
	}

	order, ok := findOrder(turn, orderID)
	if !ok {
		return invalidOrderID(orderID)
	}

	return h.OK(order), nil
}

func Create(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	gameID := event.PathParameters["gid"]
	turnID := event.PathParameters["tid"]

	var order orderRequest
	err := json.Unmarshal([]byte(event.Body), &order)
	if err != nil {
		return h.BadRequest(&h.Error{
			Message: "unable to parse request body",
		}), err
	}

	g, repo, err := getGameAndRepo(ctx, gameID)
	if err != nil {
		return h.InternalServerError(err), err
	}

	turn, ok := findTurn(g, turnID)
	if !ok {
		return invalidTurnID(turnID)
	}

	savedOrder := models.Order{
		ID:          primitive.NewObjectID().Hex(),
		PhaseID:     order.PhaseID,
		PlayerName:  order.PlayerName,
		CreatedDate: 0,
		Value:       order.Value,
	}

	statusCreated := upsertOrder(turn, &savedOrder)

	if err := g.Valid(); err != nil {
		return h.BadRequest(&h.Error{
			Message: err.Error(),
		}), err
	}

	if err := repo.Update(ctx, g); err != nil {
		return h.InternalServerError(err), err
	}

	if statusCreated {
		return h.Created(savedOrder), nil
	}
	return h.OK(savedOrder), nil
}

type orderRequest struct {
	PlayerName string `json:"playerName"`
	PhaseID    string `json:"phaseID"`
	Value      string `json:"value"`
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

func findOrder(t *models.Turn, orderID string) (*models.Order, bool) {
	return utils.Find(t.Orders, func(o *models.Order) bool {
		return o.ID == orderID
	})
}

func findPlayerOrder(t *models.Turn, playerName string) (*models.Order, bool) {
	return utils.Find(t.Orders, func(o *models.Order) bool {
		return o.PlayerName == playerName
	})
}

func upsertOrder(t *models.Turn, order *models.Order) bool {
	existing, ok := findPlayerOrder(t, order.PlayerName)
	if !ok {
		t.Orders = append(t.Orders, *order)
		return true
	}

	order.ID = existing.ID
	order.CreatedDate = existing.CreatedDate
	*existing = *order
	return false
}

func invalidTurnID(turnID string) (events.APIGatewayV2HTTPResponse, error) {
	err := fmt.Errorf("invalid turn id '%s'", turnID)
	return h.BadRequest(&h.Error{
		Message: err.Error(),
	}), err
}

func invalidOrderID(orderID string) (events.APIGatewayV2HTTPResponse, error) {
	err := fmt.Errorf("invalid order id '%s'", orderID)
	return h.BadRequest(&h.Error{
		Message: err.Error(),
	}), err
}
