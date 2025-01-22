package main

import (
	"github.com/finlleyl/shorty/internal/app"
	"github.com/finlleyl/shorty/internal/config"
	"github.com/finlleyl/shorty/internal/handlers"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	storage := app.NewStorage()
	cfg := config.ParseFlags()

	r := chi.NewRouter()

	r.Post("/", handlers.ShortenHandler(storage, cfg))
	r.Get("/{id}", handlers.RedirectHandler(storage))

	log.Fatal(http.ListenAndServe(cfg.A.Address, r))
}
