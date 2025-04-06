package routers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAccrualRouter(t *testing.T) {
	getOrderHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("order details"))
	})
	createOrderHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("order created"))
	})
	registerGoodsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("goods registered"))
	})
	r := NewAccrualRouter(getOrderHandler, createOrderHandler, registerGoodsHandler)
	req1, err := http.NewRequest("GET", "/api/orders/12345", nil)
	assert.NoError(t, err)
	rr1 := httptest.NewRecorder()
	r.ServeHTTP(rr1, req1)
	assert.Equal(t, http.StatusOK, rr1.Code)
	assert.Equal(t, "order details", rr1.Body.String())
	req2, err := http.NewRequest("POST", "/api/orders", nil)
	assert.NoError(t, err)
	rr2 := httptest.NewRecorder()
	r.ServeHTTP(rr2, req2)
	assert.Equal(t, http.StatusCreated, rr2.Code)
	assert.Equal(t, "order created", rr2.Body.String())
	req3, err := http.NewRequest("POST", "/api/goods", nil)
	assert.NoError(t, err)
	rr3 := httptest.NewRecorder()
	r.ServeHTTP(rr3, req3)
	assert.Equal(t, http.StatusOK, rr3.Code)
	assert.Equal(t, "goods registered", rr3.Body.String())
}
