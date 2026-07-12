package maps

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
		}), err
	}

	mapRepo := NewRepository(db)
	summaries, err := mapRepo.GetAllSummaries(ctx)
	if err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), err
	}

	return h.OK(summaries), nil
}

func GetByID(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	mid := event.PathParameters["mid"]

	db, err := mongo.NewMongoDB(ctx)
	if err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), err
	}

	mapRepo := NewRepository(db)
	m, err := mapRepo.Get(ctx, mid)
	if err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), err
	}

	return h.OK(m), nil
}

func Create(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       "create map",
	}, nil
}

func Update(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       "update map",
	}, nil
}

func Delete(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       "delete map",
	}, nil
}
