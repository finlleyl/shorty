package handlers

import (
	"github.com/finlleyl/shorty/internal/app"
	"github.com/finlleyl/shorty/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestJSONHandler(t *testing.T) {
	router := chi.NewRouter()
	storage := app.NewStorage()
	cfg := &config.Config{
		A: config.AddressConfig{
			Host:    "localhost",
			Port:    8080,
			Address: "localhost:8080",
		},
		B: config.BaseURLConfig{
			BaseURL: "http://localhost:8080/",
		},
	}
	router.Post("/api/shorten", JSONHandler(storage, cfg))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(`{"url":"http://google.com"}`))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)
	result := recorder.Result()
	defer result.Body.Close()

	assert.Equal(t, http.StatusCreated, result.StatusCode)
	_, err := io.ReadAll(result.Body)
	assert.NoError(t, err)

}
