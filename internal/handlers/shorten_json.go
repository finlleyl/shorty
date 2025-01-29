package handlers

import (
	"encoding/json"
	"github.com/finlleyl/shorty/internal/app"
	"github.com/finlleyl/shorty/internal/config"
	"github.com/finlleyl/shorty/internal/models"
	"net/http"
)

func JSONHandler(storage *app.Storage, config *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.Request
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		shortURL := storage.Save(req.URL)
		if shortURL == "" {
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
