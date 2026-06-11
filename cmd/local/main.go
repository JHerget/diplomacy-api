package main

import (
	"context"
	"diplomacy-api/internal/lambdas/board"
	"diplomacy-api/internal/lambdas/orders"
	"diplomacy-api/internal/lambdas/test"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type lambdaHandler func(context.Context, events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /test", adapt(test.Handler))
	mux.HandleFunc("GET /games/{gid}/board", adapt(board.Handler))
	mux.HandleFunc("POST /games/{gid}/turns/{tid}/orders", adapt(orders.Handler))

	log.Println("local API listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func adapt(handler lambdaHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		event, err := toAPIGatewayV2Request(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res, err := handler(r.Context(), event)
		if err != nil && res.StatusCode == 0 {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		writeAPIGatewayV2Response(w, res)
	}
}

func toAPIGatewayV2Request(r *http.Request) (events.APIGatewayV2HTTPRequest, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return events.APIGatewayV2HTTPRequest{}, err
	}

	headers := make(map[string]string, len(r.Header))
	for name, values := range r.Header {
		headers[strings.ToLower(name)] = strings.Join(values, ",")
	}

	query := make(map[string]string, len(r.URL.Query()))
	for name, values := range r.URL.Query() {
		query[name] = strings.Join(values, ",")
	}

	params := map[string]string{}
	for _, name := range []string{"gid", "tid"} {
		if value := r.PathValue(name); value != "" {
			params[name] = value
		}
	}

	return events.APIGatewayV2HTTPRequest{
		RawPath:               r.URL.Path,
		Headers:               headers,
		QueryStringParameters: query,
		PathParameters:        params,
		Body:                  string(body),
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method: r.Method,
				Path:   r.URL.Path,
			},
		},
	}, nil
}

func writeAPIGatewayV2Response(w http.ResponseWriter, res events.APIGatewayV2HTTPResponse) {
	for name, value := range res.Headers {
		w.Header().Set(name, value)
	}

	status := res.StatusCode
	if status == 0 {
		status = http.StatusOK
	}

	w.WriteHeader(status)
	_, _ = w.Write([]byte(res.Body))
}
