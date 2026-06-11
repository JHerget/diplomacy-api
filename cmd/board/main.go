package main

import (
	"context"
	"diplomacy-api/internal/game"
	"diplomacy-api/internal/http"
	"diplomacy-api/internal/platform/mongo"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	db, err := mongo.NewMongoDB(ctx)
	if err != nil {
		return http.InternalServerError(&http.Error{
			Message: err.Error(),
		}), err
	}

	gameRepo := game.MakeRepository(db)

	game, err := gameRepo.Get(ctx, "69564bb933c5739468982b67")
	if err != nil {
		return http.InternalServerError(&http.Error{
			Message: err.Error(),
		}), err
	}

	return http.Ok(game), nil
}

func main() {
	lambda.Start(handler)
}
