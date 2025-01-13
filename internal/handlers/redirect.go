package handlers

import (
	"github.com/finlleyl/shorty/internal/app"
	"net/http"
)

func RedirectHandler(storage *app.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[1:]

		origURL, exists := storage.Get(id)
		if !exists {
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Location", origURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
