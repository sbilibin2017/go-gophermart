package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/sbilibin2017/go-gophermart/internal/domain"
	"github.com/sbilibin2017/go-gophermart/internal/handlers/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupRegisterRewardTest(t *testing.T) (
	*gomock.Controller,
	*MockRegisterRewardService,
	*validator.Validate,
) {
	ctrl := gomock.NewController(t)
	mockService := NewMockRegisterRewardService(ctrl)
	validate := validator.New()
	return ctrl, mockService, validate
}

func TestRegisterRewardHandler(t *testing.T) {
	ctrl, mockService, validate := setupRegisterRewardTest(t)
	defer ctrl.Finish()

	tests := []struct {
		name                string
		requestBody         interface{}
		mockServiceBehavior func()
		expectedStatusCode  int
	}{
		{
			name:                "Invalid JSON body",
			requestBody:         `{"match": "Match1", "reward": 10`,
			mockServiceBehavior: func() {},
			expectedStatusCode:  http.StatusBadRequest,
		},
		{
			name: "Successful registration",
			requestBody: RegisterGoodRewardRequest{
				Match:      "Match1",
				Reward:     10,
				RewardType: "%",
			},
			mockServiceBehavior: func() {
				mockService.EXPECT().Register(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Validation error - Missing required field",
			requestBody: RegisterGoodRewardRequest{
				Match:      "",
				Reward:     10,
				RewardType: "%",
			},
			mockServiceBehavior: func() {},
			expectedStatusCode:  http.StatusBadRequest,
		},
		{
			name: "Validation error - Invalid Reward value",
			requestBody: RegisterGoodRewardRequest{
				Match:      "Match1",
				Reward:     0,
				RewardType: "%",
			},
			mockServiceBehavior: func() {},
			expectedStatusCode:  http.StatusBadRequest,
		},
		{
			name: "Conflict - Reward already registered",
			requestBody: RegisterGoodRewardRequest{
				Match:      "Match1",
				Reward:     10,
				RewardType: "%",
			},
			mockServiceBehavior: func() {
				mockService.EXPECT().Register(gomock.Any(), gomock.Any()).Return(domain.ErrRewardSearchKeyAlreadyRegistered)
			},
			expectedStatusCode: http.StatusConflict,
		},
		{
			name: "Internal Server Error",
			requestBody: RegisterGoodRewardRequest{
				Match:      "Match1",
				Reward:     10,
				RewardType: "%",
			},
			mockServiceBehavior: func() {
				mockService.EXPECT().Register(gomock.Any(), gomock.Any()).Return(utils.ErrInternal)
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockServiceBehavior()
			body, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)
			req := httptest.NewRequest(http.MethodPost, "/register-reward", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			handler := RegisterGoodRewardHandler(validate, mockService)
			handler.ServeHTTP(rr, req)
			assert.Equal(t, tt.expectedStatusCode, rr.Code)
		})
	}
}
