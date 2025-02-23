package main

import (
	"compress/gzip"
	"net/http"
	"strings"
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

func gzipMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			gr, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Failed to create gzip reader", http.StatusInternalServerError)
				return
			}
			defer gr.Close()
			r.Body = gr
		}

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Add("Vary", "Accept-Encoding")

		gzWriter := gzip.NewWriter(w)
		defer gzWriter.Close()

		grw := &gzipResponseWriter{
			ResponseWriter: w,
			gzipWriter:     gzWriter,
			headersSent:    false,
		}

		next(grw, r)
	}
}
