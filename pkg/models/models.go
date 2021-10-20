package models

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ErrNoRecord message
var ErrNoRecord = errors.New("models: no matching record found")

// Snip MongoDB model
type Snip struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Title   string             `bson:"title,omitempty"`
	Content string             `bson:"content,omitempty"`
	Created time.Time          `bson:"created,omitempty"`
	Expires time.Time          `bson:"expires,omitempty"`
}
