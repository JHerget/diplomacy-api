package phases

import (
	"context"
	"diplomacy-api/internal/http"

	"github.com/aws/aws-lambda-go/events"
)

type Handler struct {
	phaseRepo *Repository
}

func NewHandler(phaseRepo *Repository) *Handler {
	return &Handler{
		phaseRepo: phaseRepo,
	}
}

func (h *Handler) GetAll(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	allPhases, err := h.phaseRepo.GetAll(ctx)
	if err != nil {
		return http.InternalServerError(err), err
	}

	return http.OK(allPhases), nil
}
