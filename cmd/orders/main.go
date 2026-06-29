package main

import (
	"context"
	h "diplomacy-api/internal/http"
	"diplomacy-api/internal/orders"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	method := event.RequestContext.HTTP.Method
	orderID := event.PathParameters["oid"]

	switch {
	case method == http.MethodGet && orderID == "":
		return orders.GetAll(ctx, event)
	case method == http.MethodGet && orderID != "":
		return orders.GetByID(ctx, event)
	case method == http.MethodPost && orderID == "":
		return orders.Create(ctx, event)
	default:
		return h.BadRequest(&h.Error{
			Message: "method not allowed",
		}), nil
	}
}

func main() {
	lambda.Start(handler)
}
