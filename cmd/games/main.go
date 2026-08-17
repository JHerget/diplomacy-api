package main

import (
	"context"
	"diplomacy-api/internal/game"
	h "diplomacy-api/internal/http"
	"diplomacy-api/internal/maps"
	"diplomacy-api/internal/platform/mongo"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(gameHandler *game.Handler) func(context.Context, events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return func(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		switch event.RouteKey {
		case "GET /games":
			return gameHandler.GetAll(ctx, event)
		case "GET /games/{gid}":
			return gameHandler.GetByID(ctx, event)
		case "POST /games":
			return gameHandler.Create(ctx, event)
		case "PUT /games/{gid}":
			return gameHandler.Update(ctx, event)
		case "DELETE /game/{gid}":
			return gameHandler.Delete(ctx, event)
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

	gameHandler := game.NewHandler(game.NewRepository(db), maps.NewRepository(db))
	lambda.Start(handler(gameHandler))
}
