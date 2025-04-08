package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"github.com/sbilibin2017/go-gophermart/internal/log" // импортируем логгер
)

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") == "gzip" {
			log.Info("Request body is gzipped", "method", r.Method, "url", r.URL.Path)
			gzipReader, err := gzip.NewReader(r.Body)
			if err != nil {
				log.Error("Failed to create gzip reader", "method", r.Method, "url", r.URL.Path, "error", err)
				http.Error(w, "Failed to read gzip data", http.StatusBadRequest)
				return
			}
			defer gzipReader.Close()
			r.Body = gzipReader
		}
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			log.Info("Response will be gzipped", "method", r.Method, "url", r.URL.Path)
			gzipWriter := gzip.NewWriter(w)
			w.Header().Set("Content-Encoding", "gzip")
			defer gzipWriter.Close()

			gzipResponseWriter := &gzipResponseWriter{
				ResponseWriter: w,
				Writer:         gzipWriter,
			}
			next.ServeHTTP(gzipResponseWriter, r)
		} else {
			log.Info("Response will not be gzipped", "method", r.Method, "url", r.URL.Path)
			next.ServeHTTP(w, r)
		}
	})
}

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w *gzipResponseWriter) Write(p []byte) (int, error) {
	return w.Writer.Write(p)
}
