package handlers

import (
	"encoding/json"
	"github.com/finlleyl/shorty/internal/app"
	"github.com/finlleyl/shorty/internal/auth"
	"io"
	"net/http"
	"time"
)

type DeleteTask struct {
	UserID   string
	ShortURL string
}

type DeleteHandler struct {
	Store        app.Store
	DeleteTaskCh chan<- DeleteTask
}

func (h DeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	for _, url := range urls {
		h.DeleteTaskCh <- DeleteTask{
			UserID:   userID,
			ShortURL: url,
		}
	}

	w.WriteHeader(http.StatusAccepted)
}

func BatchDeleteWorker(deleteTaskCh <-chan DeleteTask, store app.Store) {
	tasks := make(map[string]map[string]struct{})
	flushInterval := 50 * time.Millisecond
	timer := time.NewTimer(flushInterval)
	defer timer.Stop()

	for {
		select {
		case task := <-deleteTaskCh:
			if tasks[task.UserID] == nil {
				tasks[task.UserID] = make(map[string]struct{})
			}
			tasks[task.UserID][task.ShortURL] = struct{}{}
		case <-timer.C:
			for userID, urlsSet := range tasks {
				var urls []string
				for url := range urlsSet {
					urls = append(urls, url)
				}
				if len(urls) > 0 {
					if err := store.BatchDelete(urls, userID); err != nil {
					}
				}
			}
			tasks = make(map[string]map[string]struct{})
			timer.Reset(flushInterval)
		}
	}
}
