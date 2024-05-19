package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *Config) routes() *chi.Mux {
	mux := chi.NewRouter()

	// specify who has access to the route
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowCredentials: true,
		Debug:            true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	mux.Use(middleware.Heartbeat("/ping"))

	mux.Post("/log", app.WriteLog)
	mux.Get("/log", app.GetLogs)
	mux.Get("/log/:logId", app.GetLog)
	mux.Put("/log/:logId", app.UpdateLog)
	mux.Delete("/log", app.DropLogs)

	return mux
}
