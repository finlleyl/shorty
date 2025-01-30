package main

import (
	"github.com/finlleyl/shorty/internal/app"
	"github.com/finlleyl/shorty/internal/config"
	"github.com/finlleyl/shorty/internal/handlers"
	"github.com/finlleyl/shorty/internal/logger"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
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

	r.Use(logger.WithLogging)
	r.Use(gzipMiddleware)

	r.Post("/", handlers.ShortenHandler(storage, cfg))
	r.Get("/{id}", handlers.RedirectHandler(storage))
	r.Post("/api/shorten", handlers.JSONHandler(storage, cfg))

	logger.Sugar.Infow("Server started", "address", cfg.A.Address)
	if err := http.ListenAndServe(cfg.A.Address, r); err != nil {
		logger.Sugar.Fatalw("Server failed", "error", err)
	}
}

func gzipMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			h.ServeHTTP(w, r)
			return
		}
		cw := newCompressWriter(w)
		defer cw.Close()

		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			if strings.Contains(r.Header.Get("Content-Type"), "application/json") ||
				strings.Contains(r.Header.Get("Content-Type"), "text/html") {
				cr, err := newCompressReader(r.Body)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				defer cr.Close()
				r.Body = cr
			}
		}

		h.ServeHTTP(cw, r)
	})
}
