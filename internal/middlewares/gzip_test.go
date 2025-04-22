package middlewares

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestGzipMiddleware_RequestCompression(t *testing.T) {
	logger.Init(zapcore.InfoLevel)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World"))
	})

	middleware := GzipMiddleware(handler)

	compressedBody := &bytes.Buffer{}
	gzipWriter := gzip.NewWriter(compressedBody)
	_, err := gzipWriter.Write([]byte("Hello, World"))
	require.NoError(t, err)
	gzipWriter.Close()

	req := httptest.NewRequest(http.MethodPost, "/", compressedBody)
	req.Header.Set("Content-Encoding", "gzip")

	rr := httptest.NewRecorder()
	middleware.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "Hello, World", rr.Body.String())
}

func TestGzipMiddleware_ResponseCompression(t *testing.T) {
	logger.Init(zapcore.InfoLevel)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World"))
	})

	middleware := GzipMiddleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rr := httptest.NewRecorder()
	middleware.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "gzip", rr.Header().Get("Content-Encoding"))

	// Проверяем, что тело ответа сжато
	gr, err := gzip.NewReader(rr.Body)
	require.NoError(t, err)

	decompressedBody, err := io.ReadAll(gr)
	require.NoError(t, err)

	require.Equal(t, "Hello, World", string(decompressedBody))
}

func TestGzipMiddleware_NoCompression(t *testing.T) {
	logger.Init(zapcore.InfoLevel)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World"))
	})

	middleware := GzipMiddleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	rr := httptest.NewRecorder()
	middleware.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Empty(t, rr.Header().Get("Content-Encoding"))
	require.Equal(t, "Hello, World", rr.Body.String())
}

func TestGzipMiddleware_DecompressionError(t *testing.T) {
	logger.Init(zapcore.InfoLevel)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World"))
	})

	middleware := GzipMiddleware(handler)

	// Некорректное тело запроса, не сжатое в gzip
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("invalid gzip data")))
	req.Header.Set("Content-Encoding", "gzip")

	rr := httptest.NewRecorder()
	middleware.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Contains(t, rr.Body.String(), "failed to read gzip data")
}
