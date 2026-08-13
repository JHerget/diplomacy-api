package game

import (
	"context"
	h "diplomacy-api/internal/http"
	"diplomacy-api/internal/maps"
	"diplomacy-api/internal/models"
	"diplomacy-api/internal/platform/mongo"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

func GetAll(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	db, err := mongo.NewMongoDB(ctx)
	if err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), err
	}

	gameRepo := NewRepository(db)
	allGames, err := gameRepo.GetAll(ctx)
	if err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), err
	}

	return h.OK(allGames), nil
}

func GetByID(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	gameID := event.PathParameters["gid"]

	db, err := mongo.NewMongoDB(ctx)
	if err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), err
	}

	gameRepo := NewRepository(db)
	g, err := gameRepo.Get(ctx, gameID)
	if err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), err
	}

	return h.OK(g), nil
}

func Create(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	var req createGameRequest
	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), err
	}

	db, err := mongo.NewMongoDB(ctx)
	if err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), err
	}

	gameRepo := NewRepository(db)
	mapRepo := maps.NewRepository(db)

	m, err := mapRepo.Get(ctx, req.MapID)
	if err != nil {
		return h.BadRequest(&h.Error{
			Message: fmt.Sprintf("map with id '%s' not found: %s", req.MapID, err),
		}), err
	}

	g := models.Game{
		OwnerID:       "1",
		Map:           m.Summary(),
		Board:         m.Providences,
		Players:       m.Players,
		DaysPerTurn:   req.DaysPerTurn,
		TurnStartHour: req.TurnStartHour,
		StartDate:     req.StartDate,
	}

	if err := g.Valid(); err != nil {
		return h.BadRequest(&h.Error{
			Message: err.Error(),
		}), err
	}

	gameRepo.Create(ctx, &g)

	return h.Created(&g), nil
}

func Update(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	var g models.Game
	if err := json.Unmarshal([]byte(event.Body), &g); err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), err
	}

	if err := g.Valid(); err != nil {
		return h.BadRequest(&h.Error{
			Message: err.Error(),
		}), err
	}

	db, err := mongo.NewMongoDB(ctx)
	if err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), err
	}

	gameRepo := NewRepository(db)
	if err := gameRepo.Update(ctx, &g); err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), err
	}

	return h.OK(g), nil
}

func Delete(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	gameID := event.PathParameters["gid"]

	db, err := mongo.NewMongoDB(ctx)
	if err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), err
	}

	gameRepo := NewRepository(db)
	if err := gameRepo.Delete(ctx, gameID); err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), err
	}

	return h.NoContent(), nil
}

type createGameRequest struct {
	MapID         string `json:"mapId"`
	DaysPerTurn   int    `json:"daysPerTurn"`
	TurnStartHour int    `json:"turnStartHour"`
	Timezone      int    `json:"timezone"`
	StartDate     int    `json:"startDate"`
}
