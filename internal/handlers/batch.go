package handlers

import (
	"encoding/json"
	"github.com/finlleyl/shorty/internal/app"
	"github.com/finlleyl/shorty/internal/auth"
	"github.com/finlleyl/shorty/internal/config"
	"net/http"
	"strings"
)

type BatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

func BatchHandler(store app.Store, config *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requests []BatchRequest

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&requests); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		defer r.Body.Close()

		response := make([]map[string]string, 0, len(requests))

		for _, request := range requests {
			id := app.GenerateID()
			originalURL := request.OriginalURL

			if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
				originalURL = "http://" + originalURL
			}

			userID, ok := auth.GetUserIDFromContext(r)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			store.Save(id, originalURL, userID)
			response = append(response, map[string]string{"correlation_id": request.CorrelationID, "short_url": config.B.BaseURL + "/" + id})

		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}
