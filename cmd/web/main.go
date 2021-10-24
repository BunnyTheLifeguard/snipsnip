package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/BunnyTheLifeguard/snipsnip/pkg/models/mongodb"
	"github.com/golangcollege/sessions"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session
	snips         *mongodb.SnipModel
	templateCache map[string]*template.Template
}

func init() {
	err := godotenv.Load("../../.env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	port := os.Getenv("PORT")
	dbName := os.Getenv("DB")
	collName := os.Getenv("COLLECTION")
	sessionSecret := os.Getenv("SECRET")

	addr := flag.String("addr", port, "HTTP network address")
	secret := flag.String("secret", sessionSecret, "Secret")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Connect to DB
	uri := os.Getenv("MONGODB_URI")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	db, err := openDB(ctx, uri)
	if err != nil {
		errorLog.Fatal(err)
	}
	coll := openCollection(db, dbName, collName)

	defer db.Disconnect(ctx)

	templateCache, err := newTemplateCache("../../ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		session:       session,
		snips:         &mongodb.SnipModel{Collection: coll},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	errSrv := srv.ListenAndServeTLS("../../tls/cert.pem", "../../tls/key.pem")
	errorLog.Fatal(errSrv)
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
	fmt.Println("Connected to MongoDB")

	return client, nil
}

func openCollection(client *mongo.Client, dbName, collectionName string) *mongo.Collection {
	coll := client.Database(dbName).Collection(collectionName)
	return coll
}
