package handlers

import (
	"github.com/finlleyl/shorty/internal/app"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestShortenHandler(t *testing.T) {
	tests := []struct {
		name         string
		storageSetup func() *app.Storage
		requestURL   string
		method       string
		expectedCode int
		expectedBody string
	}{
		{
			name: "Valid POST request",
			storageSetup: func() *app.Storage {
				return app.NewStorage()
			},
			requestURL:   "http://google.com",
			method:       http.MethodPost,
			expectedCode: http.StatusCreated,
			expectedBody: "http://localhost:8080/",
		},
		{
			name: "Empty request body",
			storageSetup: func() *app.Storage {
				return app.NewStorage()
			},
			requestURL:   "",
			method:       http.MethodPost,
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid request body\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := tt.storageSetup()
			formData := url.Values{}
			if tt.requestURL != "" {
				formData.Set("url", tt.requestURL)
			}

			r := chi.NewRouter()
			r.Post("/", ShortenHandler(storage))

			req := httptest.NewRequest(tt.method, "/", strings.NewReader(formData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			resp := rec.Result()
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					return
				}
			}(resp.Body)
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
