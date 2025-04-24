package middlewares

import (
	"compress/gzip"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"go.uber.org/zap"
)

var (
	ErrFailedToDecompressRequest = errors.New("failed to read gzip data")
)

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") == "gzip" {
			gzipReader, err := gzip.NewReader(r.Body)
			if err != nil {
				if logger.Logger != nil {
					logger.Logger.Error("failed to create gzip reader for request body", zap.Error(err))
				}
				http.Error(w, ErrFailedToDecompressRequest.Error(), http.StatusBadRequest)
				return
			}
			defer func() {
				gzipReader.Close()
			}()
			r.Body = io.NopCloser(gzipReader)
			if logger.Logger != nil {
				logger.Logger.Info("request body decompressed with gzip")
			}
		}

		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			w.Header().Set("Content-Encoding", "gzip")
			gzipWriter := gzip.NewWriter(w)
			defer func() {
				gzipWriter.Close()
			}()

			if logger.Logger != nil {
				logger.Logger.Info("response will be compressed with gzip")
			}

			grw := &gzipResponseWriter{
				ResponseWriter: w,
				Writer:         gzipWriter,
			}
			next.ServeHTTP(grw, r)

			if logger.Logger != nil {
				logger.Logger.Info("GzipMiddleware: request processed with gzip response")
			}
			return
		}

		if logger.Logger != nil {
			logger.Logger.Info("GzipMiddleware: no compression used")
		}
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
