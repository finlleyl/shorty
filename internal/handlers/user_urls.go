package handlers

import (
	"encoding/json"
	"github.com/finlleyl/shorty/internal/app"
	"github.com/finlleyl/shorty/internal/auth"
	"github.com/finlleyl/shorty/internal/config"
	"net/http"
)

func UserURLsHandler(store app.Store, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := auth.GetUserIDFromContext(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		urls, err := store.GetByUserID(userID)
		if err != nil {
			http.Error(w, "Could not get URLs", http.StatusInternalServerError)
			return
		}

		if len(urls) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		for i := range urls {
			urls[i].ShortURL = cfg.B.BaseURL + "/" + urls[i].ShortURL
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(urls)
	}
}
