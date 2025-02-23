package main

import (
	"github.com/finlleyl/shorty/db"
	"github.com/finlleyl/shorty/internal/app"
	"github.com/finlleyl/shorty/internal/config"
	"github.com/finlleyl/shorty/internal/handlers"
	"github.com/finlleyl/shorty/internal/logger"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	cfg := config.ParseFlags()

	var store app.Store

	if cfg.D.Address != "" {
		db.InitDB(cfg)
		store = db.NewPostgresStore(db.DB)
	} else {
		storage := app.NewStorage(cfg.F.Path)
		store = storage
	}

	logInstance, err := logger.InitializeLogger()
	if err != nil {
		return
	}

	defer logInstance.Sync()

	r := chi.NewRouter()

	r.Post("/", logger.WithLogging(gzipMiddleware(handlers.ShortenHandler(store, cfg))))
	r.Get("/{id}", logger.WithLogging(gzipMiddleware(handlers.RedirectHandler(store))))
	r.Post("/api/shorten", logger.WithLogging(gzipMiddleware(handlers.JSONHandler(store, cfg))))
	r.Get("/ping", logger.WithLogging(handlers.CheckConnectionHandler))
	r.Post("/api/shorten/batch", logger.WithLogging(gzipMiddleware(handlers.BatchHandler(store, cfg))))

	logger.Sugar.Infow("Server started", "address", cfg.A.Address)
	if err := http.ListenAndServe(cfg.A.Address, r); err != nil {
		logger.Sugar.Fatalw("Server failed", "error", err)
	}
}
