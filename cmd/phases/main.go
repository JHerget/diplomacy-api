package main

import (
	"context"
	h "diplomacy-api/internal/http"
	"diplomacy-api/internal/phases"
	"diplomacy-api/internal/platform/mongo"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(phaseHandler *phases.Handler) func(context.Context, events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return func(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		switch event.RouteKey {
		case "GET /phases":
			return phaseHandler.GetAll(ctx, event)
		default:
			return h.BadRequest(&h.Error{
				Message: "method not allowed",
			}), nil
		}
	}
}

func main() {
	db, err := mongo.NewMongoDB(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	phaseHandler := phases.NewHandler(phases.NewRepository(db))
	lambda.Start(handler(phaseHandler))
}
