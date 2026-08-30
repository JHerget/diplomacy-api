package orders

import (
	"context"
	"diplomacy-api/internal/game"
	"diplomacy-api/internal/http"
	"diplomacy-api/internal/models"
	"diplomacy-api/internal/utils"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

type Handler struct {
	gameRepo *game.Repository
}

func NewHandler(gameRepo *game.Repository) *Handler {
	return &Handler{
		gameRepo: gameRepo,
	}
}

func (h *Handler) GetAll(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	gameID := event.PathParameters["gid"]
	turnID := event.PathParameters["tid"]

	g, err := h.gameRepo.Get(ctx, gameID)
	if err != nil {
		return http.InternalServerError(err), err
	}

	turn, ok := g.FindTurn(turnID)
	if !ok {
		return invalidTurnID(turnID)
	}

	return http.OK(turn.Orders), nil
}

func (h *Handler) GetByID(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	gameID := event.PathParameters["gid"]
	turnID := event.PathParameters["tid"]
	orderID := event.PathParameters["oid"]

	g, err := h.gameRepo.Get(ctx, gameID)
	if err != nil {
		return http.InternalServerError(err), err
	}

	turn, ok := g.FindTurn(turnID)
	if !ok {
		return invalidTurnID(turnID)
	}

	order, ok := turn.FindOrder(orderID)
	if !ok {
		return invalidOrderID(orderID)
	}

	return http.OK(order), nil
}

func (h *Handler) Create(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	gameID := event.PathParameters["gid"]
	turnID := event.PathParameters["tid"]

	var order orderRequest
	err := json.Unmarshal([]byte(event.Body), &order)
	if err != nil {
		return http.BadRequest(&http.Error{
			Message: err.Error(),
		}), err
	}

	g, err := h.gameRepo.Get(ctx, gameID)
	if err != nil {
		return http.InternalServerError(err), err
	}

	turn, ok := g.FindTurn(turnID)
	if !ok {
		return invalidTurnID(turnID)
	}

	for i := range g.Players {
		if g.Players[i].Name == order.PlayerName && !g.Players[i].IsPlaying {
			err := fmt.Errorf("Player is not playing anymore.")
			return http.BadRequest(&http.Error{
				Message: err.Error(),
			}), err
		}
	}

	id, err := utils.RandomID()
	if err != nil {
		return http.InternalServerError(err), err
	}

	savedOrder := models.Order{
		ID:          id,
		PhaseID:     order.PhaseID,
		PlayerName:  order.PlayerName,
		CreatedDate: int(time.Now().UTC().Unix()),
		Value:       order.Value,
	}

	statusCreated := false
	existing, ok := turn.FindPlayerOrder(savedOrder.PlayerName)
	if !ok {
		turn.Orders = append(turn.Orders, savedOrder)
		statusCreated = true
	} else {
		savedOrder.ID = existing.ID
		savedOrder.CreatedDate = existing.CreatedDate
		*existing = savedOrder
	}

	if err := g.Valid(); err != nil {
		return http.BadRequest(&http.Error{
			Message: err.Error(),
		}), err
	}

	if err := h.gameRepo.Update(ctx, g); err != nil {
		return http.InternalServerError(err), err
	}

	if statusCreated {
		return http.Created(savedOrder), nil
	}

	return http.OK(savedOrder), nil
}

type orderRequest struct {
	PlayerName string `json:"playerName"`
	PhaseID    string `json:"phaseID"`
	Value      string `json:"value"`
}

func invalidTurnID(turnID string) (events.APIGatewayV2HTTPResponse, error) {
	err := fmt.Errorf("invalid turn id '%s'", turnID)
	return http.BadRequest(&http.Error{
		Message: err.Error(),
	}), err
}

func invalidOrderID(orderID string) (events.APIGatewayV2HTTPResponse, error) {
	err := fmt.Errorf("invalid order id '%s'", orderID)
	return http.BadRequest(&http.Error{
		Message: err.Error(),
	}), err
}
