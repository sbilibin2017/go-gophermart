package middlewares

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
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

		rateLimiter.Lock()
		defer rateLimiter.Unlock()

		data, exists := rateLimiter.ips[ip]
		if !exists {
			data = &requestData{count: 0, lastAccess: time.Now()}
			rateLimiter.ips[ip] = data
		}

		if time.Since(data.lastAccess) > time.Minute*time.Duration(rateLimitWindow) {
			data.count = 0
			data.lastAccess = time.Now()
		}

		if data.count >= rateLimit {
			w.Header().Set("Retry-After", "60")
			w.Header().Set("Content-Type", "text/plain")
			http.Error(w, "No more than 60000 requests per minute allowed", http.StatusTooManyRequests)
			return
		}

		data.count++
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
