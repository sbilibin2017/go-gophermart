package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/sbilibin2017/go-gophermart/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestOrderGetHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := NewMockOrderGetService(ctrl)
	handler := OrderGetHandler(mockService)
	expectedResp := &types.OrderGetResponse{
		Order:  "12345678903",
		Status: types.OrderStatusProcessed,
		Accrual: func() *int64 {
			var val int64 = 150
			return &val
		}(),
	}
	mockService.
		EXPECT().
		Get(gomock.Any(), &types.OrderGetRequest{Number: "12345678903"}).
		Return(expectedResp, nil, nil)
	r := chi.NewRouter()
	r.Get("/orders/{number}", handler)
	req := httptest.NewRequest("GET", "/orders/12345678903", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"order":"12345678903"`)
	assert.Contains(t, w.Body.String(), `"status":"PROCESSED"`)
	assert.Contains(t, w.Body.String(), `"accrual":150`)
}

func TestOrderGetHandler_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := NewMockOrderGetService(ctrl)
	handler := OrderGetHandler(mockService)
	mockService.
		EXPECT().
		Get(gomock.Any(), &types.OrderGetRequest{Number: "invalid"}).
		Return(nil, &types.APIStatus{
			Status:  http.StatusBadRequest,
			Message: "Invalid order number",
		}, errors.New("validation failed"))
	r := chi.NewRouter()
	r.Get("/orders/{number}", handler)
	req := httptest.NewRequest("GET", "/orders/invalid", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid order number")
}
