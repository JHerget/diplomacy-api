package turns

import (
	"context"
	"diplomacy-api/internal/game"
	"diplomacy-api/internal/http"
	"diplomacy-api/internal/models"
	"diplomacy-api/internal/phases"
	"diplomacy-api/internal/utils"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

type Handler struct {
	gameRepo  *game.Repository
	phaseRepo *phases.Repository
}

func NewHandler(gameRepo *game.Repository, phaseRepo *phases.Repository) *Handler {
	return &Handler{
		gameRepo:  gameRepo,
		phaseRepo: phaseRepo,
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

	phase, err := h.phaseRepo.GetByOrder(ctx, 0)
	if err != nil {
		return http.InternalServerError(err), err
	}

	id, err := utils.RandomID()
	if err != nil {
		return http.InternalServerError(err), err
	}

	startDate := g.NextTurnStartDate()
	endDate := int(time.Unix(int64(startDate), 0).UTC().AddDate(0, 0, g.DaysPerTurn).Unix())
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
		return http.BadRequest(&http.Error{
			Message: err.Error(),
		}), err
	}

	if err := h.gameRepo.Update(ctx, g); err != nil {
		return http.InternalServerError(err), err
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
