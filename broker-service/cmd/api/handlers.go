package main

import (
	"net/http"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Broker service is running",
	}

	err := app.writeJSON(w, http.StatusOK, payload)

	if err != nil {
		app.errorJSON(w, err)
	}
}
