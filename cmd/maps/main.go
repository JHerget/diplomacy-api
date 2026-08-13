package main

import (
	"context"
	h "diplomacy-api/internal/http"
	"diplomacy-api/internal/maps"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	switch event.RouteKey {
	case "GET /maps":
		return maps.GetAll(ctx, event)
	case "GET /maps/{mid}":
		return maps.GetByID(ctx, event)
	case "GET /maps/{mid}/image":
		return maps.GetImage(ctx, event)
	case "POST /maps":
		return maps.Create(ctx, event)
	case "PUT /maps/{mid}":
		return maps.Update(ctx, event)
	case "DELETE /maps/{mid}":
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
