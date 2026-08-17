package main

import (
	"context"
	h "diplomacy-api/internal/http"
	"diplomacy-api/internal/maps"
	"diplomacy-api/internal/platform/aws"
	"diplomacy-api/internal/platform/mongo"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(mapHandler *maps.Handler) func(context.Context, events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return func(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		switch event.RouteKey {
		case "GET /maps":
			return mapHandler.GetAll(ctx, event)
		case "GET /maps/{mid}":
			return mapHandler.GetByID(ctx, event)
		case "GET /maps/{mid}/image":
			return mapHandler.GetImage(ctx, event)
		case "POST /maps":
			return mapHandler.Create(ctx, event)
		case "PUT /maps/{mid}":
			return mapHandler.Update(ctx, event)
		case "DELETE /maps/{mid}":
			return mapHandler.Delete(ctx, event)
		default:
			return h.BadRequest(&h.Error{
				Message: "method not allowed",
			}), nil
		}
	}
}

func main() {
	ctx := context.Background()
	db, err := mongo.NewMongoDB(ctx)
	if err != nil {
		log.Fatal(err)
	}
	s3, err := aws.NewS3(ctx)
	if err != nil {
		log.Fatal(err)
	}

	mapHandler := maps.NewHandler(maps.NewRepository(db), s3)
	lambda.Start(handler(mapHandler))
}
