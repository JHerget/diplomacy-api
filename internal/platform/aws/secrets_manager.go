package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type SecretsManager struct {
	config config.Config
	client secretsmanager.Client
}

func NewSecretsManager(ctx context.Context) SecretsManager {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {

	}
	return SecretsManager{}
}
