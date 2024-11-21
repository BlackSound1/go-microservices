package main

import (
	"net/http"
)

// Broker handles the broker service, returning a simple JSON message
func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {

	// Create a JSONResponse
	payload := JSONResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}
