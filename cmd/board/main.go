package main

import (
	"context"
	"diplomacy-api/internal/game"
	"diplomacy-api/internal/platform/mongo"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	db, err := mongo.NewMongoDB(ctx)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
		}, err
	}

	gameRepo := game.MakeRepository(db)

	game, err := gameRepo.Get(ctx, "test")
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
		}, err
	}

	body, err := json.Marshal(game)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
		}, err
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Body:       string(body),
	}, nil
}

func main() {
	lambda.Start(handler)
}
