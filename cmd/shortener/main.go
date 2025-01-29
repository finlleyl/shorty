package main

import (
	"github.com/finlleyl/shorty/internal/app"
	"github.com/finlleyl/shorty/internal/config"
	"github.com/finlleyl/shorty/internal/handlers"
	"github.com/finlleyl/shorty/internal/logger"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	storage := app.NewStorage()
	cfg := config.ParseFlags()

	logInstance, err := logger.InitializeLogger()
	if err != nil {
		return
	}
	defer logInstance.Sync()

	r := chi.NewRouter()

	r.Post("/", logger.WithLogging(handlers.ShortenHandler(storage, cfg)))
	r.Get("/{id}", logger.WithLogging(handlers.RedirectHandler(storage)))
	r.Post("/api/shorten", logger.WithLogging(handlers.JSONHandler(storage, cfg)))

	logger.Sugar.Infow("Server started", "address", cfg.A.Address)
	if err := http.ListenAndServe(cfg.A.Address, r); err != nil {
		logger.Sugar.Fatalw("Server failed", "error", err)
	}
}
