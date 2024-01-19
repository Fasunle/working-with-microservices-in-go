package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const PORT = "80"

var count uint64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	fmt.Println("Starting authentication service on port " + PORT)

	// connect to database
	connnection := connectToDB()

	if connnection == nil {
		log.Panicln("Database connection failed")
	}

	// start application
	app := Config{
		DB:     connnection,
		Models: data.New(connnection),
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
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)

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
