package handlers

import (
	"encoding/json"
	"github.com/finlleyl/shorty/internal/app"
	"net/http"
)

type BatchRequest struct {
	CorrelationId string `json:"correlation_id"`
	OriginalUrl   string `json:"original_url"`
}

func BatchHandler(store app.Store) http.HandlerFunc {
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
			store.Save(id, request.OriginalUrl)
			response = append(response, map[string]string{"correlation_id": request.CorrelationId, "short_url": id})

		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}
