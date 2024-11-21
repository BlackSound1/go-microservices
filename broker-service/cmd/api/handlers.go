package main

import (
	"encoding/json"
	"net/http"
)

type JSONResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Broker handles the broker service, returning a simple JSON message
func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {

	// Create a JSONResponse
	payload := JSONResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	// Convert to JSON, indentend nicely
	out, _ := json.MarshalIndent(payload, "", "\t")

	// Set the headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	// Write the response to the response writer
	w.Write(out)
}
