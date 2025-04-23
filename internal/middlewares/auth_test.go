package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/stretchr/testify/assert"
)

var (
	jwtKey     = []byte("secret")
	validToken = generateValidToken()
)

func generateValidToken() string {
	claims := &claims{
		Login: "testUser",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "testIssuer",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		panic(err)
	}
	return tokenString
}

func TestAuthMiddleware(t *testing.T) {
	logger.Init()
	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectedLogin  string
	}{
		{
			name:           "Valid token",
			authHeader:     "Bearer " + validToken,
			expectedStatus: http.StatusOK,
			expectedLogin:  "testUser",
		},
		{
			name:           "Missing authorization header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectedLogin:  "",
		},
		{
			name:           "Invalid authorization header format",
			authHeader:     "BearerInvalid token",
			expectedStatus: http.StatusUnauthorized,
			expectedLogin:  "",
		},
		{
			name:           "Expired token",
			authHeader:     "Bearer " + "invalidToken",
			expectedStatus: http.StatusUnauthorized,
			expectedLogin:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.Header.Set("Authorization", tt.authHeader)

			rr := httptest.NewRecorder()

			middleware := AuthMiddleware(jwtKey)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				login, ok := GetLogin(r)
				assert.Equal(t, tt.expectedLogin, login)
				assert.True(t, ok)
				assert.Equal(t, http.StatusOK, rr.Code)
			}))

			middleware.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
