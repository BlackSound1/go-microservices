package main

import (
	"net/http"

	"github.com/BlackSound1/go-microservices/logger/data"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// WriteLog handles writing a log entry to the database.
func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {

	var requestPayload JSONPayload

	// Read the request JSON
	_ = app.readJSON(w, r, &requestPayload)

	// Create a new LogEntry for this payload
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	// Insert the log entry
	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Create response
	response := JSONResponse{
		Error:   false,
		Message: "logged",
	}

	// Write the response
	app.writeJSON(w, http.StatusAccepted, response)
}
