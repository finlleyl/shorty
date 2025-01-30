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

func gzipMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Проверить, нужно ли декомпрессировать входящий запрос (Content-Encoding: gzip)
		// 2. Проверить, нужно ли сжимать ответ (Accept-Encoding: gzip),
		//    если нет — просто вызывать next(w, r) без обёртки.
		// 3. Если да, обернуть в gzipResponseWriter только BODY, а заголовки
		//    (Location и т.п.) не трогать.

		// Пример минимальной логики:

		// Прежде: распаковать входящий request, если у него Content-Encoding: gzip
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			// Имеет смысл распаковывать только при Content-Type: "application/json"
			// или "text/html", но если хотите — можно не делать такой фильтр.
			gr, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Failed to create gzip reader", http.StatusInternalServerError)
				return
			}
			defer gr.Close()
			r.Body = gr
		}

		// Проверяем, поддерживает ли клиент gzip
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// Если нет, просто вызываем следующий хендлер без сжатия
			next(w, r)
			return
		}

		// Устанавливаем заголовки, чтобы клиент понял, что дальше будет gzip
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Add("Vary", "Accept-Encoding")

		gzWriter := gzip.NewWriter(w)
		defer gzWriter.Close()

		// Оборачиваем http.ResponseWriter в свою структуру
		// (Она будет сжимать только тело, а сами заголовки не трогать)
		grw := &gzipResponseWriter{
			ResponseWriter: w,
			gzipWriter:     gzWriter,
			headersSent:    false,
		}

		// Вызываем основной хендлер, передавая ему «обёрнутый» writer
		next(grw, r)
	}
}
