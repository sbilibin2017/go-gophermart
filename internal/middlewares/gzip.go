package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/sbilibin2017/go-gophermart/internal/log"
)

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") == "gzip" {
			log.Info("incoming request is gzip-encoded", "url", r.URL.Path)
			gzipReader, err := gzip.NewReader(r.Body)
			if err != nil {
				log.Error("failed to create gzip reader for request body", "error", err)
				http.Error(w, "Failed to read gzip data", http.StatusBadRequest)
				return
			}
			defer func() {
				gzipReader.Close()
			}()
			r.Body = gzipReader
		}

		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			log.Info("client supports gzip, compressing response", "url", r.URL.Path)
			gzipWriter := gzip.NewWriter(w)
			defer func() {
				gzipWriter.Close()
			}()
			w.Header().Set("Content-Encoding", "gzip")
			grw := &gzipResponseWriter{
				ResponseWriter: w,
				Writer:         gzipWriter,
			}
			next.ServeHTTP(grw, r)
		} else {
			log.Info("client does not support gzip", "url", r.URL.Path)
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
