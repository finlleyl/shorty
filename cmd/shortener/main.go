package main

import (
	"github.com/finlleyl/shorty/internal/app"
	"github.com/finlleyl/shorty/internal/handlers"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	storage := app.NewStorage()

	r := chi.NewRouter()

	r.Route("", func(r chi.Router) {
		r.Post("/", handlers.ShortenHandler(storage))
		r.Get("/{id}", handlers.RedirectHandler(storage))
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}
