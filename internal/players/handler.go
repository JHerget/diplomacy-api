package players

import (
	"context"
	"diplomacy-api/internal/game"
	"diplomacy-api/internal/http"
	"diplomacy-api/internal/models"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

type Handler struct {
	gameRepo *game.Repository
}

func NewHandler(gameRepo *game.Repository) *Handler {
	return &Handler{
		gameRepo: gameRepo,
	}
}

func (h *Handler) GetAll(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	gameID := event.PathParameters["gid"]

	g, err := h.gameRepo.Get(ctx, gameID)
	if err != nil {
		return http.InternalServerError(err), err
	}

	return http.OK(g.Players), nil
}

func (h *Handler) GetByID(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	gameID := event.PathParameters["gid"]
	playerID := event.PathParameters["pid"]

	g, err := h.gameRepo.Get(ctx, gameID)
	if err != nil {
		return http.InternalServerError(err), err
	}

	player, ok := g.FindPlayer(playerID)
	if !ok {
		return invalidPlayerID(playerID)
	}

	return http.OK(player), nil
}

func (h *Handler) Create(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	gameID := event.PathParameters["gid"]

	var req createPlayerRequest
	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		return http.InternalServerError(err), err
	}
	if req.UserID == nil {
		err := fmt.Errorf("Missing user id.")
		return http.BadRequest(&http.Error{
			Message: err.Error(),
		}), err
	}

	g, err := h.gameRepo.Get(ctx, gameID)
	if err != nil {
		return http.InternalServerError(err), err
	}

	player, err := g.AssignPlayer(*req.UserID)
	if err != nil {
		return http.BadRequest(&http.Error{
			Message: err.Error(),
		}), err
	}

	if err := g.Valid(); err != nil {
		return http.BadRequest(&http.Error{
			Message: err.Error(),
		}), err
	}

	if err := h.gameRepo.Update(ctx, g); err != nil {
		return http.InternalServerError(err), err
	}

	return http.Created(*player), nil
}

func (h *Handler) Update(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	gameID := event.PathParameters["gid"]
	playerID := event.PathParameters["pid"]

	var player models.Player
	if err := json.Unmarshal([]byte(event.Body), &player); err != nil {
		return http.BadRequest(&http.Error{
			Message: err.Error(),
		}), err
	}
	player.ID = playerID

	g, err := h.gameRepo.Get(ctx, gameID)
	if err != nil {
		return http.InternalServerError(err), err
	}

	existing, ok := g.FindPlayer(playerID)
	if !ok {
		return invalidPlayerID(playerID)
	}
	*existing = player

	if err := g.Valid(); err != nil {
		return http.BadRequest(&http.Error{
			Message: err.Error(),
		}), err
	}

	if err := h.gameRepo.Update(ctx, g); err != nil {
		return http.InternalServerError(err), err
	}

	return http.OK(player), nil
}

func (h *Handler) Delete(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	gameID := event.PathParameters["gid"]
	playerID := event.PathParameters["pid"]

	g, err := h.gameRepo.Get(ctx, gameID)
	if err != nil {
		return http.InternalServerError(err), err
	}

	player, ok := g.FindPlayer(playerID)
	if !ok {
		return invalidPlayerID(playerID)
	}
	player.IsPlaying = false

	if err := g.Valid(); err != nil {
		return http.BadRequest(&http.Error{
			Message: err.Error(),
		}), err
	}

	if err := h.gameRepo.Update(ctx, g); err != nil {
		return http.InternalServerError(err), err
	}

	return http.NoContent(), nil
}

type createPlayerRequest struct {
	UserID *string `json:"userID"`
}

func invalidPlayerID(playerID string) (events.APIGatewayV2HTTPResponse, error) {
	err := fmt.Errorf("invalid player id '%s'", playerID)
	return http.BadRequest(&http.Error{
		Message: err.Error(),
	}), err
}
