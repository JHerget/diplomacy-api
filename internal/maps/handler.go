package maps

import (
	"context"
	"diplomacy-api/internal/http"
	"diplomacy-api/internal/models"
	"diplomacy-api/internal/platform/aws"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

type Handler struct {
	mapRepo *Repository
	s3      *aws.S3
}

func NewHandler(mapRepo *Repository, s3 *aws.S3) *Handler {
	return &Handler{
		mapRepo: mapRepo,
		s3:      s3,
	}
}

func (h *Handler) GetAll(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	summaries, err := h.mapRepo.GetAllSummaries(ctx)
	if err != nil {
		return http.InternalServerError(err), err
	}

	return http.OK(summaries), nil
}

func (h *Handler) GetByID(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	mid := event.PathParameters["mid"]

	m, err := h.mapRepo.Get(ctx, mid)
	if err != nil {
		return http.InternalServerError(err), err
	}

	return http.OK(m), nil
}

func (h *Handler) Create(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	var req createMapRequest
	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		return http.InternalServerError(err), err
	}

	m := models.Map{
		Name:        req.Name,
		Players:     req.Players,
		Providences: req.Providences,
		IsDeleted:   false,
	}

	if err := m.Valid(); err != nil {
		return http.BadRequest(&http.Error{
			Message: err.Error(),
		}), err
	}

	h.mapRepo.Create(ctx, &m)

	return http.Created(&m), nil
}

func (h *Handler) Update(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	var m models.Map
	if err := json.Unmarshal([]byte(event.Body), &m); err != nil {
		return http.InternalServerError(err), err
	}

	if err := m.Valid(); err != nil {
		return http.BadRequest(&http.Error{
			Message: err.Error(),
		}), err
	}

	if err := h.mapRepo.Update(ctx, &m); err != nil {
		return http.InternalServerError(err), err
	}

	return http.OK(m), nil
}

func (h *Handler) Delete(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	mid := event.PathParameters["mid"]

	if err := h.mapRepo.Delete(ctx, mid); err != nil {
		return http.InternalServerError(err), err
	}

	return http.NoContent(), nil
}

func (h *Handler) GetImage(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	mid := event.PathParameters["mid"]

	filename := fmt.Sprintf("%s.png", mid)
	buf, err := h.s3.Get(ctx, "diplomacy-maps", filename)
	if err != nil {
		return http.InternalServerError(err), err
	}

	return http.OK(buf), nil
}

func (h *Handler) SetImage(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	mid := event.PathParameters["mid"]

	m, err := h.mapRepo.Get(ctx, mid)
	if err != nil {
		return http.InternalServerError(err), err
	}

	var buf []byte

	if event.IsBase64Encoded {
		buf, err = base64.StdEncoding.DecodeString(event.Body)
		if err != nil {
			return http.BadRequest(&http.Error{
				Message: err.Error(),
			}), err
		}
	} else {
		buf = []byte(event.Body)
	}

	filename := fmt.Sprintf("%s.png", m.ID)
	if err := h.s3.Put(ctx, "diplomacy-maps", filename, buf); err != nil {
		return http.InternalServerError(err), err
	}

	return http.OK(buf), nil
}

type createMapRequest struct {
	Name        string              `json:"name"`
	Players     []models.Player     `json:"players"`
	Providences []models.Providence `json:"providences"`
}
