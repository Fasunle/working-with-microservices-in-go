package data

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func New(mongo *mongo.Client) Models {
	client = mongo

	return Models{
		LogEntry: LogEntry{},
	}
}

type Models struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID        string    `bson:"_id,omitempty" "json:id,omitempty"`
	Name      string    `bson:"name,omitempty" "json:name,omitempty"`
	Data      string    `bson:"data,omitempty" "json:data,omitempty"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func (l *LogEntry) Insert(entry LogEntry) error {
	//
	MONGODB_NAME := os.Getenv("MONGODB_NAME")
	collection := client.Database(MONGODB_NAME).Collection("logs")

	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		fmt.Println("Error inserting into logs:", err)
		return err
	}

	return nil
}

func (l *LogEntry) All() ([]*LogEntry, error) {
	//
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	MONGODB_NAME := os.Getenv("MONGODB_NAME")
	collection := client.Database(MONGODB_NAME).Collection("logs")

	opts := options.Find()

	opts.SetSort(bson.D{{"created_at", -1}})

	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)

	if err != nil {
		log.Println("Finding all logs error")
		return nil, err
	}

	defer cursor.Close(ctx)

	var logs []*LogEntry

	for cursor.Next(ctx) {
		var log LogEntry

		err := cursor.Decode(&log)

		if err != nil {
			fmt.Println("Error decoding log into logs slice")
			return nil, err
		}

		logs = append(logs, &log)
	}

	return logs, nil

}

func (l *LogEntry) GetOne(id string) (*LogEntry, error) {
	// create context for database operation
	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)
	defer cancel()

	// grab a referrence to the database collection
	MONGODB_NAME := os.Getenv("MONGODB_NAME")
	collection := client.Database(MONGODB_NAME).Collection("logs")

	// convert document strig ID to mongo objectID
	docID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		log.Print("Document ID is invalid: ", id)
		return nil, err
	}

	// find the document from the collection
	var entry LogEntry
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry)

	if err != nil {
		log.Println("Error decoding log entry")
		return nil, err
	}

	return &entry, nil
}

func (l *LogEntry) DropCollection() error {
	// create context so that the operation is cancelled if taking longer time
	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)
	defer cancel()

	MONGODB_NAME := os.Getenv("MONGODB_NAME")
	collection := client.Database(MONGODB_NAME).Collection("logs")

	if err := collection.Drop(ctx); err != nil {
		return err
	}

	return nil
}

func (l *LogEntry) Update() (*mongo.UpdateResult, error) {
	// create context
	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)
	defer cancel()

	MONGODB_NAME := os.Getenv("MONGODB_NAME")
	collection := client.Database(MONGODB_NAME).Collection("logs")

	// validate the document ID
	ID, err := primitive.ObjectIDFromHex(l.ID)

	if err != nil {
		log.Print("invalid document ID")
		return nil, err
	}

	// update log
	return collection.UpdateOne(
		ctx,
		bson.M{"_id": ID},
		bson.D{
			{"$set", bson.D{
				{"name", l.Name},
				{"data", l.Data},
				{"updated_at", time.Now()},
			}},
		},
	)
}
