package main

import (
	"context"
	h "diplomacy-api/internal/http"
	"diplomacy-api/internal/turns"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	switch event.RouteKey {
	case "GET /games/{gid}/turns":
		return turns.GetAll(ctx, event)
	case "GET /games/{gid}/turns/{tid}":
		return turns.GetByID(ctx, event)
	case "POST /games/{gid}/turns":
		return turns.Create(ctx, event)
	case "PUT /games/{gid}/turns/{tid}":
		return turns.Update(ctx, event)
	case "DELETE /games/{gid}/turns/{tid}":
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
