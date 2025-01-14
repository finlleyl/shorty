package handlers

import (
	"github.com/finlleyl/shorty/internal/app"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRedirectHandler(t *testing.T) {

	tests := []struct {
		name           string
		storageSetup   func() (*app.Storage, string)
		requestPath    string
		expectedStatus int
		expectedHeader string
	}{
		{
			name: "positive test",
			storageSetup: func() (*app.Storage, string) {
				storage := app.NewStorage()
				id := storage.Save("google.com")
				return storage, id
			},
			requestPath:    "",
			expectedStatus: http.StatusTemporaryRedirect,
			expectedHeader: "google.com",
		},
		{
			name: "negative test",
			storageSetup: func() (*app.Storage, string) {
				storage := app.NewStorage()
				id := storage.Save("google.com")
				return storage, id
			},
			requestPath:    "/123",
			expectedStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage, id := tt.storageSetup()
			if tt.requestPath == "" {
				tt.requestPath = "/" + id
			}
			request := httptest.NewRequest(http.MethodGet, tt.requestPath, nil)
			recorder := httptest.NewRecorder()

			handler := RedirectHandler(storage)
			handler.ServeHTTP(recorder, request)

			result := recorder.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.expectedStatus, result.StatusCode)
			if tt.expectedHeader != "" {
				assert.Equal(t, tt.expectedHeader, result.Header.Get("Location"))
			}

		})
	}
}
