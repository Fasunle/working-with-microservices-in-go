package main

import (
	"context"
	"fmt"
	"log"
	"logger/data"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var PORT = "80"

type Config struct {
	Models data.Models
}

var client *mongo.Client

func main() {

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

	app := Config{
		Models: data.New(client),
	}

	app.serve()

}

func connectToMongo() (*mongo.Client, error) {
	MONGODB_URI := os.Getenv("MONGODB_URI")
	MONGODB_USERNAME := os.Getenv("MONGODB_USERNAME")
	MONGODB_PASSWORD := os.Getenv("MONGODB_PASSWORD")
	// options.Client().ApplyURI("mongodb://foo:bar@localhost:27017")
	// client options
	clientOptions := options.Client().ApplyURI(MONGODB_URI)
	clientOptions.SetAuth(options.Credential{
		PasswordSet: true,
		Username:    MONGODB_USERNAME,
		Password:    MONGODB_PASSWORD,
	})

	clientOptions.SetConnectTimeout(15 * time.Second)
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)

	defer cancel()

	c, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		fmt.Println("Error connecting to mongodb:", err)
		return nil, err
	}

	err = c.Ping(context.Background(), nil)

	if err != nil {
		fmt.Println("Connection unsuccessful:", err)
		return nil, err
	}

	log.Println("Connected to mongo database 👌")
	return c, nil
}

func (app *Config) serve() {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: app.routes(),
	}

	server.RegisterOnShutdown(func() {
		log.Println("Shutting down the server")
	})

	err := server.ListenAndServe()

	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Logger service started!")
}
