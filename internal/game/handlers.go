package game

import (
	"context"
	h "diplomacy-api/internal/http"
	"diplomacy-api/internal/platform/mongo"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func GetAll(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	db, err := mongo.NewMongoDB(ctx)
	if err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), nil
	}

	gameRepo := NewRepository(db)
	allGames, err := gameRepo.GetAll(ctx)
	if err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), nil
	}

	return h.OK(allGames), nil
}

func GetByID(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	gameID := event.PathParameters["gid"]

	db, err := mongo.NewMongoDB(ctx)
	if err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), nil
	}

	gameRepo := NewRepository(db)
	g, err := gameRepo.Get(ctx, gameID)
	if err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), nil
	}

	return h.OK(g), nil
}

func Create(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       "create game",
	}, nil
}

func Update(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       "update game",
	}, nil
}

func Delete(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       "delete game",
	}, nil
}
