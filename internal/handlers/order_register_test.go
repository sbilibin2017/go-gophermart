package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sbilibin2017/go-gophermart/internal/handlers/utils"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/sbilibin2017/go-gophermart/internal/services"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func init() {
	logger.Init(zapcore.InfoLevel)
}

func TestOrderRegisterHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name               string
		reqBody            []byte
		mockDecodeErr      error
		mockValidateErr    error
		mockRegisterErr    error
		expectedStatusCode int
		expectedResponse   string
		decodeTimes        int
		validateTimes      int
		registerTimes      int
	}{
		{
			name:               "failed request decoding",
			reqBody:            []byte(`{"order": 12345}`),
			mockDecodeErr:      errors.New("invalid JSON"),
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   utils.ErrUnprocessableJSON.Error(),
			decodeTimes:        1,
			validateTimes:      0,
			registerTimes:      0,
		},
		{
			name:               "validation failed",
			reqBody:            []byte(`{"order": 12345, "goods": [{"description": "Item 1", "price": 0}]}`),
			mockDecodeErr:      nil,
			mockValidateErr:    errors.New("validation failed"),
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   "validation failed",
			decodeTimes:        1,
			validateTimes:      1,
			registerTimes:      0,
		},
		{
			name:               "order already registered",
			reqBody:            []byte(`{"order": 12345, "goods": [{"description": "Item 1", "price": 100}]}`),
			mockDecodeErr:      nil,
			mockValidateErr:    nil,
			mockRegisterErr:    services.ErrOrderAlreadyRegistered,
			expectedStatusCode: http.StatusConflict,
			expectedResponse:   services.ErrOrderAlreadyRegistered.Error(),
			decodeTimes:        1,
			validateTimes:      1,
			registerTimes:      1,
		},
		{
			name:               "internal server error during registration",
			reqBody:            []byte(`{"order": 12345, "goods": [{"description": "Item 1", "price": 100}]}`),
			mockDecodeErr:      nil,
			mockValidateErr:    nil,
			mockRegisterErr:    errors.New("internal server error"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   ErrOrderIsNotRegistered.Error(),
			decodeTimes:        1,
			validateTimes:      1,
			registerTimes:      1,
		},
		{
			name:               "successful order registration",
			reqBody:            []byte(`{"order": 12345, "goods": [{"description": "Item 1", "price": 100}]}`),
			mockDecodeErr:      nil,
			mockValidateErr:    nil,
			mockRegisterErr:    nil,
			expectedStatusCode: http.StatusAccepted,
			expectedResponse:   "Order registered successfully",
			decodeTimes:        1,
			validateTimes:      1,
			registerTimes:      1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockOrderRegisterService(ctrl)
			mockDecoder := NewMockOrderRegisterDecoder(ctrl)
			mockValidator := NewMockOrderRegisterValidator(ctrl)

			req := &http.Request{
				Method: "POST",
				Body:   io.NopCloser(bytes.NewReader(tt.reqBody)),
			}

			rr := httptest.NewRecorder()

			if tt.mockDecodeErr != nil {
				mockDecoder.EXPECT().Decode(req, gomock.Any()).Return(tt.mockDecodeErr).Times(tt.decodeTimes)
			} else {
				mockDecoder.EXPECT().Decode(req, gomock.Any()).DoAndReturn(func(r *http.Request, v any) error {
					return json.Unmarshal(tt.reqBody, v)
				}).Times(tt.decodeTimes)
			}

			if tt.mockValidateErr != nil {
				mockValidator.EXPECT().Struct(gomock.Any()).Return(tt.mockValidateErr).Times(tt.validateTimes)
			} else {
				mockValidator.EXPECT().Struct(gomock.Any()).Return(nil).Times(tt.validateTimes)
			}

			if tt.mockRegisterErr != nil {
				mockService.EXPECT().Register(context.Background(), gomock.Any()).Return(tt.mockRegisterErr).Times(tt.registerTimes)
			} else {
				mockService.EXPECT().Register(context.Background(), gomock.Any()).Return(nil).Times(tt.registerTimes)
			}

			handler := OrderRegisterHandler(mockService, mockDecoder, mockValidator)

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedResponse)
		})
	}
}
