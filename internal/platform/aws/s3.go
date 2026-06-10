package aws

import (
	"context"
	"io"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3 struct {
	config config.Config
	client *s3.Client
}

func NewS3(ctx context.Context) S3 {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Println("failed to load default s3 config", err)
	}

	return S3{
		config: cfg,
		client: s3.NewFromConfig(cfg),
	}
}

func (s *S3) Get(ctx context.Context, bucket string, filename string) ([]byte, error) {
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		log.Printf("failed to get file '%s': %v", filename, err)
		return nil, err
	}
	defer result.Body.Close()

	buf, err := io.ReadAll(result.Body)
	if err != nil {
		log.Printf("failed to read body of file '%s': %v", filename, err)
		return nil, err
	}

	return buf, nil
}
