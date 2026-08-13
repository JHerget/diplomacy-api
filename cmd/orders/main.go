package main

import (
	"context"
	h "diplomacy-api/internal/http"
	"diplomacy-api/internal/orders"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	switch event.RouteKey {
	case "GET /games/{gid}/turns/{tid}/orders":
		return orders.GetAll(ctx, event)
	case "GET /games/{gid}/turns/{tid}/orders/{oid}":
		return orders.GetByID(ctx, event)
	case "POST /games/{gid}/turns/{tid}/orders":
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
