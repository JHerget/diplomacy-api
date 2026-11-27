package main

import (
	"encoding/json"
	"testing"
)

func TestScheduledTurnEventUsesGameID(t *testing.T) {
	event, err := decodeTurnEvent(json.RawMessage(`{"gameId":"game-123"}`))
	if err != nil {
		t.Fatal(err)
	}
	if got := event.RouteKey; got != "POST /games/{gid}/turns" {
		t.Errorf("route key = %q, want POST /games/{gid}/turns", got)
	}
	if got := event.PathParameters["gid"]; got != "game-123" {
		t.Errorf("gid = %q, want game-123", got)
	}
}

func TestAPIGatewayTurnEventKeepsRouteAndGameID(t *testing.T) {
	input := json.RawMessage(`{"routeKey":"POST /games/{gid}/turns","pathParameters":{"gid":"game-456"}}`)
	event, err := decodeTurnEvent(input)
	if err != nil {
		t.Fatal(err)
	}
	if got := event.RouteKey; got != "POST /games/{gid}/turns" {
		t.Errorf("route key = %q, want POST /games/{gid}/turns", got)
	}
	if got := event.PathParameters["gid"]; got != "game-456" {
		t.Errorf("gid = %q, want game-456", got)
	}
}

func TestScheduledTurnEventRequiresGameID(t *testing.T) {
	_, err := decodeTurnEvent(json.RawMessage(`{}`))
	if err == nil {
		t.Fatal("expected an error for a scheduled turn without a game ID")
	}
}

func TestTurnEventRejectsInvalidJSON(t *testing.T) {
	_, err := decodeTurnEvent(json.RawMessage(`{`))
	if err == nil {
		t.Fatal("expected an error for malformed JSON")
	}
}
