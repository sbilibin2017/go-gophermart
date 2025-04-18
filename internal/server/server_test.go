package server

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServer_Run(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	srv := NewServer(":8081", handler)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	start := time.Now()
	srv.Run(ctx)

	duration := time.Since(start)
	assert.Less(t, duration.Seconds(), 1.0, "Run должен завершиться после отмены контекста")

	client := &http.Client{}
	resp, err := client.Get("http://localhost:8081")

	if err == nil {
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
}
