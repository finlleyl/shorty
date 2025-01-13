package handlers

import (
	"github.com/finlleyl/shorty/internal/app"
	"io"
	"net/http"
	"strconv"
)

func ShortenHandler(storage *app.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil || len(body) == 0 {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		url := string(body)
		id := storage.Save(url)
		shortURL := "http://localhost:8080/" + id

		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", strconv.Itoa(len(shortURL)))
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(shortURL))
	}
}
