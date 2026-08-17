package main

import (
	"context"
	"diplomacy-api/internal/game"
	h "diplomacy-api/internal/http"
	"diplomacy-api/internal/phases"
	"diplomacy-api/internal/platform/mongo"
	"diplomacy-api/internal/turns"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(turnHandler *turns.Handler) func(context.Context, events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return func(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		switch event.RouteKey {
		case "GET /games/{gid}/turns":
			return turnHandler.GetAll(ctx, event)
		case "GET /games/{gid}/turns/{tid}":
			return turnHandler.GetByID(ctx, event)
		case "POST /games/{gid}/turns":
			return turnHandler.Create(ctx, event)
		case "PUT /games/{gid}/turns/{tid}":
			return turnHandler.Update(ctx, event)
		case "DELETE /games/{gid}/turns/{tid}":
			return turnHandler.Delete(ctx, event)
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

	turnHandler := turns.NewHandler(game.NewRepository(db), phases.NewRepository(db))
	lambda.Start(handler(turnHandler))
}
