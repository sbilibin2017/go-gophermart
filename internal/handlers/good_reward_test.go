package handlers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sbilibin2017/go-gophermart/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestGoodRewardHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMockGoodRewardService(ctrl)

	tests := []struct {
		name            string
		requestBody     []byte
		mockResponse    *types.APIStatus
		mockError       error
		expectedCode    int
		expectedMessage string
	}{
		{
			name: "Valid Request - Success",
			requestBody: []byte(`{
				"match": "12345",
				"reward": 100,
				"reward_type": "%"
			}`),
			mockResponse: &types.APIStatus{
				Status:  200,
				Message: "Good reward registered successfully",
			},
			mockError:       nil,
			expectedCode:    http.StatusOK,
			expectedMessage: "Good reward registered successfully",
		},
		{
			name: "Error during register",
			requestBody: []byte(`{
				"match": "12345",
				"reward": 100,
				"reward_type": "%"
			}`),
			mockResponse: &types.APIStatus{
				Status:  500,
				Message: "Internal Server Error",
			},
			mockError:       errors.New("Internal Server Error"),
			expectedCode:    http.StatusInternalServerError,
			expectedMessage: "Internal Server Error\n",
		},
		{
			name: "Invalid JSON - Malformed",
			requestBody: []byte(`{
				"match": "12345",
				"reward": 100,
				"reward_type": "%"
			`),
			mockResponse:    nil,
			mockError:       nil,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "Invalid request body\n",
		},
		{
			name:            "Empty Body",
			requestBody:     []byte(``),
			mockResponse:    nil,
			mockError:       nil,
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "Invalid request body\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockResponse != nil {
				mockSvc.EXPECT().Register(context.Background(), gomock.Any()).Return(tt.mockResponse, tt.mockError).Times(1)
			}

			req := httptest.NewRequest(http.MethodPost, "/good/reward", bytes.NewReader(tt.requestBody))
			rec := httptest.NewRecorder()

			handler := RewardHandler(mockSvc)
			handler.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
			assert.Equal(t, tt.expectedMessage, rec.Body.String())
		})
	}
}
