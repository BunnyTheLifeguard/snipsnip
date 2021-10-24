package mongodb

import (
	"context"
	"time"

	"github.com/BunnyTheLifeguard/snipsnip/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SnipModel wraps the db collection
type SnipModel struct {
	Collection *mongo.Collection
}

// Insert adds a new snip to the DB and returns its ObjectID
func (m *SnipModel) Insert(title, content string, created, expires time.Time) (interface{}, error) {
	oid := primitive.NewObjectID()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	snip := models.Snip{
		OID:     oid,
		ID:      oid.Hex(),
		Title:   title,
		Content: content,
		Created: created,
		Expires: expires,
	}
	defer cancel()

	insertResult, err := m.Collection.InsertOne(ctx, snip)
	if err != nil {
		return 0, err
	}

	return insertResult.InsertedID, nil
}

// Get a snip via its ObjectID
func (m *SnipModel) Get(id string) (*models.Snip, error) {
	var result *models.Snip
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{Key: "_id", Value: oid}}
	resErr := m.Collection.FindOne(ctx, filter).Decode(&result)
	if resErr != nil {
		if resErr == mongo.ErrNoDocuments {
			return nil, resErr
		}
		return nil, resErr
	}

	return result, err
}

// Latest shows the most recently created 10 snips unless expired
func (m *SnipModel) Latest() ([]*models.Snip, error) {
	var results []*models.Snip
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"expires": bson.M{"$gt": time.Now()}}
	opts := options.Find().SetLimit(10)
	cursor, err := m.Collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	resErr := cursor.All(ctx, &results)
	if resErr != nil {
		return nil, resErr
	}

	return results, nil
}
