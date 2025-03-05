package handlers

import (
	"encoding/json"
	"github.com/finlleyl/shorty/internal/app"
	"github.com/finlleyl/shorty/internal/auth"
	"io"
	"net/http"
)

func DeleteHandler(store app.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return
		}
		defer r.Body.Close()

		var urls []string
		if err = json.Unmarshal(body, &urls); err != nil {
			return
		}

		userID, ok := auth.GetUserIDFromContext(r)
		if !ok {
			return
		}
		store.BatchDelete(urls, userID)

		w.WriteHeader(http.StatusAccepted)

	}
}
