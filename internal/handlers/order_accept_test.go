package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sbilibin2017/go-gophermart/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestOrderAcceptHandler_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := NewMockOrderAcceptService(ctrl)
	handler := OrderAcceptHandler(mockService)
	invalidJSON := `{"order":"12345678903","goods":[{"description":"Item","price":1000}`
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request body")
}

func TestOrderAcceptHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := NewMockOrderAcceptService(ctrl)
	handler := OrderAcceptHandler(mockService)
	orderJSON := `{"order":"12345678903","goods":[{"description":"Item","price":1000}]}`
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBufferString(orderJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mockService.EXPECT().
		Accept(gomock.Any(), &types.OrderAcceptRequest{
			Order: "12345678903",
			Goods: []types.Good{{Description: "Item", Price: 1000}},
		}).
		Return(&types.APIStatus{
			Status:  http.StatusAccepted,
			Message: "Order accepted",
		}, nil, nil)
	handler(w, req)
	assert.Equal(t, http.StatusAccepted, w.Code)
	assert.Equal(t, "text/plain", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Body.String(), "Order accepted")
}

func TestOrderAcceptHandler_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := NewMockOrderAcceptService(ctrl)
	handler := OrderAcceptHandler(mockService)
	orderJSON := `{"order":"invalid","goods":[{"description":"Item","price":1000}]}`
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBufferString(orderJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mockService.EXPECT().
		Accept(gomock.Any(), &types.OrderAcceptRequest{
			Order: "invalid",
			Goods: []types.Good{{Description: "Item", Price: 1000}},
		}).
		Return(nil, &types.APIStatus{
			Status:  http.StatusBadRequest,
			Message: "Invalid order",
		}, errors.New("bad order"))
	handler(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid order")
}
