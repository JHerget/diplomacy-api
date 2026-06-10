package game

import (
	"context"
	"diplomacy-api/internal/board"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	collection *mongo.Collection
}

func MakeRepository(db *mongo.Client) *Repository {
	return &Repository{
		collection: db.Database("diplomacy").Collection("games"),
	}
}

func (r *Repository) Get(ctx context.Context, id string) (*board.Game, error) {
	var game board.Game

	err := r.collection.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&game)
	if err != nil {
		return nil, err
	}

	return &game, nil
}

func (r *Repository) GetAll(ctx context.Context) ([]*board.Game, error) {
	var games []*board.Game

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
