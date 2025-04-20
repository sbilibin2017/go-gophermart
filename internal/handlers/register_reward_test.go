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

	"github.com/stretchr/testify/assert"
)

func sendRegisterRewardRequest(t *testing.T, handler http.HandlerFunc, requestBody *RegisterRewardRequest) *httptest.ResponseRecorder {
	reqBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("could not marshal request body: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/register/reward", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	handler(w, req)
	return w
}

func compareRegisterRewardResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatusCode int, expectedResponse string) {
	assert.Equal(t, expectedStatusCode, w.Code)
	actualResponse := strings.TrimSpace(w.Body.String())
	expectedResponse = strings.TrimSpace(expectedResponse)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestRegisterRewardHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name               string
		request            *RegisterRewardRequest
		setupMocks         func(mockValidator *MockRegisterRewardValidator, mockService *MockRegisterRewardService)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "successful reward registration",
			request: &RegisterRewardRequest{
				Match:      "match1",
				Reward:     100,
				RewardType: "type1",
			},
			setupMocks: func(mockValidator *MockRegisterRewardValidator, mockService *MockRegisterRewardService) {
				mockValidator.EXPECT().Struct(gomock.Any()).Return(nil)
				mockService.EXPECT().Register(gomock.Any(), "match1", uint64(100), "type1").Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   "Reward registered successfully",
		},
		{
			name: "validation error",
			request: &RegisterRewardRequest{
				Match:      "",
				Reward:     100,
				RewardType: "type1",
			},
			setupMocks: func(mockValidator *MockRegisterRewardValidator, mockService *MockRegisterRewardService) {
				mockValidator.EXPECT().Struct(gomock.Any()).Return(errors.New("validation error"))
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   "Validation error",
		},
		{
			name: "service error: reward already exists",
			request: &RegisterRewardRequest{
				Match:      "match1",
				Reward:     100,
				RewardType: "type1",
			},
			setupMocks: func(mockValidator *MockRegisterRewardValidator, mockService *MockRegisterRewardService) {
				mockValidator.EXPECT().Struct(gomock.Any()).Return(nil)
				mockService.EXPECT().Register(gomock.Any(), "match1", uint64(100), "type1").Return(services.ErrRewardAlreadyExists)
			},
			expectedStatusCode: http.StatusConflict,
			expectedResponse:   "Reward already exists",
		},
		{
			name: "service error: reward not registered",
			request: &RegisterRewardRequest{
				Match:      "match1",
				Reward:     100,
				RewardType: "type1",
			},
			setupMocks: func(mockValidator *MockRegisterRewardValidator, mockService *MockRegisterRewardService) {
				mockValidator.EXPECT().Struct(gomock.Any()).Return(nil)
				mockService.EXPECT().Register(gomock.Any(), "match1", uint64(100), "type1").Return(services.ErrRewardIsNotRegistered)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   "Reward is not registered",
		},
		{
			name:    "invalid request body (decoding error)",
			request: nil,
			setupMocks: func(mockValidator *MockRegisterRewardValidator, mockService *MockRegisterRewardService) {

			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   "Invalid request body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockValidator := NewMockRegisterRewardValidator(ctrl)
			mockService := NewMockRegisterRewardService(ctrl)

			tt.setupMocks(mockValidator, mockService)

			handler := RegisterRewardHandler(mockValidator, mockService)

			var w *httptest.ResponseRecorder
			if tt.request == nil {
				req := httptest.NewRequest(http.MethodPost, "/register/reward", bytes.NewReader([]byte("{ invalid_json }")))
				w = httptest.NewRecorder()
				handler(w, req)
			} else {
				w = sendRegisterRewardRequest(t, handler, tt.request)
			}

			compareRegisterRewardResponse(t, w, tt.expectedStatusCode, tt.expectedResponse)
		})
	}
}
