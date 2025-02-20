package handlers

import (
	"github.com/finlleyl/shorty/internal/app"
	"github.com/finlleyl/shorty/internal/config"
	"io"
	"net/http"
)

func ShortenHandler(storage *app.Storage, config *config.Config) http.HandlerFunc {
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

		id := app.GenerateID()
		shortURL := config.B.BaseURL + "/" + id

		storage.Save(id, longURL)

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(shortURL))
	}
}
