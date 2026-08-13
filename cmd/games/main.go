package main

import (
	"context"
	"diplomacy-api/internal/game"
	h "diplomacy-api/internal/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	switch event.RouteKey {
	case "GET /games":
		return game.GetAll(ctx, event)
	case "GET /games/{gid}":
		return game.GetByID(ctx, event)
	case "POST /games":
		return game.Create(ctx, event)
	case "PUT /games/{gid}":
		return game.Update(ctx, event)
	case "DELETE /game/{gid}":
		return game.Delete(ctx, event)
	default:
		return h.BadRequest(&h.Error{
			Message: "method not allowed",
		}), nil
	}
}

func main() {
	lambda.Start(handler)
}
