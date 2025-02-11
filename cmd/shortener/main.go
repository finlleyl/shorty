package main

import (
	"compress/gzip"
	"github.com/finlleyl/shorty/internal/app"
	"github.com/finlleyl/shorty/internal/config"
	"github.com/finlleyl/shorty/internal/handlers"
	"github.com/finlleyl/shorty/internal/logger"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
)

func main() {
	cfg := config.ParseFlags()
	storage := app.NewStorage(cfg.F.Path)

	logInstance, err := logger.InitializeLogger()
	if err != nil {
		return
	}
	defer logInstance.Sync()

	r := chi.NewRouter()

	r.Post("/", logger.WithLogging(gzipMiddleware(handlers.ShortenHandler(storage, cfg))))
	r.Get("/{id}", logger.WithLogging(gzipMiddleware(handlers.RedirectHandler(storage))))
	r.Post("/api/shorten", logger.WithLogging(gzipMiddleware(handlers.JSONHandler(storage, cfg))))

	logger.Sugar.Infow("Server started", "address", cfg.A.Address)
	if err := http.ListenAndServe(cfg.A.Address, r); err != nil {
		logger.Sugar.Fatalw("Server failed", "error", err)
	}
}

func gzipMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			gr, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Failed to create gzip reader", http.StatusInternalServerError)
				return
			}
			defer gr.Close()
			r.Body = gr
		}

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Add("Vary", "Accept-Encoding")

		gzWriter := gzip.NewWriter(w)
		defer gzWriter.Close()

		grw := &gzipResponseWriter{
			ResponseWriter: w,
			gzipWriter:     gzWriter,
			headersSent:    false,
		}

		next(grw, r)
	}
}
