package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/jwt"
)

func TestAuthMiddleware(t *testing.T) {
	jwtConfig := &configs.JWTConfig{
		SecretKey: "testsecret",
		Exp:       time.Hour,
	}
	validLogin := "testuser"
	validToken, err := jwt.Generate(jwtConfig, validLogin)
	require.NoError(t, err)
	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectLogin    bool
	}{
		{
			name:           "no authorization header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectLogin:    false,
		},
		{
			name:           "invalid auth format",
			authHeader:     "Basic token123",
			expectedStatus: http.StatusUnauthorized,
			expectLogin:    false,
		},
		{
			name:           "invalid token",
			authHeader:     "Bearer invalidtoken",
			expectedStatus: http.StatusUnauthorized,
			expectLogin:    false,
		},
		{
			name:           "valid token",
			authHeader:     "Bearer " + validToken,
			expectedStatus: http.StatusOK,
			expectLogin:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := AuthMiddleware(jwtConfig)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				login := r.Context().Value(userLoginKey)
				if tc.expectLogin {
					assert.Equal(t, validLogin, login)
				} else {
					assert.Nil(t, login)
				}
				w.WriteHeader(http.StatusOK)
			}))
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedStatus, rec.Code)
		})
	}
}
