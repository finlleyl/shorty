package handlers

import (
	"github.com/finlleyl/shorty/internal/app"
	"net/http"
)

func ShortenHandler(storage *app.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		longURL := r.FormValue("url")
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
