package main

import (
	"fmt"
	"log"
	"net/http"
)

type Config struct{}

const PORT = "80"

func main() {
	app := Config{}

	fmt.Println("Starting broker service on port " + PORT)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: app.routes(),
	}

	// start server
	err := server.ListenAndServe()

	if err != nil {
		log.Fatalln(err)
	}

}
