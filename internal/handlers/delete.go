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

type DeleteTask struct {
	UserID   string
	ShortURL string
}

// DeleteHandler отправляет задачи удаления в канал.
type DeleteHandler struct {
	Store        app.Store
	DeleteTaskCh chan<- DeleteTask
}

func (h *DeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	// Для каждого URL отправляем задачу в канал.
	for _, url := range urls {
		select {
		case h.DeleteTaskCh <- DeleteTask{UserID: userID, ShortURL: url}:
		default:
			// Если канал заполнен, логируем и немедленно обрабатываем задачу
			log.Printf("Task channel is full, executing delete immediately for user %s, url %s", userID, url)
			if err := h.Store.BatchDelete([]string{url}, userID); err != nil {
				log.Printf("Immediate delete error: %v", err)
			}
		}
	}

	w.WriteHeader(http.StatusAccepted)
}

// BatchDeleteWorker агрегирует задачи и выполняет batch update каждые 10 мс.
func BatchDeleteWorker(deleteTaskCh <-chan DeleteTask, store app.Store) {
	tasks := make(map[string]map[string]struct{})
	var mu sync.Mutex
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case task := <-deleteTaskCh:
			mu.Lock()
			if tasks[task.UserID] == nil {
				tasks[task.UserID] = make(map[string]struct{})
			}
			tasks[task.UserID][task.ShortURL] = struct{}{}
			mu.Unlock()

		case <-ticker.C:
			mu.Lock()
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
			tasks = make(map[string]map[string]struct{})
			mu.Unlock()
		}
	}
}
