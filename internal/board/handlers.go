package board

import (
	"context"
	"diplomacy-api/internal/game"
	h "diplomacy-api/internal/http"
	"diplomacy-api/internal/platform/aws"
	"diplomacy-api/internal/platform/mongo"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

func Get(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	gameId := event.PathParameters["gid"]

	db, err := mongo.NewMongoDB(ctx)
	if err != nil {
		return h.InternalServerError(err), err
	}

	s3, err := aws.NewS3(ctx)
	if err != nil {
		return h.InternalServerError(err), err
	}

	gameRepo := game.NewRepository(db)
	g, err := gameRepo.Get(ctx, gameId)
	if err != nil {
		return h.InternalServerError(err), err
	}

	filename := fmt.Sprintf("%s.png", g.Map.ID)
	buf, err := s3.Get(ctx, "diplomacy-maps", filename)
	if err != nil {
		return h.InternalServerError(err), err
	}

	buf, err = Draw(buf, g)
	if err != nil {
		return h.InternalServerError(err), err
	}

	return h.OK(buf), nil
}
