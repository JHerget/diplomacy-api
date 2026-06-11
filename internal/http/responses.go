package http

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

func Ok[T any](body *T) events.APIGatewayV2HTTPResponse {
	return ResponseWithBody(200, body)
}

func BadRequest[T any](body *T) events.APIGatewayV2HTTPResponse {
	return ResponseWithBody(400, body)
}

func InternalServerError[T any](body *T) events.APIGatewayV2HTTPResponse {
	return ResponseWithBody(500, body)
}

func ResponseNoBody(statusCode int) events.APIGatewayV2HTTPResponse {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: statusCode,
	}
}

func ResponseWithBody[T any](statusCode int, body *T) events.APIGatewayV2HTTPResponse {
	if body == nil {
		return ResponseNoBody(statusCode)
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return ResponseNoBody(statusCode)
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: statusCode,
		Body: string(jsonBody),
		Headers: map[string]string {
			"Content-Type": "application/json",
		},
	}
}
