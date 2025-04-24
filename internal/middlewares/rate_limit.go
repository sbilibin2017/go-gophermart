package middlewares

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"go.uber.org/zap"
)

const rateLimit = 60000
const rateLimitWindow = 1

var rateLimiter = struct {
	sync.RWMutex
	ips map[string]*requestData
}{
	ips: make(map[string]*requestData),
}

type requestData struct {
	count      int
	lastAccess time.Time
}

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getIP(r)

		if logger.Logger != nil {
			logger.Logger.Info("RateLimitMiddleware: Received request", zap.String("IP", ip))
		}

		rateLimiter.Lock()
		defer rateLimiter.Unlock()

		data, exists := rateLimiter.ips[ip]
		if !exists {
			data = &requestData{count: 0, lastAccess: time.Now()}
			rateLimiter.ips[ip] = data
			if logger.Logger != nil {
				logger.Logger.Info("RateLimitMiddleware: New IP added", zap.String("IP", ip))
			}
		}

		if time.Since(data.lastAccess) > time.Minute*time.Duration(rateLimitWindow) {
			data.count = 0
			data.lastAccess = time.Now()
			if logger.Logger != nil {
				logger.Logger.Info("RateLimitMiddleware: Rate limit window reset", zap.String("IP", ip))
			}
		}

		if data.count >= rateLimit {
			if logger.Logger != nil {
				logger.Logger.Warn("RateLimitMiddleware: Rate limit exceeded", zap.String("IP", ip), zap.Int("count", data.count))
			}
			w.Header().Set("Retry-After", "60")
			w.Header().Set("Content-Type", "text/plain")
			http.Error(w, "No more than 60000 requests per minute allowed", http.StatusTooManyRequests)
			return
		}

		data.count++
		if logger.Logger != nil {
			logger.Logger.Info("RateLimitMiddleware: Request allowed", zap.String("IP", ip), zap.Int("count", data.count))
		}

		next.ServeHTTP(w, r)
	})
}

func getIP(r *http.Request) string {
	xfHeader := r.Header.Get("X-Forwarded-For")
	if xfHeader != "" {
		ips := strings.Split(xfHeader, ",")
		return strings.TrimSpace(ips[0])
	}
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	return host
}
