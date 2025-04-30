package middlewares

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGzipMiddleware_DecompressRequest(t *testing.T) {
	var body string
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := io.ReadAll(r.Body)
		body = string(data)
	})

	middleware := GzipMiddleware(handler)

	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	gz.Write([]byte("compressed request"))
	gz.Close()

	req := httptest.NewRequest("POST", "/", &b)
	req.Header.Set("Content-Encoding", "gzip")

	rec := httptest.NewRecorder()
	middleware.ServeHTTP(rec, req)

	assert.Equal(t, "compressed request", body)
}

func TestGzipMiddleware_CompressResponse(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("original response"))
	})

	middleware := GzipMiddleware(handler)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	middleware.ServeHTTP(rec, req)

	assert.Equal(t, "gzip", rec.Header().Get("Content-Encoding"))

	gr, err := gzip.NewReader(rec.Body)
	assert.NoError(t, err)
	defer gr.Close()

	data, err := io.ReadAll(gr)
	assert.NoError(t, err)
	assert.Equal(t, "original response", string(data))
}

func TestGzipMiddleware_NoCompression(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("plain response"))
	})

	middleware := GzipMiddleware(handler)

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	middleware.ServeHTTP(rec, req)

	assert.Equal(t, "", rec.Header().Get("Content-Encoding"))
	assert.Equal(t, "plain response", rec.Body.String())
}

func TestGzipMiddleware_BadGzipRequest(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called on bad gzip data")
	})

	middleware := GzipMiddleware(handler)

	badGzipData := []byte("not really gzip")
	req := httptest.NewRequest("POST", "/", bytes.NewReader(badGzipData))
	req.Header.Set("Content-Encoding", "gzip")

	rec := httptest.NewRecorder()
	middleware.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "Failed to read gzip data\n", rec.Body.String())
}
