package mongo

import (
	"context"
	"diplomacy-api/internal/platform/aws"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDB(ctx context.Context) (*mongo.Client, error) {
	secretsManager, err := aws.NewSecretsManager(ctx)
	if err != nil {
		return nil, err
	}

	secrets, err := secretsManager.GetJson(ctx, "diplomacy-credentials")
	if err != nil {
		return nil, err
	}

	opts := options.Client().ApplyURI(secrets.MongoDBConnectionString)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	return client, nil
}
