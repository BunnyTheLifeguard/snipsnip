package mongodb

import (
	"github.com/BunnyTheLifeguard/snipsnip/pkg/models"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserModel wraps the db collection
type UserModel struct {
	Collection *mongo.Collection
}

// Insert adds a new user
func (m *UserModel) Insert(name, email, password string) error {
	return nil
}

// Authenticate verifies if the user exists and returns the user id
func (m *UserModel) Authenticate(email, password string) (string, error) {
	return "", nil
}

// Get fetches details of user with specified user id
func (m *UserModel) Get(id string) (*models.User, error) {
	return nil, nil
}
