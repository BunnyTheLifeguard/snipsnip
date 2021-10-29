package mock

import (
	"time"

	"github.com/BunnyTheLifeguard/snipsnip/pkg/models"
	"go.mongodb.org/mongo-driver/mongo"
)

var mockUser = &models.User{
	ID:      "6176ec1c013d4f8b76647d69",
	Name:    "TestUser",
	Email:   "test@snipsnip.com",
	Created: time.Now(),
}

type UserModel struct{}

func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "dupe@snipsnip.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) Authenticate(email, password string) (string, error) {
	switch email {
	case "test@snipsnip.com":
		return "6176ec1c013d4f8b76647d69", nil
	default:
		return "", models.ErrInvalidCredentials
	}
}

func (m *UserModel) Get(id string) (*models.User, error) {
	switch id {
	case "6176ec1c013d4f8b76647d69":
		return mockUser, nil
	default:
		return nil, mongo.ErrNoDocuments
	}
}
