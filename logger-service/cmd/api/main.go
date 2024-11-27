package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/BlackSound1/go-microservices/logger/data"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	WEB_PORT  = "80"
	RPC_PORT  = "5001"
	MONGO_URL = "mongodb://mongo-service:27017"
	GRPC_PORT = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {

	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}

	client = mongoClient

	// Create context to allow us to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Disconnect whenver it is time to
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	// Start web server
	app.serve()
}

// serve starts the web server and listens on port WEB_PORT (default 80)
func (app *Config) serve() {

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + WEB_PORT,
		Handler: app.routes(),
	}

	// Start server
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

// connectToMongo establishes a connection to a MongoDB instance
func connectToMongo() (*mongo.Client, error) {

	// Set options for connecting to Mongo
	clientOptions := options.Client().ApplyURI(MONGO_URL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// Connect
	conn, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting:", err)
		return nil, err
	}

	return conn, nil
}
