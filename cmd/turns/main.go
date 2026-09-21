package main

import (
	"context"
	"diplomacy-api/internal/game"
	h "diplomacy-api/internal/http"
	"diplomacy-api/internal/phases"
	platformaws "diplomacy-api/internal/platform/aws"
	"diplomacy-api/internal/platform/mongo"
	"diplomacy-api/internal/turns"
	"encoding/json"
	"errors"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func decodeTurnEvent(input json.RawMessage) (events.APIGatewayV2HTTPRequest, error) {
	var envelope struct {
		RouteKey string `json:"routeKey"`
		GameID   string `json:"gameId"`
	}
	if err := json.Unmarshal(input, &envelope); err != nil {
		return events.APIGatewayV2HTTPRequest{}, err
	}
	if envelope.RouteKey != "" {
		var event events.APIGatewayV2HTTPRequest
		if err := json.Unmarshal(input, &event); err != nil {
			return events.APIGatewayV2HTTPRequest{}, err
		}
		return event, nil
	}
	if envelope.GameID == "" {
		return events.APIGatewayV2HTTPRequest{}, errors.New("missing gameId")
	}
	return events.APIGatewayV2HTTPRequest{
		RouteKey:       "POST /games/{gid}/turns",
		PathParameters: map[string]string{"gid": envelope.GameID},
	}, nil
}

func handler(turnHandler *turns.Handler) func(context.Context, json.RawMessage) (events.APIGatewayV2HTTPResponse, error) {
	return func(ctx context.Context, input json.RawMessage) (events.APIGatewayV2HTTPResponse, error) {
		event, err := decodeTurnEvent(input)
		if err != nil {
			return h.BadRequest(&h.Error{Message: err.Error()}), err
		}
		switch event.RouteKey {
		case "GET /games/{gid}/turns":
			return turnHandler.GetAll(ctx, event)
		case "GET /games/{gid}/turns/{tid}":
			return turnHandler.GetByID(ctx, event)
		case "POST /games/{gid}/turns":
			return turnHandler.Create(ctx, event)
		case "PUT /games/{gid}/turns/{tid}":
			return turnHandler.Update(ctx, event)
		case "DELETE /games/{gid}/turns/{tid}":
			return turnHandler.Delete(ctx, event)
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

	sqs, err := platformaws.NewSQS(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	turnHandler := turns.NewHandler(game.NewRepository(db), phases.NewRepository(db), sqs)
	lambda.Start(handler(turnHandler))
}
