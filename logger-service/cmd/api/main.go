package main

import (
	"fmt"
	"log"
	"net/http"
)

var PORT = "80"

type Config struct{}

func main() {

	app := Config{}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: app.routes(),
	}

	err := server.ListenAndServe()

	if err != nil {
		log.Fatalln(err)
	}

}
