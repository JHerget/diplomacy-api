package board

import (
	"context"
	"diplomacy-api/internal/game"
	h "diplomacy-api/internal/http"
	"diplomacy-api/internal/platform/aws"
	"diplomacy-api/internal/platform/mongo"

	"github.com/aws/aws-lambda-go/events"
)

func Get(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	gameId := event.PathParameters["gid"]

	db, err := mongo.NewMongoDB(ctx)
	if err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), err
	}

	s3, err := aws.NewS3(ctx)
	if err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), err
	}

	gameRepo := game.NewRepository(db)
	game, err := gameRepo.Get(ctx, gameId)
	if err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), err
	}

	buf, err := s3.Get(ctx, "diplomacy-maps", game.Map.Filename)
	if err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), err
	}

	buf, err = Draw(buf, game)
	if err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), nil
	}

	return h.OK(&buf), nil
}
