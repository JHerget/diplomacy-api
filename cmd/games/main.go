package main

import (
	"context"
	"diplomacy-api/internal/game"
	h "diplomacy-api/internal/http"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	method := event.RequestContext.HTTP.Method
	gameID := event.PathParameters["gid"]

	switch {
	case method == http.MethodGet && gameID == "":
		return game.GetAll(ctx, event)
	case method == http.MethodGet && gameID != "":
		return game.GetByID(ctx, event)
	case method == http.MethodPost && gameID == "":
		return game.Create(ctx, event)
	case method == http.MethodPut && gameID != "":
		return game.Update(ctx, event)
	case method == http.MethodDelete && gameID != "":
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
