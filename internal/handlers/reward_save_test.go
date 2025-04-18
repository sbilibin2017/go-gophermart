package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/sbilibin2017/go-gophermart/internal/services"
	"github.com/sbilibin2017/go-gophermart/internal/types"
	"github.com/sbilibin2017/go-gophermart/internal/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupRegisterRewardTest(t *testing.T) (
	*gomock.Controller,
	*MockRegisterGoodRewardService,
	*validator.Validate,
) {
	ctrl := gomock.NewController(t)
	mockService := NewMockRegisterGoodRewardService(ctrl)
	validate := validator.New()
	validate.RegisterValidation("reward_type", validation.ValidateRewardType)
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
			requestBody: types.Reward{
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
			requestBody: types.Reward{
				Match:      "",
				Reward:     10,
				RewardType: "%",
			},
			mockServiceBehavior: func() {},
			expectedStatusCode:  http.StatusBadRequest,
		},
		{
			name: "Validation error - Invalid Reward value",
			requestBody: types.Reward{
				Match:      "Match1",
				Reward:     0,
				RewardType: "%",
			},
			mockServiceBehavior: func() {},
			expectedStatusCode:  http.StatusBadRequest,
		},
		{
			name: "Conflict - Reward already registered",
			requestBody: types.Reward{
				Match:      "Match1",
				Reward:     10,
				RewardType: "%",
			},
			mockServiceBehavior: func() {
				mockService.EXPECT().Register(gomock.Any(), gomock.Any()).Return(services.ErrRewardAlreadyExists)
			},
			expectedStatusCode: http.StatusConflict,
		},
		{
			name: "Reward is not registered",
			requestBody: types.Reward{
				Match:      "Match1",
				Reward:     10,
				RewardType: "%",
			},
			mockServiceBehavior: func() {
				mockService.EXPECT().Register(gomock.Any(), gomock.Any()).Return(services.ErrRewardIsNotRegistered)
			},
			expectedStatusCode: http.StatusBadRequest,
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
			handler := RegisterRewardSaveHandler(validate, mockService)
			handler.ServeHTTP(rr, req)
			assert.Equal(t, tt.expectedStatusCode, rr.Code)
		})
	}
}
