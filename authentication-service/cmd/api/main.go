package main

import (
	"authentication/data"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const PORT = "5000"

var count uint64

var dbpool *pgxpool.Pool

type Config struct {
	DB     *pgxpool.Pool
	Models data.Models
}

func main() {
	fmt.Println("Starting authentication service on port " + PORT)

	// connect to database
	connection := connectToDB()

	if connection == nil {
		log.Panicln("Database connection failed")
	}

	// start application
	app := Config{
		DB:     connection,
		Models: data.New(connection),
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: app.routes(),
	}

	// start server
	err := server.ListenAndServe()

	if err != nil {
		log.Panicln(err)
	}

	defer dbpool.Close()
}

func openDB(config *pgxpool.Config) (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.NewWithConfig(context.Background(), config)

	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
		return nil, err
	}

	// Test the connection
	err = dbpool.Ping(context.Background())

	if err != nil {
		fmt.Printf("Unable to ping database: %v\n", err)
		return nil, err
	}

	return dbpool, nil
}

func connectToDB() *pgxpool.Pool {
	dsn := os.Getenv("DSN")

	// confirm if the connection string is valid
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		fmt.Printf("Unable to parse DSN string: %v\n", err)
		return nil
	}

	for {
		connection, err := openDB(config)

		if err != nil {
			fmt.Println("PostgreSQL is not ready yet...")
			count++
		} else {
			fmt.Println("Database connect successfully...")
			return connection
		}

		if count > 10 {
			fmt.Println(err)
			return nil
		}

		fmt.Println("Restart the database connection after 2 seconds")
		time.Sleep(2 * time.Second)

		continue
	}

}
