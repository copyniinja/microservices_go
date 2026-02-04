package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (a *Config) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization", "X-CSRF-TOKEN"},
		AllowCredentials: true,
	}))
	r.Use(middleware.Heartbeat("/health"))
	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello from logger server")
	})
	r.Post("/log", a.WriteLog)

	return r

}
