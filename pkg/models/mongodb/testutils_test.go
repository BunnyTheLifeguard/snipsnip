package mongodb

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func newTestDB(t *testing.T) (*mongo.Collection, func()) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		t.Fatal("Error loading .env file")
	}
	uri := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("DB")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	db, err := openDB(ctx, uri)
	if err != nil {
		t.Fatal(err)
	}

	testUserColl := openCollection(db, dbName, "test-users")
	oid, _ := primitive.ObjectIDFromHex("6176ec1c013d4f8b76612345")

	testUser := bson.D{
		{Key: "_id", Value: oid},
		{Key: "id", Value: "6176ec1c013d4f8b76612345"},
		{Key: "name", Value: "Cpt. Teemo"},
		{Key: "email", Value: "teemo@snipsnip.com"},
		{Key: "hashedPassword", Value: "$2a$12$NuTjWXm3KKntReFwyBVHyuf/to.HEwTy.eS206TNfkGfr6HzGJSWG"},
		{Key: "created", Value: time.Date(2018, time.December, 23, 17, 25, 22, 0, time.UTC)},
	}

	// Insert testuser into test collection
	_, err = testUserColl.InsertOne(ctx, testUser)
	if err != nil {
		t.Fatal(err)
	}

	// Return testcollection and anonymous function which drops the collection after testing
	return testUserColl, func() {
		err := testUserColl.Drop(ctx)
		if err != nil {
			t.Fatal(err)
		}

		db.Disconnect(ctx)
	}
}

func openDB(ctx context.Context, uri string) (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return client, nil
}

func openCollection(client *mongo.Client, dbName, collectionName string) *mongo.Collection {
	coll := client.Database(dbName).Collection(collectionName)
	return coll
}
