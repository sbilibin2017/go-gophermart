package middlewares

import (
	"net/http"
	"time"

	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func init() {
	log = logger.NewLogger(zapcore.InfoLevel)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		ww := &ResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(ww, r)
		duration := time.Since(startTime)

		log.Info("Request",
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.Duration("duration", duration),
		)

		log.Info("Response",
			zap.Int("status", ww.statusCode),
			zap.Int("response_size", ww.responseSize),
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
