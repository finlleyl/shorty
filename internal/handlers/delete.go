package handlers

import (
	"encoding/json"
	"github.com/finlleyl/shorty/internal/app"
	"github.com/finlleyl/shorty/internal/auth"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	deleteTasksMu sync.Mutex
	// Для каждого userID накапливаем множество short URL для удаления
	deleteTasks = make(map[string]map[string]struct{})
)

func enqueueDeletion(userID string, urls []string) {
	deleteTasksMu.Lock()
	defer deleteTasksMu.Unlock()
	if deleteTasks[userID] == nil {
		deleteTasks[userID] = make(map[string]struct{})
	}
	for _, url := range urls {
		deleteTasks[userID][url] = struct{}{}
	}
}

func flushDeletions(store app.Store) {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		deleteTasksMu.Lock()
		for userID, urlsSet := range deleteTasks {
			var urls []string
			for url := range urlsSet {
				urls = append(urls, url)
			}
			if len(urls) > 0 {
				if err := store.BatchDelete(urls, userID); err != nil {
					log.Printf("Batch delete error for user %s: %v", userID, err)
				} else {
					log.Printf("Deleted URLs for user %s: %v", userID, urls)
				}
			}
		}
		deleteTasks = make(map[string]map[string]struct{})
		deleteTasksMu.Unlock()
	}
}

func DeleteHandler(store app.Store) http.HandlerFunc {
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

		userID, ok := auth.GetUserIDFromContext(r)
		if !ok || userID == "" {
			http.Error(w, "User not authenticated", http.StatusUnauthorized)
			return
		}

		go flushDeletions(store)

		enqueueDeletion(userID, urls)

		w.WriteHeader(http.StatusAccepted)
	}
}
