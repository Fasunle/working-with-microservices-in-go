package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var PORT = "80"

type Config struct{}

var client *mongo.Client

func main() {

	app := Config{}

	// connect to mongo
	mongoClient, err := connectToMongo()

	if err != nil {
		log.Panicln(err)
	}

	client = mongoClient

	// create context to disconnect the mongo connection when the application stops
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()

	// disconnect from database whenever the server stops
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Panic(err)
		}
	}()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: app.routes(),
	}

	err = server.ListenAndServe()

	if err != nil {
		log.Fatalln(err)
	}

}

func connectToMongo() (*mongo.Client, error) {
	MONGODB_URI := os.Getenv("MONGO_URI")
	MONGODB_USERNAME := os.Getenv("MONGODB_USERNAME")
	MONGODB_PASSWORD := os.Getenv("MONGODB_PASSWORD")

	// client options
	clientOptions := options.Client().ApplyURI(MONGODB_URI)
	clientOptions.SetAuth(options.Credential{
		Username: MONGODB_USERNAME,
		Password: MONGODB_PASSWORD,
	})

	c, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		fmt.Println("Error connecting to mongodb:", err)
		return nil, err
	}

	return c, nil
}
