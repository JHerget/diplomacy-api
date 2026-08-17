package phases

import (
	"context"
	h "diplomacy-api/internal/http"
	"diplomacy-api/internal/platform/mongo"

	"github.com/aws/aws-lambda-go/events"
)

func GetAll(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	db, err := mongo.NewMongoDB(ctx)
	if err != nil {
		return h.InternalServerError(err), err
	}

	phaseRepo := NewRepository(db)
	allPhases, err := phaseRepo.GetAll(ctx)
	if err != nil {
		return h.InternalServerError(err), err
	}

	return h.OK(allPhases), nil
}
