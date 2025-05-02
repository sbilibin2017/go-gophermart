package middlewares

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGzipMiddleware_DecompressRequest(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body) // Replaced ioutil.ReadAll with io.ReadAll
		assert.NoError(t, err)
		assert.Equal(t, "Hello, World!", string(body))
		w.WriteHeader(http.StatusOK)
	})
	middleware := GzipMiddleware(handler)
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	_, err := gzipWriter.Write([]byte("Hello, World!"))
	assert.NoError(t, err)
	err = gzipWriter.Close()
	assert.NoError(t, err)
	req, err := http.NewRequest("POST", "/test", &buf)
	assert.NoError(t, err)
	req.Header.Set("Content-Encoding", "gzip")
	rr := httptest.NewRecorder()
	middleware.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestGzipMiddleware_CompressResponse(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
	middleware := GzipMiddleware(handler)
	req, err := http.NewRequest("GET", "/test", nil)
	assert.NoError(t, err)
	req.Header.Set("Accept-Encoding", "gzip")
	rr := httptest.NewRecorder()
	middleware.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "gzip", rr.Header().Get("Content-Encoding"))
	gzipReader, err := gzip.NewReader(rr.Body)
	assert.NoError(t, err)
	defer gzipReader.Close()
	decompressedBody, err := io.ReadAll(gzipReader) // Replaced ioutil.ReadAll with io.ReadAll
	assert.NoError(t, err)
	assert.Equal(t, "Hello, World!", string(decompressedBody))
}

func TestGzipMiddleware_NoCompression(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
	middleware := GzipMiddleware(handler)
	req, err := http.NewRequest("GET", "/test", nil)
	assert.NoError(t, err)
	rr := httptest.NewRecorder()
	middleware.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Empty(t, rr.Header().Get("Content-Encoding"))
}

func TestGzipMiddleware_BadGzipRequest(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
	middleware := GzipMiddleware(handler)
	req, err := http.NewRequest("POST", "/test", strings.NewReader("Invalid Gzip Data"))
	assert.NoError(t, err)
	req.Header.Set("Content-Encoding", "gzip")
	rr := httptest.NewRecorder()
	middleware.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Failed to read gzip data")
}
