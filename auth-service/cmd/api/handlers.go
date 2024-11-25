package main

import (
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

	// Create response to send back
	payload := JSONResponse{
		Error: false,
		Message: "Logged in user " + user.Email + " successfully",
		Data: user,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}
