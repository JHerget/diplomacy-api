package board

import (
	"context"
	"diplomacy-api/internal/game"
	"diplomacy-api/internal/http"
	"diplomacy-api/internal/platform/aws"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

type Handler struct {
	gameRepo *game.Repository
	s3       *aws.S3
}

func NewHandler(gameRepo *game.Repository, s3 *aws.S3) *Handler {
	return &Handler{
		gameRepo: gameRepo,
		s3:       s3,
	}
}

func (h *Handler) Get(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	gameId := event.PathParameters["gid"]

	g, err := h.gameRepo.Get(ctx, gameId)
	if err != nil {
		return http.InternalServerError(err), err
	}

	filename := fmt.Sprintf("%s.png", g.Map.ID)
	buf, err := h.s3.Get(ctx, "diplomacy-maps", filename)
	if err != nil {
		return http.InternalServerError(err), err
	}

	buf, err = Draw(buf, g)
	if err != nil {
		return http.InternalServerError(err), err
	}

	return http.OK(buf), nil
}
