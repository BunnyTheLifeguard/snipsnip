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
	OID     primitive.ObjectID `bson:"_id,omitempty"`
	ID      string             `json:"id" bson:"id,omitempty"`
	Title   string             `json:"title" bson:"title,omitempty"`
	Content string             `json:"content" bson:"content,omitempty"`
	Created time.Time          `json:"created" bson:"created,omitempty"`
	Expires time.Time          `json:"expires" bson:"expires,omitempty"`
}
