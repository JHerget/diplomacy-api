package game

import (
	"context"
	"diplomacy-api/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	collection *mongo.Collection
}

type GameFilter struct {
	ExternalID *string
	InProgress *bool
}

func NewRepository(db *mongo.Client) *Repository {
	return &Repository{
		collection: db.Database("diplomacy").Collection("games"),
	}
}

func (r *Repository) Get(ctx context.Context, id string) (*models.Game, error) {
	var game models.Game

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = r.collection.FindOne(ctx, bson.M{
		"_id": objectID,
	}).Decode(&game)
	if err != nil {
		return nil, err
	}

	return &game, nil
}

func (r *Repository) GetAll(ctx context.Context, filter GameFilter) ([]*models.Game, error) {
	var games []*models.Game

	query := bson.M{
		"isDeleted": bson.M{"$ne": true},
	}
	if filter.ExternalID != nil {
		query["externalId"] = *filter.ExternalID
	}
	if filter.InProgress != nil {
		query["inProgress"] = *filter.InProgress
	}

	cursor, err := r.collection.Find(ctx, query)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &games); err != nil {
		return nil, err
	}

	return games, nil
}

func (r *Repository) Create(ctx context.Context, game *models.Game) error {
	result, err := r.collection.InsertOne(ctx, game)
	if err != nil {
		return err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		game.ID = oid.Hex()
	}

	return nil
}

func (r *Repository) Update(ctx context.Context, game *models.Game) error {
	objectID, err := primitive.ObjectIDFromHex(game.ID)
	if err != nil {
		return err
	}

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{
			"$set": bson.M{
				"ownerID":       game.OwnerID,
				"map":           game.Map,
				"board":         game.Board,
				"players":       game.Players,
				"turns":         game.Turns,
				"daysPerTurn":   game.DaysPerTurn,
				"turnStartHour": game.TurnStartHour,
				"timezone":      game.Timezone,
				"startDate":     game.StartDate,
				"endDate":       game.EndDate,
				"inProgress":    game.InProgress,
				"isDeleted":     game.IsDeleted,
			},
		},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{"isDeleted": true}},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
