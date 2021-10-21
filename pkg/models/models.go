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
	ID      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title   string             `json:"title" bson:"title,omitempty"`
	Content string             `json:"content" bson:"content,omitempty"`
	Created time.Time          `json:"created" bson:"created,omitempty"`
	Expires time.Time          `json:"expires" bson:"expires,omitempty"`
}
