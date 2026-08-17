package main

import (
	"context"
	"diplomacy-api/internal/board"
	"diplomacy-api/internal/game"
	h "diplomacy-api/internal/http"
	"diplomacy-api/internal/platform/aws"
	"diplomacy-api/internal/platform/mongo"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(boardHandler *board.Handler) func(context.Context, events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return func(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		switch event.RouteKey {
		case "GET /games/{gid}/board":
			return boardHandler.Get(ctx, event)
		default:
			return h.BadRequest(&h.Error{
				Message: "method not allowed",
			}), nil
		}
	}
}

func main() {
	ctx := context.Background()
	db, err := mongo.NewMongoDB(ctx)
	if err != nil {
		log.Fatal(err)
	}
	s3, err := aws.NewS3(ctx)
	if err != nil {
		log.Fatal(err)
	}

	boardHandler := board.NewHandler(game.NewRepository(db), s3)
	lambda.Start(handler(boardHandler))
}
