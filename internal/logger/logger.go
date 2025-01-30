package logger

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

var Sugar *zap.SugaredLogger

type (
	responseData struct {
		statusCode int
		size       int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.statusCode = statusCode
}

func InitializeLogger() (*zap.Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	Sugar = logger.Sugar()

	return logger, nil
}

func WithLogging(h http.HandlerFunc) http.HandlerFunc {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			statusCode: http.StatusOK,
			size:       0,
		}

		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		h.ServeHTTP(&lw, r)

		duration := time.Since(start)

		Sugar.Infow("HTTP Request",
			"uri", r.RequestURI,
			"method", r.Method,
			"status", responseData.statusCode,
			"duration", duration,
			"size", responseData.size,
		)

	}
	return logFn
}
