package game

import (
	"context"
	"diplomacy-api/internal/http"
	"diplomacy-api/internal/maps"
	"diplomacy-api/internal/models"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
)

type Handler struct {
	gameRepo *Repository
	mapRepo  *maps.Repository
}

func NewHandler(gameRepo *Repository, mapRepo *maps.Repository) *Handler {
	return &Handler{
		gameRepo: gameRepo,
		mapRepo:  mapRepo,
	}
}

func (h *Handler) GetAll(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	filter := GameFilter{}
	if externalID, ok := event.QueryStringParameters["external_id"]; ok {
		filter.ExternalID = &externalID
	}
	if inProgress, ok := event.QueryStringParameters["in_progress"]; ok {
		value, err := strconv.ParseBool(inProgress)
		if err != nil {
			err := fmt.Errorf("in_progress must be convertible to bool")
			return http.BadRequest(&http.Error{
				Message: err.Error(),
			}), err
		}
		filter.InProgress = &value
	}

	allGames, err := h.gameRepo.GetAll(ctx, filter)
	if err != nil {
		return http.InternalServerError(err), err
	}

	return http.OK(allGames), nil
}

func (h *Handler) GetByID(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	gameID := event.PathParameters["gid"]

	g, err := h.gameRepo.Get(ctx, gameID)
	if err != nil {
		return http.InternalServerError(err), err
	}

	return http.OK(g), nil
}

func (h *Handler) Create(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	var req createGameRequest
	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		return http.InternalServerError(err), err
	}

	m, err := h.mapRepo.Get(ctx, req.MapID)
	if err != nil {
		return http.BadRequest(&http.Error{
			Message: fmt.Sprintf("map with id '%s' not found: %s", req.MapID, err),
		}), err
	}

	g := models.Game{
		OwnerID:       "1",
		ExternalID:    req.ExternalID,
		Map:           m.Summary(),
		Board:         m.Providences,
		Players:       m.Players,
		DaysPerTurn:   req.DaysPerTurn,
		Timezone:      req.Timezone,
		TurnStartHour: req.TurnStartHour,
		StartDate:     req.StartDate,
		InProgress:    true,
	}

	if err := g.Valid(); err != nil {
		return http.BadRequest(&http.Error{
			Message: err.Error(),
		}), err
	}

	h.gameRepo.Create(ctx, &g)

	return http.Created(&g), nil
}

func (h *Handler) Update(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	var g models.Game
	if err := json.Unmarshal([]byte(event.Body), &g); err != nil {
		return http.InternalServerError(err), err
	}

	if err := g.Valid(); err != nil {
		return http.BadRequest(&http.Error{
			Message: err.Error(),
		}), err
	}

	if err := h.gameRepo.Update(ctx, &g); err != nil {
		return http.InternalServerError(err), err
	}

	return http.OK(g), nil
}

func (h *Handler) Delete(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	gameID := event.PathParameters["gid"]

	if err := h.gameRepo.Delete(ctx, gameID); err != nil {
		return http.InternalServerError(err), err
	}

	return http.NoContent(), nil
}

type createGameRequest struct {
	ExternalID    *string `json:"externalId"`
	MapID         string  `json:"mapId"`
	DaysPerTurn   int     `json:"daysPerTurn"`
	TurnStartHour int     `json:"turnStartHour"`
	Timezone      int     `json:"timezone"`
	StartDate     int     `json:"startDate"`
}
