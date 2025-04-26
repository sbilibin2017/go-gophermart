package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Decompress request body if it's gzipped
		if r.Header.Get("Content-Encoding") == "gzip" {
			gzipReader, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Failed to read gzip data", http.StatusBadRequest)
				return
			}
			defer gzipReader.Close()
			r.Body = io.NopCloser(gzipReader)
		}

		// Compress response body if client accepts gzip
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			w.Header().Set("Content-Encoding", "gzip")
			gzipWriter := gzip.NewWriter(w)
			defer gzipWriter.Close()

			grw := &gzipResponseWriter{
				ResponseWriter: w,
				Writer:         gzipWriter,
			}
			next.ServeHTTP(grw, r)
			return
		}

		// Serve without gzip compression
		next.ServeHTTP(w, r)
	})
}

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w *gzipResponseWriter) Write(p []byte) (int, error) {
	return w.Writer.Write(p)
}
