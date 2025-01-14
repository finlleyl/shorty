package handlers

import (
	"github.com/finlleyl/shorty/internal/app"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestShortenHandler(t *testing.T) {
	tests := []struct {
		name         string
		storageSetup func() *app.Storage
		requestBody  string
		method       string
		expectedCode int
		expectedBody string
	}{
		{
			name: "Valid POST request",
			storageSetup: func() *app.Storage {
				return app.NewStorage()
			},
			requestBody:  "http://google.com",
			method:       http.MethodPost,
			expectedCode: http.StatusCreated,
			expectedBody: "http://localhost:8080/",
		},
		{
			name: "Invalid HTTP method",
			storageSetup: func() *app.Storage {
				return app.NewStorage()
			},
			requestBody:  "http://google.com",
			method:       http.MethodGet,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: "Method not allowed\n",
		},
		{
			name: "Empty request body",
			storageSetup: func() *app.Storage {
				return app.NewStorage()
			},
			requestBody:  "",
			method:       http.MethodPost,
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid request body\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := tt.storageSetup()

			req := httptest.NewRequest(tt.method, "/", strings.NewReader(tt.requestBody))
			rec := httptest.NewRecorder()

			handler := ShortenHandler(storage)
			handler.ServeHTTP(rec, req)

			resp := rec.Result()
			defer resp.Body.Close()
			if resp.StatusCode != tt.expectedCode {
				t.Errorf("unexpected status: got %v, want %v", resp.StatusCode, tt.expectedCode)
			}

			body, _ := io.ReadAll(resp.Body)
			if tt.expectedCode == http.StatusCreated {
				if !strings.HasPrefix(string(body), tt.expectedBody) {
					t.Errorf("unexpected body: got %v, want prefix %v", string(body), tt.expectedBody)
				}
			} else {
				if string(body) != tt.expectedBody {
					t.Errorf("unexpected body: got %v, want %v", string(body), tt.expectedBody)
				}
			}
		})
	}
}
