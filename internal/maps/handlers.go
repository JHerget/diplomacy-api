package maps

import (
	"context"
	h "diplomacy-api/internal/http"
	"diplomacy-api/internal/models"
	"diplomacy-api/internal/platform/aws"
	"diplomacy-api/internal/platform/mongo"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func GetAll(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	db, err := mongo.NewMongoDB(ctx)
	if err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), err
	}

	mapRepo := NewRepository(db)
	summaries, err := mapRepo.GetAllSummaries(ctx)
	if err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), err
	}

	return h.OK(summaries), nil
}

func GetByID(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	mid := event.PathParameters["mid"]

	db, err := mongo.NewMongoDB(ctx)
	if err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), err
	}

	mapRepo := NewRepository(db)
	m, err := mapRepo.Get(ctx, mid)
	if err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), err
	}

	return h.OK(m), nil
}

func Create(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	var req createMapRequest
	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), err
	}

	db, err := mongo.NewMongoDB(ctx)
	if err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), err
	}

	mapRepo := NewRepository(db)
	m := models.Map{
		Name:        req.Name,
		Players:     req.Players,
		Providences: req.Providences,
		IsDeleted:   false,
	}

	if err := m.Valid(); err != nil {
		return h.BadRequest(&h.Error{
			Message: err.Error(),
		}), err
	}

	mapRepo.Create(ctx, &m)

	return h.Created(&m), nil
}

func Update(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	var m models.Map
	if err := json.Unmarshal([]byte(event.Body), &m); err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), err
	}

	if err := m.Valid(); err != nil {
		return h.BadRequest(&h.Error{
			Message: err.Error(),
		}), err
	}

	db, err := mongo.NewMongoDB(ctx)
	if err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), err
	}

	mapRepo := NewRepository(db)
	if err := mapRepo.Update(ctx, &m); err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), err
	}

	return h.OK(m), nil
}

func Delete(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       "delete map",
	}, nil
}

func GetImage(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	mid := event.PathParameters["mid"]

	s3, err := aws.NewS3(ctx)
	if err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), err
	}

	filename := fmt.Sprintf("%s.png", mid)
	buf, err := s3.Get(ctx, "diplomacy-maps", filename)
	if err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), err
	}

	return h.OK(buf), nil
}

func SetImage(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	mid := event.PathParameters["mid"]

	db, err := mongo.NewMongoDB(ctx)
	if err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), err
	}

	mapRepo := NewRepository(db)
	m, err := mapRepo.Get(ctx, mid)
	if err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), err
	}

	var buf []byte

	if event.IsBase64Encoded {
		buf, err = base64.StdEncoding.DecodeString(event.Body)
		if err != nil {
			return h.BadRequest(&h.Error{
				Message: err.Error(),
			}), err
		}
	} else {
		buf = []byte(event.Body)
	}

	s3, err := aws.NewS3(ctx)
	if err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), err
	}

	filename := fmt.Sprintf("%s.png", m.ID)
	if err := s3.Put(ctx, "diplomacy-maps", filename, buf); err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), err
	}

	return h.OK(buf), nil
}

type createMapRequest struct {
	Name        string              `json:"name"`
	Players     []models.Player     `json:"players"`
	Providences []models.Providence `json:"providences"`
}
