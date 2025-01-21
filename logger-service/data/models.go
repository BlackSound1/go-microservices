package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type Models struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// New creates a Models instance with a given MongoDB client.
func New(mongo *mongo.Client) Models {
	client = mongo

	return Models{
		LogEntry: LogEntry{},
	}
}

// Insert adds a new log entry to the database.
func (l *LogEntry) Insert(entry LogEntry) error {

	// Create a database called "logs" and a collection in that database called "logs".
	// Will just use it if it exists and create it if it doesn't.
	collection := client.Database("logs").Collection("logs")

	// Insert a log entry into the collection
	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("Error inserting logs:", err)
		return err
	}

	return nil
}

// GetAll gets all log entries from the database, sorted by creation date in descending order.
func (l *LogEntry) GetAll() ([]*LogEntry, error) {

	// Create a timeout to prevent long execution
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	// Set some options to sort the logs by creation date (ASC)
	opts := options.Find()
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	// Execute find command on the logs collections
	cursor, err := collection.Find(context.TODO(), bson.D{{}}, opts)
	if err != nil {
		log.Println("Finding all docs error: ", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	// Start creating output
	var logs []*LogEntry

	// Iterate through cursor containing all results
	for cursor.Next(ctx) {

		// Decode each log entry
		var entry LogEntry

		err := cursor.Decode(&entry)
		if err != nil {
			log.Println("Error decoding log entry: ", err)
			return nil, err
		}

		// Append each log entry
		logs = append(logs, &entry)
	}

	return logs, nil
}

// GetOne gets a log entry by ID from the database.
func (l *LogEntry) GetOne(id string) (*LogEntry, error) {

	// Create a timeout to prevent long execution
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	// Convert given id into an ObjectID
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Start creating output
	var entry LogEntry

	// Find the log entry by id
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

// DropCollection drops the given collection in the given database. It
// is a potentially destructive method and should only be used during
// development.
func (l *LogEntry) DropCollection(db, collection string) error {

	// Create a timeout to prevent long execution
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Grab the requested collection
	coll := client.Database(db).Collection(collection)

	// Try to drop it
	if err := coll.Drop(ctx); err != nil {
		return err
	}

	return nil
}

// Update modifies an existing log entry in the database by its ID.
func (l *LogEntry) Update() (*mongo.UpdateResult, error) {

	// Create a timeout to prevent long execution
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	// Convert the id of this LogEntry into an ObjectID
	docID, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil {
		return nil, err
	}

	// Update this LogEntry
	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": docID},
		bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "name", Value: l.Name},
				{Key: "data", Value: l.Data},
				{Key: "updated_at", Value: time.Now()},
			}},
		},
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}
