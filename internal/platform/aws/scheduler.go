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
	client    *awsscheduler.Client
	roleARN   string
	lambdaARN string
}

type Schedule struct {
	Name       string
	Expression string
	StartDate  time.Time
	Input      string
}

func NewScheduler(ctx context.Context) (*Scheduler, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	roleARN := os.Getenv("TURN_SCHEDULE_ROLE_ARN")
	if roleARN == "" {
		return nil, errors.New("missing TURN_SCHEDULE_ROLE_ARN")
	}
	lambdaARN := os.Getenv("TURN_SCHEDULE_LAMBDA_ARN")
	if lambdaARN == "" {
		return nil, errors.New("missing TURN_SCHEDULE_LAMBDA_ARN")
	}

	return &Scheduler{
		client:    awsscheduler.NewFromConfig(cfg),
		roleARN:   roleARN,
		lambdaARN: lambdaARN,
	}, nil
}

func (s *Scheduler) Create(ctx context.Context, schedule Schedule) error {
	_, err := s.client.CreateSchedule(ctx, &awsscheduler.CreateScheduleInput{
		Name:               aws.String(schedule.Name),
		ScheduleExpression: aws.String(schedule.Expression),
		StartDate:          aws.Time(schedule.StartDate),
		FlexibleTimeWindow: &schedulertypes.FlexibleTimeWindow{
			Mode: schedulertypes.FlexibleTimeWindowModeOff,
		},
		Target: &schedulertypes.Target{
			Arn:     aws.String(s.lambdaARN),
			RoleArn: aws.String(s.roleARN),
			Input:   aws.String(schedule.Input),
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
