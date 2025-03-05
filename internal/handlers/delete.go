package handlers

import (
	"encoding/json"
	"github.com/finlleyl/shorty/internal/app"
	"io"
	"log"
	"net/http"
)

func DeleteHandler(store *app.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var urls []string
		if err := json.Unmarshal(body, &urls); err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		userID, ok := r.Context().Value("userID").(string)
		if !ok || userID == "" {
			http.Error(w, "User not authenticated", http.StatusUnauthorized)
			return
		}

		go func(urls []string, userID string) {
			if err := store.BatchDelete(urls, userID); err != nil {
				log.Printf("Batch deletion error: %v", err)
			}
		}(urls, userID)

		w.WriteHeader(http.StatusAccepted)
	}
}
