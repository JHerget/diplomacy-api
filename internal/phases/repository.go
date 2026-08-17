package phases

import (
	"context"
	"diplomacy-api/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Client) *Repository {
	return &Repository{
		collection: db.Database("diplomacy").Collection("phases"),
	}
}

func (r *Repository) GetAll(ctx context.Context) ([]*models.Phase, error) {
	var phases []*models.Phase

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &phases); err != nil {
		return nil, err
	}

	return phases, nil
}

func (r *Repository) GetByOrder(ctx context.Context, phaseOrder int8) (*models.Phase, error) {
	var phase models.Phase

	err := r.collection.FindOne(ctx, bson.M{
		"phaseOrder": phaseOrder,
	}).Decode(&phase)
	if err != nil {
		return nil, err
	}

	return &phase, nil
}
