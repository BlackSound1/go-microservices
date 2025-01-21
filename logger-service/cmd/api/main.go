package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
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

	// Disconnect whenever it is time to
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	// Register our custom RPC server type
	err = rpc.Register(new(RPCServer))
	if err != nil {
		log.Fatal("error registering RPC server: ", err)
	}

	// Start RPC server. Needs to be on separate goroutine because it is blocking
	go app.rpcListen()

	// Start gRPC server. Needs to be on separate goroutine because it is blocking
	go app.grpcListen()

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

// rpcListen starts an RPC server on the configured RPC_PORT (default 5001). It
// continually listens for incoming TCP connections and serves each connection on
// a separate goroutine
func (app *Config) rpcListen() error {

	log.Println("Starting RPC server on port: " + RPC_PORT)

	// Start listen for TCP connections on RPC_PORT
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", RPC_PORT))
	if err != nil {
		return err
	}
	defer listen.Close()

	// Continually try to accept connections. When connection is accepted,
	// serve the connection on a separate goroutine (because ServeConn is blocking)
	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}

		go rpc.ServeConn(rpcConn)
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
