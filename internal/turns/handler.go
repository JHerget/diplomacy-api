package turns

import (
	"context"
	"diplomacy-api/internal/game"
	"diplomacy-api/internal/http"
	"diplomacy-api/internal/models"
	"diplomacy-api/internal/phases"
	"diplomacy-api/internal/platform/aws"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

type Handler struct {
	gameRepo  *game.Repository
	phaseRepo *phases.Repository
	notifier  *aws.SQS
}

func NewHandler(gameRepo *game.Repository, phaseRepo *phases.Repository, notifier *aws.SQS) *Handler {
	return &Handler{
		gameRepo:  gameRepo,
		phaseRepo: phaseRepo,
		notifier:  notifier,
	}
}

func (h *Handler) GetAll(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	gameID := event.PathParameters["gid"]

	g, err := h.gameRepo.Get(ctx, gameID)
	if err != nil {
		return http.InternalServerError(err), err
	}

	return http.OK(g.Turns), nil
}

func (h *Handler) GetByID(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
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

	return http.OK(turn), nil
}

func (h *Handler) Create(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	gameID := event.PathParameters["gid"]

	g, err := h.gameRepo.Get(ctx, gameID)
	if err != nil {
		return http.InternalServerError(err), err
	}

	currentTurn := g.CurrentTurn()
	if currentTurn != nil && !currentTurn.IsFinished(g.Players) {
		err := fmt.Errorf("current turn is not finished")
		return http.BadRequest(&http.Error{
			Message: err.Error(),
		}), err
	}

	turn, err := g.NewTurn()
	if err != nil {
		return http.InternalServerError(err), err
	}

	if err := h.gameRepo.Update(ctx, g); err != nil {
		return http.InternalServerError(err), err
	}

	if g.ExternalID != nil && *g.ExternalID != "" {
		turnID := turn.ID
		if currentTurn != nil {
			turnID = currentTurn.ID
		}

		err := h.notifier.Send(ctx, models.NotificationMessage{
			ChannelID: *g.ExternalID,
			GameID:    gameID,
			TurnID:    turnID,
		})
		if err != nil {
			return http.InternalServerError(err), err
		}
	}

	return http.Created(turn), nil
}

func (h *Handler) Update(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	gameID := event.PathParameters["gid"]
	turnID := event.PathParameters["tid"]

	var turn models.Turn
	if err := json.Unmarshal([]byte(event.Body), &turn); err != nil {
		return http.BadRequest(&http.Error{
			Message: err.Error(),
		}), err
	}
	turn.ID = turnID

	g, err := h.gameRepo.Get(ctx, gameID)
	if err != nil {
		return http.InternalServerError(err), err
	}

	existing, ok := g.FindTurn(turnID)
	if !ok {
		return invalidTurnID(turnID)
	}
	*existing = turn

	if err := g.Valid(); err != nil {
		return http.BadRequest(&http.Error{
			Message: err.Error(),
		}), err
	}

	if err := h.gameRepo.Update(ctx, g); err != nil {
		return http.InternalServerError(err), err
	}

	return http.OK(turn), nil
}

func (h *Handler) Delete(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	gameID := event.PathParameters["gid"]
	turnID := event.PathParameters["tid"]

	g, err := h.gameRepo.Get(ctx, gameID)
	if err != nil {
		return http.InternalServerError(err), err
	}

	if err := g.RemoveTurn(turnID); err != nil {
		return http.BadRequest(&http.Error{
			Message: err.Error(),
		}), err
	}

	if err := h.gameRepo.Update(ctx, g); err != nil {
		return http.InternalServerError(err), err
	}

	return http.NoContent(), nil
}

func invalidTurnID(turnID string) (events.APIGatewayV2HTTPResponse, error) {
	err := fmt.Errorf("invalid turn id '%s'", turnID)
	return http.BadRequest(&http.Error{
		Message: err.Error(),
	}), err
}
