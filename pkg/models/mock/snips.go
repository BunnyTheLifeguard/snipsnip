package mock

import (
	"time"

	"github.com/BunnyTheLifeguard/snipsnip/pkg/models"
	"go.mongodb.org/mongo-driver/mongo"
)

var mockSnip = &models.Snip{
	ID:      "617ab3734902e994b0d646ec",
	Title:   "0 Test Snip",
	Content: "0 Test Content",
	Created: time.Now(),
	Expires: time.Now(),
}

type SnipModel struct{}

func (m *SnipModel) Insert(title, content string, expires, created time.Time) (interface{}, error) {
	return 2, nil
}

func (m *SnipModel) Get(id string) (*models.Snip, error) {
	switch id {
	case "617ab3734902e994b0d646ec":
		return mockSnip, nil
	default:
		return nil, mongo.ErrNoDocuments
	}
}

func (m *SnipModel) Latest() ([]*models.Snip, error) {
	return []*models.Snip{mockSnip}, nil
}
