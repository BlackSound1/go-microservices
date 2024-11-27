package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

// Authenticate validates a user's credentials
func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {

	// A request should look like this
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Read the request and save it into the payload
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Validate user
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	// Validate password
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	// Log auth request
	err = app.logRequest("auth", user.Email+" logged in")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Create response to send back
	payload := JSONResponse{
		Error:   false,
		Message: "Logged in user " + user.Email + " successfully",
		Data:    user,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

// logRequest sends a log entry with the specified name and data to the logger service.
func (app *Config) logRequest(name, data string) error {

	// Create object for given entry data
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}
	entry.Name = name
	entry.Data = data

	// Convert it to JSON
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	// Create the request
	logServiceURL := "http://logger-service/log"
	req, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// Do the request
	client := http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
