package main

import (
	"context"
	h "diplomacy-api/internal/http"
	"diplomacy-api/internal/players"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	method := event.RequestContext.HTTP.Method
	playerID := event.PathParameters["pid"]

	switch {
	case method == http.MethodGet && playerID == "":
		return players.GetAll(ctx, event)
	case method == http.MethodGet && playerID != "":
		return players.GetByID(ctx, event)
	case method == http.MethodPost && playerID == "":
		return players.Create(ctx, event)
	case method == http.MethodPut && playerID != "":
		return players.Update(ctx, event)
	case method == http.MethodDelete && playerID != "":
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
