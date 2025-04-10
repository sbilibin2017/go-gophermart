package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestApplyMiddlewares(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMiddleware1 := NewMockMiddleware(ctrl)
	mockMiddleware2 := NewMockMiddleware(ctrl)

	finalHandlerCalled := false
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		finalHandlerCalled = true
	})

	mw1Called := false
	mw2Called := false

	mockMiddleware1.
		EXPECT().
		Apply(gomock.Any()).
		DoAndReturn(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mw1Called = true
				next.ServeHTTP(w, r)
			})
		})

	mockMiddleware2.
		EXPECT().
		Apply(gomock.Any()).
		DoAndReturn(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mw2Called = true
				next.ServeHTTP(w, r)
			})
		})

	handler := ApplyMiddlewares(finalHandler, mockMiddleware1, mockMiddleware2)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.True(t, mw1Called, "Middleware 1 should have been called")
	assert.True(t, mw2Called, "Middleware 2 should have been called")
	assert.True(t, finalHandlerCalled, "Final handler should have been called")
}

func TestWrapHandler(t *testing.T) {
	called := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	wrapped := WrapHandler(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	wrapped(w, req, nil)

	assert.True(t, called, "Handler should have been called")
}

func TestWrap(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMiddleware := NewMockMiddleware(ctrl)

	handlerCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})

	mockMiddleware.
		EXPECT().
		Apply(gomock.Any()).
		DoAndReturn(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			})
		})

	wrapped := Wrap(handler, mockMiddleware)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	wrapped(w, req, nil)

	assert.True(t, handlerCalled, "Handler should have been called after wrapping")
}
