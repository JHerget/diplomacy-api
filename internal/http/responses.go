package http

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func OK[T any](body T) events.APIGatewayV2HTTPResponse {
	return ResponseWithBody(http.StatusOK, &body)
}

func BadRequest[T any](body T) events.APIGatewayV2HTTPResponse {
	return ResponseWithBody(http.StatusBadRequest, &body)
}

func InternalServerError[T any](body T) events.APIGatewayV2HTTPResponse {
	return ResponseWithBody(http.StatusInternalServerError, &body)
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

	switch v := any(*body).(type) {
	case []byte:
		return events.APIGatewayV2HTTPResponse{
			StatusCode: statusCode,
			Body:       base64.StdEncoding.EncodeToString(v),
			Headers: map[string]string{
				"Content-Type":  "image/png",
				"Cache-Control": "no-store",
			},
			IsBase64Encoded: true,
		}

	default:
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return ResponseNoBody(statusCode)
		}

		return events.APIGatewayV2HTTPResponse{
			StatusCode: statusCode,
			Body:       string(jsonBody),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
	}
}
