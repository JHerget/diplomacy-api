package main

import (
	"context"
	h "diplomacy-api/internal/http"
	"diplomacy-api/internal/turns"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	method := event.RequestContext.HTTP.Method
	turnID := event.PathParameters["tid"]

	switch {
	case method == http.MethodGet && turnID == "":
		return turns.GetAll(ctx, event)
	case method == http.MethodGet && turnID != "":
		return turns.GetByID(ctx, event)
	case method == http.MethodPost && turnID == "":
		return turns.Create(ctx, event)
	case method == http.MethodPut && turnID != "":
		return turns.Update(ctx, event)
	case method == http.MethodDelete && turnID != "":
		return turns.Delete(ctx, event)
	default:
		return h.BadRequest(&h.Error{
			Message: "method not allowed",
		}), nil
	}
}

func main() {
	lambda.Start(handler)
}
