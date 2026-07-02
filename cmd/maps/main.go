package main

import (
	"context"
	h "diplomacy-api/internal/http"
	"diplomacy-api/internal/maps"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	method := event.RequestContext.HTTP.Method
	mapID := event.PathParameters["mid"]

	switch {
	case method == http.MethodGet && mapID == "":
		return maps.GetAll(ctx, event)
	case method == http.MethodGet && mapID != "":
		return maps.GetByID(ctx, event)
	case method == http.MethodPost && mapID == "":
		return maps.Create(ctx, event)
	case method == http.MethodPut && mapID != "":
		return maps.Update(ctx, event)
	case method == http.MethodDelete && mapID != "":
		return maps.Delete(ctx, event)
	default:
		return h.BadRequest(&h.Error{
			Message: "method not allowed",
		}), nil
	}
}

func main() {
	lambda.Start(handler)
}
