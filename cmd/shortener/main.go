package main

import (
	"io"
	"math/rand"
	"net/http"
)

type Storage map[string]string

var storage Storage

func (s *Storage) Save(url string) string {
	id := generateID()

	(*s)[id] = url

	return id
}

func (s *Storage) Get(id string) (string, bool) {
	origURL, exists := (*s)[id]
	return origURL, exists
}

func generateID() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}

func generateIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	url := string(body)
	id := storage.Save(url)
	shortURL := "http://localhost:8080/" + id

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

func main() {
	mux := http.NewServeMux()
	storage = make(Storage)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		generateIDHandler(w, r)
	})

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		return
	}
}
