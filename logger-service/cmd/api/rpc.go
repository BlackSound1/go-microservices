package main

import (
	"context"
	"log"
	"time"

	"github.com/BlackSound1/go-microservices/logger/data"
)

type RPCServer struct{}

type RPCPayload struct {
	Name string
	Data string
}

// LogInfo processes the given RPCPayload by inserting a log entry into the MongoDB
// database and returns a response message
func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {

	// Create a logs collection in Mongo
	collection := client.Database("logs").Collection("logs")

	// Insert a log entry
	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Println("error writing to Mongo", err)
		return err
	}

	// Set the response message
	*resp = "Processed payload via RPC: " + payload.Name
	
	return nil
}
