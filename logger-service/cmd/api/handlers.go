package main

import (
	"fmt"
	"logger/data"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	// read request body (json) into a variable
	var requestPayload JSONPayload
	_ = app.readJSON(w, r, &requestPayload)

	// insert data into the logger db
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := app.Models.LogEntry.Insert(event)

	if err != nil {
		app.errorJSON(w, err)
		return
	}
	// build success response
	log := jsonResponse{
		Error:   false,
		Message: "Logged ✅",
		Data:    requestPayload,
	}

	app.writeJSON(w, http.StatusAccepted, log)
}

func (app *Config) GetLog(w http.ResponseWriter, r *http.Request) {
	// read part param from the URL
	logId := chi.URLParam(r, "logId")
	log, err := app.Models.LogEntry.GetOne(logId)

	if err != nil {
		app.errorJSON(w, err)
		return
	}
	// build success response
	response := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Fetch log with ID %s", logId),
		Data:    log,
	}

	app.writeJSON(w, http.StatusAccepted, response)
}

func (app *Config) GetLogs(w http.ResponseWriter, r *http.Request) {
	logs, err := app.Models.LogEntry.All()

	if err != nil {
		app.errorJSON(w, err)
		return
	}
	// build success response
	response := jsonResponse{
		Error:   false,
		Message: "Fetch all logged events 🆕",
		Data:    logs,
	}

	app.writeJSON(w, http.StatusAccepted, response)
}
func (app *Config) UpdateLog(w http.ResponseWriter, r *http.Request) {
	// read part param from the URL
	logId := chi.URLParam(r, "logId")

	// read update into a variable
	var logEntry data.LogEntry

	err := app.readJSON(w, r, &logEntry)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_, err = logEntry.Update()

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	res, err := logEntry.GetOne(logId)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// build success response
	response := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Update log with ID %s", logId),
		Data:    res,
	}

	app.writeJSON(w, http.StatusAccepted, response)
}
func (app *Config) DropLogs(w http.ResponseWriter, r *http.Request) {
	// clear the collection
	err := app.Models.LogEntry.DropCollection()

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	response := jsonResponse{
		Error:   false,
		Message: "Clear all the logged events 🙌",
	}

	app.writeJSON(w, http.StatusAccepted, response)
}
