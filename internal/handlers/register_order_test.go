package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sbilibin2017/go-gophermart/internal/services"
	"github.com/sbilibin2017/go-gophermart/internal/types"
	"github.com/stretchr/testify/assert"
)

func sendRegisterOrderRequest(t *testing.T, handler http.HandlerFunc, requestBody *types.RegisterOrderRequest) *httptest.ResponseRecorder {
	reqBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("could not marshal request body: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/orders", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	handler(w, req)
	return w
}

func compareRegisterOrderResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatusCode int, expectedResponse string) {
	assert.Equal(t, expectedStatusCode, w.Code)
	actualResponse := strings.TrimSpace(w.Body.String())
	expectedResponse = strings.TrimSpace(expectedResponse)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestRegisterOrderHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name               string
		request            *types.RegisterOrderRequest
		setupMocks         func(mockValidator *MockRegisterOrderValidator, mockService *MockRegisterOrderService)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "successful order registration",
			request: &types.RegisterOrderRequest{
				Order: "12345678903",
			},
			setupMocks: func(mockValidator *MockRegisterOrderValidator, mockService *MockRegisterOrderService) {
				mockValidator.EXPECT().Struct(gomock.Any()).Return(nil)
				mockService.EXPECT().Register(gomock.Any(), &types.RegisterOrderRequest{
					Order: "12345678903",
				}).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   "Order registered successfully",
		},
		{
			name: "validation error",
			request: &types.RegisterOrderRequest{
				Order: "",
			},
			setupMocks: func(mockValidator *MockRegisterOrderValidator, mockService *MockRegisterOrderService) {
				mockValidator.EXPECT().Struct(gomock.Any()).Return(errors.New("validation error"))
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   "Validation error",
		},
		{
			name: "service error: order already exists",
			request: &types.RegisterOrderRequest{
				Order: "12345678903",
			},
			setupMocks: func(mockValidator *MockRegisterOrderValidator, mockService *MockRegisterOrderService) {
				mockValidator.EXPECT().Struct(gomock.Any()).Return(nil)
				mockService.EXPECT().Register(gomock.Any(), &types.RegisterOrderRequest{
					Order: "12345678903",
				}).Return(services.ErrRegisterOrderAlreadyExists)
			},
			expectedStatusCode: http.StatusConflict,
			expectedResponse:   "Order already exists",
		},
		{
			name: "service error: order is not registered",
			request: &types.RegisterOrderRequest{
				Order: "12345678903",
			},
			setupMocks: func(mockValidator *MockRegisterOrderValidator, mockService *MockRegisterOrderService) {
				mockValidator.EXPECT().Struct(gomock.Any()).Return(nil)
				mockService.EXPECT().Register(gomock.Any(), &types.RegisterOrderRequest{
					Order: "12345678903",
				}).Return(services.ErrRegisterOrderIsNotRegistered)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   "Order is not registered",
		},
		{
			name:               "invalid request body (decoding error)",
			request:            nil,
			setupMocks:         func(mockValidator *MockRegisterOrderValidator, mockService *MockRegisterOrderService) {},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   "Invalid request body",
		},
		{
			name: "service error: unknown",
			request: &types.RegisterOrderRequest{
				Order: "12345678903",
			},
			setupMocks: func(mockValidator *MockRegisterOrderValidator, mockService *MockRegisterOrderService) {
				mockValidator.EXPECT().Struct(gomock.Any()).Return(nil)
				mockService.EXPECT().Register(gomock.Any(), &types.RegisterOrderRequest{
					Order: "12345678903",
				}).Return(errors.New("unexpected error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   "Unexpected error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockValidator := NewMockRegisterOrderValidator(ctrl)
			mockService := NewMockRegisterOrderService(ctrl)

			tt.setupMocks(mockValidator, mockService)

			handler := RegisterOrderHandler(mockValidator, mockService)

			var w *httptest.ResponseRecorder
			if tt.request == nil {
				req := httptest.NewRequest(http.MethodPost, "/api/orders", bytes.NewReader([]byte("{ invalid_json }")))
				w = httptest.NewRecorder()
				handler(w, req)
			} else {
				w = sendRegisterOrderRequest(t, handler, tt.request)
			}

			compareRegisterOrderResponse(t, w, tt.expectedStatusCode, tt.expectedResponse)
		})
	}
}
