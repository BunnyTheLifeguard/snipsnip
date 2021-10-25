package mongodb

import (
	"context"
	"strings"
	"time"

	"github.com/BunnyTheLifeguard/snipsnip/pkg/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// UserModel wraps the db collection
type UserModel struct {
	Collection *mongo.Collection
}

// Insert adds a new user
func (m *UserModel) Insert(name, email, password string) error {
	oid := primitive.NewObjectID()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	user := models.User{
		OID:            oid,
		ID:             oid.Hex(),
		Name:           name,
		Email:          email,
		HashedPassword: hashedPassword,
		Created:        time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = m.Collection.InsertOne(ctx, user)
	if err != nil {
		mongoErr := err.(mongo.WriteException)
		if strings.Contains(mongoErr.WriteErrors[0].Message, "name") {
			return models.ErrDuplicateName
		} else if strings.Contains(mongoErr.WriteErrors[0].Message, "email") {
			return models.ErrDuplicateEmail
		}
		return err
	}

	return err
}

// Authenticate verifies if the user exists and returns the user id
func (m *UserModel) Authenticate(email, password string) (string, error) {
	return "", nil
}

// Get fetches details of user with specified user id
func (m *UserModel) Get(id string) (*models.User, error) {
	return nil, nil
}
