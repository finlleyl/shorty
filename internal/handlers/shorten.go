package handlers

import (
	"errors"
	"github.com/finlleyl/shorty/internal/app"
	"github.com/finlleyl/shorty/internal/apperrors"
	"github.com/finlleyl/shorty/internal/auth"
	"github.com/finlleyl/shorty/internal/config"
	"io"
	"net/http"
)

func ShortenHandler(store app.Store, config *config.Config) http.HandlerFunc {
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
		userID, ok := auth.GetUserIDFromContext(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		_, err = store.Save(id, longURL, userID)

		if err != nil {
			var conflictErr *apperrors.ConflictError

			if errors.As(err, &conflictErr) {
				w.WriteHeader(http.StatusConflict)
				_, _ = w.Write([]byte(config.B.BaseURL + "/" + conflictErr.ShortURL))
				return
			}
			http.Error(w, "Could not save URL", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(shortURL))
	}
}
