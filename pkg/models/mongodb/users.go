package mongodb

import (
	"context"
	"strings"
	"time"

	"github.com/BunnyTheLifeguard/snipsnip/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
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
	// var id string
	var result *models.User

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{Key: "email", Value: email}}
	err := m.Collection.FindOne(ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return "", models.ErrInvalidCredentials
	} else if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword(result.HashedPassword, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return "", models.ErrInvalidCredentials
	} else if err != nil {
		return "", err
	}

	return result.ID, err
}

// Get fetches details of user with specified user id
func (m *UserModel) Get(id string) (*models.User, error) {
	var result *models.User

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{Key: "_id", Value: oid}}
	resErr := m.Collection.FindOne(ctx, filter).Decode(&result)
	if resErr != nil {
		return nil, resErr
	}

	return result, nil
}
