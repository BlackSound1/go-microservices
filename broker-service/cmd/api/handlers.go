package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

// Broker handles the broker service, returning a simple JSON message
func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {

	// Create a JSONResponse
	payload := JSONResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

// HandleSubmission handles any incoming requests to the broker service
//
// It expects to receive a JSON payload with an "action" parameter, which
// determines the action to take. The following actions are currently
// supported:
//
// - "auth": Authenticate the user using the credentials provided in the request body
//
// Any other action will result in an error response being sent
func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {

	var requestPayload RequestPayload

	// Read the JSON from the request
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Different behaviour depending on the action specified
	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "log":
		app.logItem(w, requestPayload.Log)
	case "mail":
		app.sendMail(w, requestPayload.Mail)
	default:
		app.errorJSON(w, errors.New("unknown action"))
	}
}

// sendMail sends an email by forwarding the MailPayload to the mail service.
func (app *Config) sendMail(w http.ResponseWriter, msg MailPayload) {

	jsonData, _ := json.MarshalIndent(msg, "", "\t")

	mailServiceURL := "http://mail-service/send"

	req, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		app.errorJSON(w, err)
		return
	}

	var payload JSONResponse
	payload.Error = false
	payload.Message = "sent mail to " + msg.To

	app.writeJSON(w, http.StatusAccepted, payload)
}

// logItem sends a request to the log service to log the given entry
func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {

	// Convert given entry into JSON
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log"

	// Create a new request to the log service
	req, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Set the necessary header
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling log service"))
		return
	}

	// Create a JSON response
	var payload JSONResponse
	payload.Error = false
	payload.Message = "logged"

	// Write the response
	app.writeJSON(w, http.StatusAccepted, payload)
}

// authenticate sends a request to the auth service to verify the user's credentials.
func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {

	// Read the JSON
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	// Create a new custom request to the auth service
	request, err := http.NewRequest("POST", "http://auth-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Actually perform the request by creating a client to do so
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// Make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling auth service"))
		return
	}

	// Read response body
	var jsonFromService JSONResponse
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// If we get an error in the response
	if jsonFromService.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	// We have valid login, so write a proper response
	var payload JSONResponse
	payload.Error = false
	payload.Message = "Successfully authenticated"
	payload.Data = jsonFromService.Data
	app.writeJSON(w, http.StatusAccepted, payload)
}
