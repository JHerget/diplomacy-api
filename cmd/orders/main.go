package main

import (
	"context"
	"diplomacy-api/internal/game"
	h "diplomacy-api/internal/http"
	"diplomacy-api/internal/orders"
	"diplomacy-api/internal/platform/mongo"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(orderHandler *orders.Handler) func(context.Context, events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return func(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		switch event.RouteKey {
		case "GET /games/{gid}/turns/{tid}/orders":
			return orderHandler.GetAll(ctx, event)
		case "GET /games/{gid}/turns/{tid}/orders/{oid}":
			return orderHandler.GetByID(ctx, event)
		case "POST /games/{gid}/turns/{tid}/orders":
			return orderHandler.Create(ctx, event)
		default:
			return h.BadRequest(&h.Error{
				Message: "method not allowed",
			}), nil
		}
	}
}

func main() {
	db, err := mongo.NewMongoDB(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	orderHandler := orders.NewHandler(game.NewRepository(db))
	lambda.Start(handler(orderHandler))
}
