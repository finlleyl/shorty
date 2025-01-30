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

	r.Post("/", gzipMiddleware(logger.WithLogging(handlers.ShortenHandler(storage, cfg))))
	r.Get("/{id}", gzipMiddleware(logger.WithLogging(handlers.RedirectHandler(storage))))
	r.Post("/api/shorten", gzipMiddleware(logger.WithLogging(handlers.JSONHandler(storage, cfg))))

	logger.Sugar.Infow("Server started", "address", cfg.A.Address)
	if err := http.ListenAndServe(cfg.A.Address, r); err != nil {
		logger.Sugar.Fatalw("Server failed", "error", err)
	}
}

func gzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		cw := newCompressWriter(w)
		defer cw.Close()

		next.ServeHTTP(cw, r)
	})
}
