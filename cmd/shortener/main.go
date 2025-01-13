package main

import (
	"github.com/finlleyl/shorty/internal/app"
	"github.com/finlleyl/shorty/internal/handlers"
	"log"
	"net/http"
)

func main() {
	storage := app.NewStorage()

	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.ShortenHandler(storage))
	mux.HandleFunc("/{id}", handlers.RedirectHandler(storage))

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
