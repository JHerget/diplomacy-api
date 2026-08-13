package main

import (
	"context"
	h "diplomacy-api/internal/http"
	"diplomacy-api/internal/players"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	switch event.RouteKey {
	case "GET /games/{gid}/players":
		return players.GetAll(ctx, event)
	case "GET /games/{gid}/players/{pid}":
		return players.GetByID(ctx, event)
	case "POST /games/{gid}/players":
		return players.Create(ctx, event)
	case "PUT /games/{gid}/players/{pid}":
		return players.Update(ctx, event)
	case "DELETE /games/{gid}/players/{pid}":
		return players.Delete(ctx, event)
	default:
		return h.BadRequest(&h.Error{
			Message: "method not allowed",
		}), nil
	}
}

func main() {
	lambda.Start(handler)
}
