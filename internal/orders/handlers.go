package orders

import (
	"context"
	"diplomacy-api/internal/board"
	"diplomacy-api/internal/game"
	h "diplomacy-api/internal/http"
	"diplomacy-api/internal/models"
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
	gameID := event.PathParameters["gid"]
	turnID := event.PathParameters["tid"]

	var order orderRequest
	err := json.Unmarshal([]byte(event.Body), &order)
	if err != nil {
		return h.BadRequest(&h.Error{
			Message: "unable to parse request body",
		}), err
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
	g, err := gameRepo.Get(ctx, gameID)
	if err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), err
	}

	turn, ok := utils.Find(g.Turns, func(t *models.Turn) bool {
		return t.ID == turnID
	})
	if !ok {
		return h.BadRequest(&h.Error{
			Message: fmt.Sprintf("invalid turn id '%s'", turnID),
		}), fmt.Errorf("invalid turn id '%s'", turnID)
	}

	turn.Orders = append(turn.Orders, models.Order{
		ID:          "",
		PhaseID:     order.PhaseID,
		PlayerName:  order.PlayerName,
		CreatedDate: 0,
		Value:       order.Value,
	})

	allOrders := utils.Filter(turn.Orders, func(o *models.Order) bool {
		return o.PhaseID == turn.PhaseID
	})
	allCommands := GetCommands(allOrders)
	validCommands := ValidateCommands(allCommands, g.Board)

	if turn.TurnNumber%2 != 0 {
		validCommands.Reinforce = nil
		validCommands.Disband = nil
	}

	newBoard, err := board.ApplyTurn(g.Board, validCommands)
	if err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), err
	}
	g.Board = newBoard

	filename := fmt.Sprintf("%s.png", g.Map.ID)
	buf, err := s3.Get(ctx, "diplomacy-maps", filename)
	if err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), err
	}

	buf, err = board.Draw(buf, g)
	if err != nil {
		return h.InternalServerError(&h.Error{
			Message: err.Error(),
		}), err
	}

	return h.OK(&buf), nil
}

type orderRequest struct {
	PlayerName string `json:"playerName"`
	PhaseID    string `json:"phaseID"`
	Value      string `json:"value"`
}
