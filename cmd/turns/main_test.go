package main

import (
	"encoding/json"
	"testing"
)

func TestDecodeTurnEvent(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		routeKey string
		gameID   string
		wantErr  bool
	}{
		{
			name:     "scheduled invocation",
			input:    `{"gameId":"game-123"}`,
			routeKey: "POST /games/{gid}/turns",
			gameID:   "game-123",
		},
		{
			name:     "API Gateway invocation",
			input:    `{"routeKey":"POST /games/{gid}/turns","pathParameters":{"gid":"game-456"}}`,
			routeKey: "POST /games/{gid}/turns",
			gameID:   "game-456",
		},
		{
			name:    "missing scheduled game ID",
			input:   `{}`,
			wantErr: true,
		},
		{
			name:    "invalid JSON",
			input:   `{`,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			event, err := decodeTurnEvent(json.RawMessage(test.input))
			if (err != nil) != test.wantErr {
				t.Fatalf("decodeTurnEvent() error = %v, wantErr = %v", err, test.wantErr)
			}
			if test.wantErr {
				return
			}
			if event.RouteKey != test.routeKey || event.PathParameters["gid"] != test.gameID {
				t.Fatalf("decodeTurnEvent() = route %q, game %q; want route %q, game %q", event.RouteKey, event.PathParameters["gid"], test.routeKey, test.gameID)
			}
		})
	}
}
