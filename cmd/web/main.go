package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/BunnyTheLifeguard/snipsnip/pkg/models"
	"github.com/BunnyTheLifeguard/snipsnip/pkg/models/mongodb"
	"github.com/golangcollege/sessions"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type contextKey string

var contextKeyUser = contextKey("user")

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session
	templateCache map[string]*template.Template
	snips         interface {
		Insert(string, string, time.Time, time.Time) (interface{}, error)
		Get(string) (*models.Snip, error)
		Latest() ([]*models.Snip, error)
	}
	users interface {
		Insert(string, string, string) error
		Authenticate(string, string) (string, error)
		Get(string) (*models.User, error)
	}
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
	dataCollName := os.Getenv("DATA")
	userCollName := os.Getenv("USERS")
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
	defer db.Disconnect(ctx)

	// Get respective collections from DB
	dataColl := openCollection(db, dbName, dataCollName)
	userColl := openCollection(db, dbName, userCollName)

	// Create TTL & unique indexes
	_, err = dataColl.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "expires", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(1),
		},
	)
	if err != nil {
		errorLog.Fatal(err)
	}

	_, err = userColl.Indexes().CreateMany(
		context.Background(),
		[]mongo.IndexModel{
			{
				Keys:    bson.D{{Key: "name", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
			{
				Keys:    bson.D{{Key: "email", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
		},
	)
	if err != nil {
		errorLog.Fatal(err)
	}

	templateCache, err := newTemplateCache("../../ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true
	session.SameSite = http.SameSiteStrictMode

	// Add to application dependencies
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		session:       session,
		snips:         &mongodb.SnipModel{Collection: dataColl},
		users:         &mongodb.UserModel{Collection: userColl},
		templateCache: templateCache,
	}

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		},
		MinVersion:       tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
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
