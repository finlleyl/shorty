package main

import (
	"compress/gzip"
	"net/http"
)

type gzipResponseWriter struct {
	http.ResponseWriter
	gzipWriter  *gzip.Writer
	headersSent bool
}

// Переопределяем WriteHeader, чтобы не слать его дважды
func (g *gzipResponseWriter) WriteHeader(statusCode int) {
	if !g.headersSent {
		g.ResponseWriter.WriteHeader(statusCode) // Отправляем «чистый» статус и заголовки
		g.headersSent = true
	}
}

// Переопределяем Write, чтобы писать в g.gzipWriter (тело) — оно будет в gzip
func (g *gzipResponseWriter) Write(b []byte) (int, error) {
	if !g.headersSent {
		// Если код не был отправлен явно, выставим по умолчанию 200
		g.WriteHeader(http.StatusOK)
	}
	return g.gzipWriter.Write(b) // записываем в gzip
}
