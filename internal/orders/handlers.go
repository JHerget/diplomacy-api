package orders

import (
	"context"
	"diplomacy-api/internal/board"
	"diplomacy-api/internal/game"
	h "diplomacy-api/internal/http"
	"diplomacy-api/internal/platform/aws"
	"diplomacy-api/internal/platform/mongo"
	"diplomacy-api/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func GetAll(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       "get all orders",
	}, nil
}

func GetByID(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       "get orders by id",
	}, nil
}

func Create(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	gameId := event.PathParameters["gid"]
	turnId := event.PathParameters["tid"]

	var order orderRequest
	err := json.Unmarshal([]byte(event.Body), &order)
	if err != nil {
		return h.BadRequest(&h.Error{
			Message: "unable to parse request body",
		}), nil
	}

	db, err := mongo.NewMongoDB(ctx)
	if err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), err
	}

	s3, err := aws.NewS3(ctx)
	if err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), err
	}

	gameRepo := game.NewRepository(db)
	g, err := gameRepo.Get(ctx, gameId)
	if err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), err
	}

	turn, ok := utils.Find(g.Turns, func(t *game.Turn) bool {
		return t.Id == turnId
	})
	if !ok {
		return h.BadRequest(&h.Error{
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
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), nil
	}
	g.Board = newBoard

	buf, err := s3.Get(ctx, "diplomacy-maps", g.Map.Filename)
	if err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), err
	}

	buf, err = board.Draw(buf, g)
	if err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), nil
	}

	return h.OK(&buf), nil
}

type orderRequest struct {
	PlayerName string `json:"playerName"`
	PhaseId    string `json:"phaseId"`
	Value      string `json:"value"`
}
