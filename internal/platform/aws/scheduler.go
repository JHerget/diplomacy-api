package aws

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awsscheduler "github.com/aws/aws-sdk-go-v2/service/scheduler"
	schedulertypes "github.com/aws/aws-sdk-go-v2/service/scheduler/types"
)

type Scheduler struct {
	client      *awsscheduler.Client
	roleARN     string
	eventBusARN string
}

type Schedule struct {
	Name       string
	Expression string
	StartDate  time.Time
	Source     string
	DetailType string
	Input      string
}

func NewScheduler(ctx context.Context) (*Scheduler, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	return &Scheduler{
		client:      awsscheduler.NewFromConfig(cfg),
		roleARN:     os.Getenv("TURN_SCHEDULE_ROLE_ARN"),
		eventBusARN: os.Getenv("TURN_SCHEDULE_EVENT_BUS_ARN"),
	}, nil
}

func (s *Scheduler) Create(ctx context.Context, schedule Schedule) error {
	if s.roleARN == "" {
		return errors.New("missing scheduler role ARN")
	}
	if s.eventBusARN == "" {
		return errors.New("missing scheduler event bus ARN")
	}

	_, err := s.client.CreateSchedule(ctx, &awsscheduler.CreateScheduleInput{
		Name:               aws.String(schedule.Name),
		ScheduleExpression: aws.String(schedule.Expression),
		StartDate:          aws.Time(schedule.StartDate),
		FlexibleTimeWindow: &schedulertypes.FlexibleTimeWindow{
			Mode: schedulertypes.FlexibleTimeWindowModeOff,
		},
		Target: &schedulertypes.Target{
			Arn:     aws.String(s.eventBusARN),
			RoleArn: aws.String(s.roleARN),
			EventBridgeParameters: &schedulertypes.EventBridgeParameters{
				DetailType: aws.String(schedule.DetailType),
				Source:     aws.String(schedule.Source),
			},
			Input: aws.String(schedule.Input),
		},
	})
	return err
}

func (s *Scheduler) Delete(ctx context.Context, name string) error {
	if s.client == nil {
		return errors.New("missing scheduler client")
	}

	_, err := s.client.DeleteSchedule(ctx, &awsscheduler.DeleteScheduleInput{
		Name: aws.String(name),
	})
	var notFound *schedulertypes.ResourceNotFoundException
	if errors.As(err, &notFound) {
		return nil
	}
	return err
}
