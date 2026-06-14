package game

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Client) *Repository {
	return &Repository{
		collection: db.Database("diplomacy").Collection("games"),
	}
}

func (r *Repository) Get(ctx context.Context, id string) (*Game, error) {
	var game Game

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = r.collection.FindOne(ctx, bson.M{
		"_id": objectId,
	}).Decode(&game)
	if err != nil {
		return nil, err
	}

	return &game, nil
}

func (r *Repository) GetAll(ctx context.Context) ([]*Game, error) {
	var games []*Game

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &games); err != nil {
		return nil, err
	}

	return games, nil
}
