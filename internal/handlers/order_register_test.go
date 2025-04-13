package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	gomock "github.com/golang/mock/gomock"
	j "github.com/sbilibin2017/go-gophermart/internal/json"
	"github.com/sbilibin2017/go-gophermart/internal/services"
	"github.com/sbilibin2017/go-gophermart/internal/usecases"
	"github.com/stretchr/testify/assert"
)

func TestOrderRegisterHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := NewMockOrderRegisterUsecase(ctrl)
	mockDecoder := NewMockOrderRegisterRequestDecoder(ctrl)

	tests := []struct {
		name           string
		requestBody    interface{}
		mockDecode     func()
		mockExecute    func()
		expectedStatus int
	}{
		{
			name: "Valid request",
			requestBody: &usecases.OrderRegisterRequest{
				Order: 123,
				Goods: []usecases.Good{
					{Description: "Item 1", Price: 100},
					{Description: "Item 2", Price: 200},
				},
			},
			mockDecode: func() {
				mockDecoder.EXPECT().Decode(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			mockExecute: func() {
				mockUsecase.EXPECT().Execute(gomock.Any(), gomock.Any()).
					Return(&usecases.OrderRegisterResponse{Message: "Order registered successfully"}, nil)
			},
			expectedStatus: http.StatusAccepted,
		},
		{
			name:        "Invalid JSON",
			requestBody: "{invalid json",
			mockDecode: func() {
				mockDecoder.EXPECT().Decode(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(j.ErrRequestDecoderUnprocessableJSON)
			},
			mockExecute:    func() {}, // не вызывается
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Order already registered",
			requestBody: &usecases.OrderRegisterRequest{
				Order: 123,
				Goods: []usecases.Good{
					{Description: "Item 1", Price: 100},
				},
			},
			mockDecode: func() {
				mockDecoder.EXPECT().Decode(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			mockExecute: func() {
				mockUsecase.EXPECT().Execute(gomock.Any(), gomock.Any()).
					Return(nil, services.ErrOrderAlreadyRegistered)
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name: "Internal server error",
			requestBody: &usecases.OrderRegisterRequest{
				Order: 123,
				Goods: []usecases.Good{
					{Description: "Item 1", Price: 100},
				},
			},
			mockDecode: func() {
				mockDecoder.EXPECT().Decode(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			mockExecute: func() {
				mockUsecase.EXPECT().Execute(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Invalid request by business logic",
			requestBody: &usecases.OrderRegisterRequest{
				Order: 0,
				Goods: []usecases.Good{
					{Description: "Invalid good", Price: 100},
				},
			},
			mockDecode: func() {
				mockDecoder.EXPECT().Decode(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			mockExecute: func() {
				mockUsecase.EXPECT().Execute(gomock.Any(), gomock.Any()).
					Return(nil, usecases.ErrOrderRegisterInvalidRequest)
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockDecode()
			tt.mockExecute()

			var reqBody []byte
			var err error
			switch v := tt.requestBody.(type) {
			case string:
				reqBody = []byte(v)
			default:
				reqBody, err = json.Marshal(v)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(reqBody))
			rr := httptest.NewRecorder()

			handler := OrderRegisterHandler(mockUsecase, mockDecoder)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
