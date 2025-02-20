package handlers

import (
	"github.com/finlleyl/shorty/db"
	"net/http"
)

func CheckConnectionHandler(w http.ResponseWriter, r *http.Request) {
	if err := db.PingDB; err != nil {
		http.Error(w, "Database is not reachable", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Database is healthy"))
}
