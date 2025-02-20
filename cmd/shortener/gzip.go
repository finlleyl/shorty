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

func (g *gzipResponseWriter) WriteHeader(statusCode int) {
	if !g.headersSent {
		g.ResponseWriter.WriteHeader(statusCode)
		g.headersSent = true
	}
}

func (g *gzipResponseWriter) Write(b []byte) (int, error) {
	if !g.headersSent {
		g.WriteHeader(http.StatusOK)
	}
	return g.gzipWriter.Write(b)
}
