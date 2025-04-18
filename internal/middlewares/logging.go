package middlewares

import (
	"net/http"
	"time"

	"github.com/sbilibin2017/go-gophermart/internal/logger"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		ww := &ResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(ww, r)
		duration := time.Since(startTime)
		logger.Logger.Info("Request",
			"uri", r.RequestURI,
			"method", r.Method,
			"duration", duration.Seconds(),
		)
		logger.Logger.Info("Response",
			"status", ww.statusCode,
			"response_size", ww.responseSize,
		)
	})
}

type ResponseWriter struct {
	http.ResponseWriter
	statusCode   int
	responseSize int
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *ResponseWriter) Write(p []byte) (n int, err error) {
	n, err = rw.ResponseWriter.Write(p)
	rw.responseSize += n
	return n, err
}
