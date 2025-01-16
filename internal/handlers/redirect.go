package handlers

import (
	"github.com/finlleyl/shorty/internal/app"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func RedirectHandler(storage *app.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		origURL, exists := storage.Get(id)
		if !exists {
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Location", origURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
