package aws

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type SecretsManager struct {
	config config.Config
	client *secretsmanager.Client
}

func NewSecretsManager(ctx context.Context) (*SecretsManager, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	client := secretsmanager.NewFromConfig(cfg)

	return &SecretsManager{
		config: cfg,
		client: client,
	}, nil
}

func (s *SecretsManager) Get(ctx context.Context, secretName string) (*string, error) {
	result, err := s.client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	})
	if err != nil {
		return nil, err
	}

	return result.SecretString, nil
}

func (s *SecretsManager) GetJson(ctx context.Context, secretName string) (*Secrets, error) {
	rawSecrets, err := s.Get(ctx, secretName)
	if err != nil {
		return nil, err
	}

	var secrets Secrets
	err = json.Unmarshal([]byte(*rawSecrets), &secrets)
	if err != nil {
		return nil, err
	}

	return &secrets, nil
}
