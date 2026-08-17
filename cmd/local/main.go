package main

import (
	"context"
	"diplomacy-api/internal/board"
	"diplomacy-api/internal/game"
	"diplomacy-api/internal/maps"
	"diplomacy-api/internal/orders"
	"diplomacy-api/internal/phases"
	"diplomacy-api/internal/platform/aws"
	"diplomacy-api/internal/platform/mongo"
	"diplomacy-api/internal/turns"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type lambdaHandler func(context.Context, events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error)

func main() {
	ctx := context.Background()
	db, err := mongo.NewMongoDB(ctx)
	if err != nil {
		log.Fatal(err)
	}
	s3, err := aws.NewS3(ctx)
	if err != nil {
		log.Fatal(err)
	}

	gameRepo := game.NewRepository(db)
	mapRepo := maps.NewRepository(db)
	phaseRepo := phases.NewRepository(db)
	boardHandler := board.NewHandler(gameRepo, s3)
	gameHandler := game.NewHandler(gameRepo, mapRepo)
	mapHandler := maps.NewHandler(mapRepo, s3)
	orderHandler := orders.NewHandler(gameRepo)
	turnHandler := turns.NewHandler(gameRepo, phaseRepo)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/maps", adapt(mapHandler.GetAll))
	mux.HandleFunc("GET /v1/maps/{mid}", adapt(mapHandler.GetByID))
	mux.HandleFunc("POST /v1/maps", adapt(mapHandler.Create))
	mux.HandleFunc("PUT /v1/maps/{mid}", adapt(mapHandler.Update))
	mux.HandleFunc("GET /v1/maps/{mid}/image", adapt(mapHandler.GetImage))
	mux.HandleFunc("POST /v1/maps/{mid}/image", adapt(mapHandler.SetImage))

	mux.HandleFunc("GET /v1/games", adapt(gameHandler.GetAll))
	mux.HandleFunc("GET /v1/games/{gid}", adapt(gameHandler.GetByID))
	mux.HandleFunc("POST /v1/games", adapt(gameHandler.Create))
	mux.HandleFunc("PUT /v1/games/{gid}", adapt(gameHandler.Update))
	mux.HandleFunc("DELETE /v1/games/{gid}", adapt(gameHandler.Delete))
	mux.HandleFunc("GET /v1/games/{gid}/board", adapt(boardHandler.Get))
	mux.HandleFunc("GET /v1/games/{gid}/turns", adapt(turnHandler.GetAll))
	mux.HandleFunc("GET /v1/games/{gid}/turns/{tid}", adapt(turnHandler.GetByID))
	mux.HandleFunc("POST /v1/games/{gid}/turns", adapt(turnHandler.Create))
	mux.HandleFunc("PUT /v1/games/{gid}/turns/{tid}", adapt(turnHandler.Update))
	mux.HandleFunc("DELETE /v1/games/{gid}/turns/{tid}", adapt(turnHandler.Delete))
	mux.HandleFunc("GET /v1/games/{gid}/turns/{tid}/orders", adapt(orderHandler.GetAll))
	mux.HandleFunc("GET /v1/games/{gid}/turns/{tid}/orders/{oid}", adapt(orderHandler.GetByID))
	mux.HandleFunc("POST /v1/games/{gid}/turns/{tid}/orders", adapt(orderHandler.Create))

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
	for _, name := range []string{"gid", "tid", "oid", "mid"} {
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

	body := []byte(res.Body)
	if res.IsBase64Encoded {
		decodedBody, err := base64.StdEncoding.DecodeString(res.Body)
		if err != nil {
			http.Error(w, "failed to decode response body", http.StatusInternalServerError)
			return
		}
		body = decodedBody
	}

	w.WriteHeader(status)
	_, _ = w.Write(body)
}
