package mongodb

import (
	"reflect"
	"testing"
	"time"

	"github.com/BunnyTheLifeguard/snipsnip/pkg/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestUserModelGet(t *testing.T) {
	// Skip testing if '-shot' flag provided
	if testing.Short() {
		t.Skip("mongodb: skipping integration test")
	}

	oid, _ := primitive.ObjectIDFromHex("6176ec1c013d4f8b76612345")

	tests := []struct {
		id        string
		name      string
		wantUser  *models.User
		wantError error
	}{
		{
			name: "Valid ID",
			id:   "6176ec1c013d4f8b76612345",
			wantUser: &models.User{
				OID:     oid,
				ID:      "6176ec1c013d4f8b76612345",
				Name:    "Cpt. Teemo",
				Email:   "teemo@snipsnip.com",
				Created: time.Date(2018, 12, 23, 17, 25, 22, 0, time.UTC),
			},
			wantError: nil,
		},
		{
			name:      "No ID",
			id:        "",
			wantUser:  nil,
			wantError: primitive.ErrInvalidHex,
		},
		{
			name:      "Wrong ID",
			id:        "6176ec1c013d4f8b76612340",
			wantUser:  nil,
			wantError: mongo.ErrNoDocuments,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			// Create an instance of the UserModel
			m := UserModel{db}

			user, err := m.Get(tt.id)

			if err != tt.wantError {
				t.Errorf("want %v got %s", tt.wantError, err)
			}

			if !reflect.DeepEqual(user, tt.wantUser) {
				t.Errorf("want %v got %v", tt.wantUser, user)
			}
		})
	}
}
