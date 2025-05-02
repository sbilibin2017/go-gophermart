package middlewares

import (
	"log"
	"net/http"
	"sync"
	"time"
)

const maxRequestsPerSecond = 10000

func RateLimitMiddleware(next http.Handler) http.Handler {
	var mu sync.Mutex
	var requestCount int
	var lastTimestamp time.Time

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		currentTimestamp := time.Now()

		if currentTimestamp.Sub(lastTimestamp).Seconds() >= 1 {
			lastTimestamp = currentTimestamp
			requestCount = 0
		}

		requestCount++

		if requestCount > maxRequestsPerSecond {
			w.Header().Set("Retry-After", "1") // Retry-After header
			http.Error(w, "Too many requests, please try again later.", http.StatusTooManyRequests)
			log.Println("Rate limit exceeded")
			return
		}

		next.ServeHTTP(w, r)
	})
}
