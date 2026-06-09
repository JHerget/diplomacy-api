package aws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3 struct {
	config config.Config
	client *s3.Client
}

func NewS3() S3 {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Println("failed to load default s3 config", err)
	}

	return S3{
		config: cfg,
		client: s3.NewFromConfig(cfg),
	}
}

func (s *S3) Get() ([]byte, error) {
	return nil, nil
}
