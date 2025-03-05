package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/finlleyl/shorty/internal/app"
	"github.com/finlleyl/shorty/internal/auth"
)

// DeleteTask описывает задачу удаления одного short URL для конкретного пользователя.
type DeleteTask struct {
	UserID   string
	ShortURL string
}

// DeleteHandler отправляет задачи удаления в общий канал.
func DeleteHandler(store app.Store, deleteTaskCh chan<- DeleteTask) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var urls []string
		if err = json.Unmarshal(body, &urls); err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		userID, ok := auth.GetUserIDFromContext(r)
		if !ok || userID == "" {
			http.Error(w, "User not authenticated", http.StatusUnauthorized)
			return
		}

		// Для каждого URL отправляем задачу в канал
		for _, url := range urls {
			deleteTaskCh <- DeleteTask{
				UserID:   userID,
				ShortURL: url,
			}
		}

		w.WriteHeader(http.StatusAccepted)
	}
}

// BatchDeleteWorker агрегирует задачи из канала и каждые flushInterval выполняет один batch‑update.
// Для каждого userID он собирает уникальные short URL и вызывает store.BatchDelete.
func BatchDeleteWorker(store app.Store, deleteTaskCh <-chan DeleteTask, flushInterval time.Duration) {
	// tasks группирует задачи по userID: для каждого пользователя хранится множество short URL.
	tasks := make(map[string]map[string]struct{})
	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

	for {
		select {
		case task := <-deleteTaskCh:
			if tasks[task.UserID] == nil {
				tasks[task.UserID] = make(map[string]struct{})
			}
			tasks[task.UserID][task.ShortURL] = struct{}{}
		case <-ticker.C:
			// По истечении flushInterval выполняем batch‑update для каждого пользователя
			for userID, urlsSet := range tasks {
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
			// Обнуляем накопленные задачи
			tasks = make(map[string]map[string]struct{})
		}
	}
}
