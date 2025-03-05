package handlers

import (
	"github.com/finlleyl/shorty/internal/app"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func RedirectHandler(store app.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		origURL, exists := store.Get(id)
		if !exists {
			if origURL == "alpha" {
				w.WriteHeader(http.StatusGone)
				return
			}
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Location", origURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
