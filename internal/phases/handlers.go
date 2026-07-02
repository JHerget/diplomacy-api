package phases

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func GetAll(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       "get all phases",
	}, nil
}
