package maps

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

func NewRepository(db *mongo.Client) *Repository {
	return &Repository{
		collection: db.Database("diplomacy").Collection("maps"),
	}
}

func (r *Repository) Get(ctx context.Context, id string) (*models.Map, error) {
	var m models.Map

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = r.collection.FindOne(ctx, bson.M{
		"_id": objectID,
	}).Decode(&m)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (r *Repository) GetAllSummaries(ctx context.Context) ([]*models.MapSummary, error) {
	var summaries []*models.MapSummary

	cursor, err := r.collection.Find(ctx, bson.M{
		"isDeleted": bson.M{"$ne": true},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &summaries); err != nil {
		return nil, err
	}

	return summaries, nil
}

func (r *Repository) Create(ctx context.Context, m *models.Map) error {
	result, err := r.collection.InsertOne(ctx, m)
	if err != nil {
		return err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		m.ID = oid.Hex()
	}

	return nil
}

func (r *Repository) Update(ctx context.Context, m *models.Map) error {
	objectID, err := primitive.ObjectIDFromHex(m.ID)
	if err != nil {
		return err
	}

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{
			"$set": bson.M{
				"name":        m.Name,
				"players":     m.Players,
				"providences": m.Providences,
				"isDeleted":   m.IsDeleted,
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
