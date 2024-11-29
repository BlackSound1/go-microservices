package main

import "net/http"

// SendMail handles sending an email by reading the request payload,
// constructing a Message object, and sending it via the Mailer.
func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {

	// mailMessage struct holds the email data received in the request
	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	// Read the JSON request body into requestPayload
	var requestPayload mailMessage
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Create a Message object from the request payload
	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	// Send the email using the Mailer
	err = app.Mailer.SendSMTPMessage(msg)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Set up a response payload
	payload := JSONResponse{
		Error:   false,
		Message: "sent mail to " + requestPayload.To,
	}

	// Send a success JSON response
	app.writeJSON(w, http.StatusAccepted, payload)
}
