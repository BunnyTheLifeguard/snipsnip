package models

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	// ErrInvalidCredentials message
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	// ErrDuplicateName message
	ErrDuplicateName = errors.New("models: duplicate name")
	// ErrDuplicateEmail message
	ErrDuplicateEmail = errors.New("models: duplicate email")
)

// Snip MongoDB model
type Snip struct {
	OID     primitive.ObjectID `bson:"_id,omitempty"`
	ID      string             `json:"id" bson:"id,omitempty"`
	Title   string             `json:"title" bson:"title,omitempty"`
	Content string             `json:"content" bson:"content,omitempty"`
	Created time.Time          `json:"created" bson:"created,omitempty"`
	Expires time.Time          `json:"expires" bson:"expires,omitempty"`
}

// User MongoDB model
type User struct {
	OID            primitive.ObjectID `bson:"_id,omitempty"`
	ID             string             `json:"id" bson:"id,omitempty"`
	Name           string             `json:"name" bson:"name,omitempty"`
	Email          string             `json:"email" bson:"email,omitempty"`
	HashedPassword []byte             `json:"password" bson:"password,omitempty"`
	Created        time.Time          `json:"created" bson:"created,omitempty"`
}
