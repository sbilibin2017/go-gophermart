package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"errors"

	"github.com/golang/mock/gomock"
	"github.com/sbilibin2017/go-gophermart/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestRegisterRewardSaveHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockValidator := NewMockRewardSaveValidator(ctrl)
	mockService := NewMockRegisterRewardSaveService(ctrl)

	handler := RegisterRewardSaveHandler(mockValidator, mockService)

	tests := []struct {
		name           string
		requestBody    []byte
		setupMocks     func()
		expectedStatus int
	}{
		{
			name:        "valid JSON",
			requestBody: []byte(`{"match": "test_match", "reward": 100, "reward_type": "%"}`),
			setupMocks: func() {
				mockValidator.EXPECT().Struct(gomock.Any()).Return(nil)
				mockService.EXPECT().Register(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "invalid JSON request",
			requestBody: []byte(`{"match": "test_match", "reward": 100, "reward_type": %}`),
			setupMocks: func() {
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:        "invalid reward structure",
			requestBody: []byte(`{"match": "test_match", "reward": 100, "reward_type": "%"}`),
			setupMocks: func() {
				mockValidator.EXPECT().Struct(gomock.Any()).Return(errors.New(""))
				mockService.EXPECT().Register(gomock.Any(), gomock.Any()).Times(0)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "reward already exists",
			requestBody: []byte(`{"match": "test_match", "reward": 100, "reward_type": "%"}`),
			setupMocks: func() {
				mockValidator.EXPECT().Struct(gomock.Any()).Return(nil)
				mockService.EXPECT().Register(gomock.Any(), gomock.Any()).Return(services.ErrRewardAlreadyExists)
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name:        "error in registering reward",
			requestBody: []byte(`{"match": "test_match", "reward": 100, "reward_type": "%"}`),
			setupMocks: func() {
				mockValidator.EXPECT().Struct(gomock.Any()).Return(nil)
				mockService.EXPECT().Register(gomock.Any(), gomock.Any()).Return(services.ErrRewardIsNotRegistered)
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMocks != nil {
				tt.setupMocks()
			}
			req := httptest.NewRequest(http.MethodPost, "/reward", bytes.NewReader(tt.requestBody))
			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, req)
			assert.Equal(t, tt.expectedStatus, recorder.Code)
		})
	}
}
