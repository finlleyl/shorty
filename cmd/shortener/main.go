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

	r.Post("/", logger.WithLogging(gzipMiddleware(handlers.ShortenHandler(storage, cfg))))
	r.Get("/{id}", logger.WithLogging(gzipMiddleware(handlers.RedirectHandler(storage))))
	r.Post("/api/shorten", logger.WithLogging(gzipMiddleware(handlers.JSONHandler(storage, cfg))))

	logger.Sugar.Infow("Server started", "address", cfg.A.Address)
	if err := http.ListenAndServe(cfg.A.Address, r); err != nil {
		logger.Sugar.Fatalw("Server failed", "error", err)
	}
}

func gzipMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// По умолчанию ответ пишем в ow = w
		ow := w

		// Определяем, хочет ли клиент gzip
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")

		// Если да, оборачиваем w в gzip.Writer
		if supportsGzip {
			w.Header().Set("Content-Encoding", "gzip")
			w.Header().Add("Vary", "Accept-Encoding")

			cw := newCompressWriter(w)
			ow = cw
			defer cw.Close()
		}

		// А теперь на уровне запроса (входящих данных) проверяем,
		// не пришел ли он к нам gzipped
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")

		// Если клиент отправил gzip и это данные типа JSON или HTML,
		// мы распаковываем Body (необязательная логика, зависит от API).
		if sendsGzip && (strings.Contains(r.Header.Get("Content-Type"), "application/json") ||
			strings.Contains(r.Header.Get("Content-Type"), "text/html")) {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				// Ошибка при создании gzip.Reader
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			// Переопределяем r.Body
			r.Body = cr
			defer cr.Close()
		}

		// И вызываем основной хендлер
		h.ServeHTTP(ow, r)
	}
}
