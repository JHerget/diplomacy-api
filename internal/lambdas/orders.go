package lambdas

import (
	"context"
	"diplomacy-api/internal/board"
	"diplomacy-api/internal/game"
	"diplomacy-api/internal/http"
	"diplomacy-api/internal/platform/aws"
	"diplomacy-api/internal/platform/mongo"
	"diplomacy-api/internal/utils"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

type orderRequest struct {
	PlayerName string `json:"playerName"`
	PhaseId    string `json:"phaseId"`
	Value      string `json:"value"`
}

func OrdersHandler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	gameId := event.PathParameters["gid"]
	turnId := event.PathParameters["tid"]

	var order orderRequest
	err := json.Unmarshal([]byte(event.Body), &order)
	if err != nil {
		return http.BadRequest(&http.Error{
			Message: "unable to parse request body",
		}), nil
	}

	db, err := mongo.NewMongoDB(ctx)
	if err != nil {
		return http.InternalServerError(&http.Error{
			Message: err.Error(),
		}), err
	}

	s3, err := aws.NewS3(ctx)
	if err != nil {
		return http.InternalServerError(&http.Error{
			Message: err.Error(),
		}), err
	}

	gameRepo := game.NewRepository(db)
	g, err := gameRepo.Get(ctx, gameId)
	if err != nil {
		return http.InternalServerError(&http.Error{
			Message: err.Error(),
		}), err
	}

	turn, ok := utils.Find(g.Turns, func(t *game.Turn) bool {
		return t.Id == turnId
	})
	if !ok {
		return http.BadRequest(&http.Error{
			Message: fmt.Sprintf("invalid turn id '%s'", turnId),
		}), nil
	}

	turn.Orders = append(turn.Orders, game.Order{
		Id:          "",
		PhaseId:     order.PhaseId,
		PlayerName:  order.PlayerName,
		CreatedDate: 0,
		Value:       order.Value,
	})

	newBoard, err := board.ApplyTurn(g.Board, turn)
	if err != nil {
		return http.InternalServerError(&http.Error{
			Message: err.Error(),
		}), nil
	}
	g.Board = newBoard

	buf, err := s3.Get(ctx, "diplomacy-maps", g.Map.Filename)
	if err != nil {
		return http.InternalServerError(&http.Error{
			Message: err.Error(),
		}), err
	}

	buf, err = board.Draw(buf, g)
	if err != nil {
		return http.InternalServerError(&http.Error{
			Message: err.Error(),
		}), nil
	}

	return http.OK(&buf), nil
}
