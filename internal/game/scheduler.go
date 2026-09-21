package game

import (
	"context"
	"diplomacy-api/internal/models"
	"diplomacy-api/internal/platform/aws"
	"encoding/json"
	"fmt"
	"time"
)

type TurnScheduler struct {
	scheduler *aws.Scheduler
}

func NewTurnScheduler(scheduler *aws.Scheduler) *TurnScheduler {
	return &TurnScheduler{
		scheduler: scheduler,
	}
}

func (ts *TurnScheduler) Create(ctx context.Context, game *models.Game) error {
	input, err := json.Marshal(map[string]string{
		"gameId": game.ID,
	})
	if err != nil {
		return err
	}

	return ts.scheduler.Create(ctx, aws.Schedule{
		Name:       turnScheduleName(game.ID),
		Expression: turnScheduleExpression(game),
		StartDate:  time.Unix(int64(game.StartDate), 0).UTC(),
		Input:      string(input),
	})
}

func (ts *TurnScheduler) Delete(ctx context.Context, gameID string) error {
	return ts.scheduler.Delete(ctx, turnScheduleName(gameID))
}

func turnScheduleName(gameID string) string {
	return fmt.Sprintf("diplomacy-game-%s", gameID)
}

func turnScheduleExpression(game *models.Game) string {
	utcHour := (game.TurnStartHour - game.Timezone) % 24
	if utcHour < 0 {
		utcHour += 24
	}

	return fmt.Sprintf("cron(0 %d * * ? *)", utcHour)
}
