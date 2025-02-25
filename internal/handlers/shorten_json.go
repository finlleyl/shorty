package handlers

import (
	"encoding/json"
	"errors"
	"github.com/finlleyl/shorty/db"
	"github.com/finlleyl/shorty/internal/app"
	"github.com/finlleyl/shorty/internal/config"
	"github.com/finlleyl/shorty/internal/models"
	"net/http"
)

func JSONHandler(store app.Store, config *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.Request
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		shortURL := app.GenerateID()
		if shortURL == "" {
			http.Error(w, "could not save URL", http.StatusInternalServerError)
			return
		}

		_, err := store.Save(shortURL, req.URL)
		if err != nil {
			if errors.Is(err, db.ErrConflict) {
				http.Error(w, "URL already exists", http.StatusConflict)
				return
			}

			http.Error(w, "could not save URL", http.StatusInternalServerError)

			return
		}

		resp := models.Response{
			Result: config.B.BaseURL + "/" + shortURL,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		enc := json.NewEncoder(w)
		if err := enc.Encode(resp); err != nil {
			return
		}
	}
}
