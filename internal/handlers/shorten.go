package handlers

import (
	"github.com/finlleyl/shorty/internal/app"
	"io"
	"net/http"
)

func ShortenHandler(storage *app.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		body, err := io.ReadAll(r.Body)
		if err != nil || len(body) == 0 {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		longURL := string(body)
		if longURL == "" {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		id := storage.Save(longURL)
		shortURL := "http://localhost:8080/" + id

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(shortURL))
	}
}
