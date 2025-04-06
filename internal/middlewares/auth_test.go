package middlewares

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sbilibin2017/go-gophermart/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := NewMockUserRepository(ctrl)
	mockDecoder := NewMockJWTDecoder(ctrl)
	middleware := AuthMiddleware(mockRepo, mockDecoder)
	tests := []struct {
		name               string
		token              string
		mockDecodeResp     *types.Claims
		mockDecodeErr      error
		mockGetByIDResp    *types.User
		mockGetByIDErr     error
		expectedStatusCode int
		expectDecodeCall   bool
	}{
		{
			name:               "valid token and user",
			token:              "valid_token",
			mockDecodeResp:     types.NewClaims(123),
			mockDecodeErr:      nil,
			mockGetByIDResp:    types.NewUser(123),
			mockGetByIDErr:     nil,
			expectedStatusCode: http.StatusOK,
			expectDecodeCall:   true,
		},
		{
			name:               "missing token",
			token:              "",
			mockDecodeResp:     nil,
			mockDecodeErr:      nil,
			mockGetByIDResp:    nil,
			mockGetByIDErr:     nil,
			expectedStatusCode: http.StatusUnauthorized,
			expectDecodeCall:   false,
		},
		{
			name:               "invalid token",
			token:              "invalid_token",
			mockDecodeResp:     nil,
			mockDecodeErr:      errors.New("invalid token"),
			mockGetByIDResp:    nil,
			mockGetByIDErr:     nil,
			expectedStatusCode: http.StatusUnauthorized,
			expectDecodeCall:   true,
		},
		{
			name:               "user not found",
			token:              "valid_token",
			mockDecodeResp:     types.NewClaims(123),
			mockDecodeErr:      nil,
			mockGetByIDResp:    nil,
			mockGetByIDErr:     errors.New("user not found"),
			expectedStatusCode: http.StatusUnauthorized,
			expectDecodeCall:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectDecodeCall {
				mockDecoder.EXPECT().Decode(tt.token).Return(tt.mockDecodeResp, tt.mockDecodeErr)
			} else {
				mockDecoder.EXPECT().Decode(tt.token).Times(0)
			}
			if tt.mockDecodeResp != nil && tt.mockDecodeResp.UserID != nil {
				mockRepo.EXPECT().GetByID(gomock.Any(), tt.mockDecodeResp.UserID).Return(tt.mockGetByIDResp, tt.mockGetByIDErr)
			} else {
				mockRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Times(0)
			}
			req, err := http.NewRequest(http.MethodGet, "/", nil)
			require.NoError(t, err)
			if tt.token != "" {
				req.Header.Set("Authorization", tt.token)
			}
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})
			rr := httptest.NewRecorder()
			middleware(nextHandler).ServeHTTP(rr, req)
			assert.Equal(t, tt.expectedStatusCode, rr.Code)
		})
	}
}
