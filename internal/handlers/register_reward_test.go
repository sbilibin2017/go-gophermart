package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/sbilibin2017/go-gophermart/internal/domain"
)

func TestRegisterRewardHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockRegisterRewardService(ctrl)
	mockValidator := NewMockRegisterRewardValidator(ctrl)

	handler := RegisterRewardHandler(mockService, mockValidator)

	tests := []struct {
		name         string
		body         string
		setupMocks   func()
		expectedCode int
	}{
		{
			name: "success",
			body: `{"match": "order123", "reward": 100, "reward_type": "points"}`,
			setupMocks: func() {
				mockValidator.EXPECT().
					Struct(gomock.Any()).
					Return(nil)

				mockService.EXPECT().
					Register(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "invalid_json",
			body: `{"match": "order123", "reward": "oops"}`,
			setupMocks: func() {
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "validation_failed",
			body: `{"match": "order123", "reward": 100, "reward_type": "points"}`,
			setupMocks: func() {
				mockValidator.EXPECT().
					Struct(gomock.Any()).
					Return(errors.New("some validation error"))
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "duplicate_reward",
			body: `{"match": "order123", "reward": 100, "reward_type": "points"}`,
			setupMocks: func() {
				mockValidator.EXPECT().
					Struct(gomock.Any()).
					Return(nil)
				mockService.EXPECT().
					Register(gomock.Any(), gomock.Any()).
					Return(domain.ErrRewardKeyAlreadyRegistered)
			},
			expectedCode: http.StatusConflict,
		},
		{
			name: "unexpected_error",
			body: `{"match": "order123", "reward": 100, "reward_type": "points"}`,
			setupMocks: func() {
				mockValidator.EXPECT().
					Struct(gomock.Any()).
					Return(nil)
				mockService.EXPECT().
					Register(gomock.Any(), gomock.Any()).
					Return(errors.New("db error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks()
			}

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(tc.body))
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)

		})
	}
}
