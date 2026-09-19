package aws

import (
	"context"
	"diplomacy-api/internal/models"
	"encoding/json"
	"errors"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQS struct {
	config   config.Config
	client   *sqs.Client
	queueURL string
}

func NewSQS(ctx context.Context) (*SQS, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	return &SQS{
		config:   cfg,
		client:   sqs.NewFromConfig(cfg),
		queueURL: os.Getenv("EVENTS_QUEUE_URL"),
	}, nil
}

func (s *SQS) Send(ctx context.Context, message models.NotificationMessage) error {
	if s.queueURL == "" {
		return errors.New("missing EVENTS_QUEUE_URL")
	}

	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = s.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(s.queueURL),
		MessageBody: aws.String(string(body)),
	})
	return err
}
