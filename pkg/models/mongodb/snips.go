package mongodb

import (
	"context"
	"time"

	"github.com/BunnyTheLifeguard/snipsnip/pkg/models"

	"go.mongodb.org/mongo-driver/mongo"
)

// SnipModel wraps the db collection
type SnipModel struct {
	Collection *mongo.Collection
}

// Insert adds a new snip to the DB and returns its ObjectID
func (m *SnipModel) Insert(title, content string, expires time.Time) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	snip := models.Snip{
		Title:   title,
		Content: content,
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
	return nil, nil
}

// Latest shows the last inserted snip
func (m *SnipModel) Latest() ([]*models.Snip, error) {
	return nil, nil
}
