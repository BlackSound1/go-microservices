package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	// Specify who can connect
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},                                   // Anyone can connect to this service
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},                 // Many methods are allowed
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"}, // These headers are allowed
		ExposedHeaders:   []string{"Link"},                                                    // This header is exposed
		AllowCredentials: true,                                                                // Cookies, other credentials are allowed
		MaxAge:           300,                                                                 // Cache for 5 minutes
	}))

	mux.Use(middleware.Heartbeat("/ping")) // Health check

	// Set up handlers
	mux.Post("/log", app.WriteLog)

	return mux
}
