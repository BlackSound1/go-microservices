package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type JSONResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"` // Omit this field if empty
}

// readJSON reads a JSON from a request body into the target data.
//
// It checks that the body is not more than 1 MB and that it only contains 1 JSON value.
func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data any) error {

	// Make sure JSON file is less than 1 MB. Shouldn't be more than that
	maxBytes := int64(1048576) // 1 MB
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

	// Create decoder for JSON
	decoder := json.NewDecoder(r.Body)

	// Decode JSON
	err := decoder.Decode(data)
	if err != nil {
		return err
	}

	// Make sure there is only 1 JSON value
	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

// writeJSON sends a JSON response with the given status code and data.
//
// If any headers are passed in, they are added to the response.
func (app *Config) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {

	// Convert data to JSON
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// If there are any headers, add them to the response
	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	// Set the JSON header
	w.Header().Set("Content-Type", "application/json")

	// Write the status code to the response
	w.WriteHeader(status)

	// Write the response
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

// errorJSON sends a JSON response with an error message and optional status code.
// If a status code is provided, it uses the first one; otherwise, it defaults to http.StatusBadRequest.
func (app *Config) errorJSON(w http.ResponseWriter, err error, status ...int) error {

	// Set a status code for default
	statusCode := http.StatusBadRequest

	// If there are any status codes, use the first one
	if len(status) > 0 {
		statusCode = status[0]
	}

	// Set up the payload
	var payload JSONResponse
	payload.Error = true
	payload.Message = err.Error()

	// Send the response
	return app.writeJSON(w, statusCode, payload)
}
