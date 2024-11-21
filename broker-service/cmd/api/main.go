package main

import (
	"log"
	"net/http"
)

const WEB_PORT = "80"

type Config struct {
}


// main is the main entry point for the broker service.
// It starts the web server on the WEB_PORT constant port,
// and listens for incoming requests. It also sets up the
// routes defined in the Config struct.
func main() {
	app := Config{}

	log.Println("Starting broker service on port ", WEB_PORT)

	// Define HTTP server
	srv := &http.Server{
		Addr:    ":" + WEB_PORT,
		Handler: app.routes(),
	}

	// Start server
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
