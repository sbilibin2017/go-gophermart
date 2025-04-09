package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// Function to generate a valid JWT token
func generateTestToken(secretKey string) (string, error) {
	// Create the claims
	claims := jwt.MapClaims{
		"login": "testUser",
		"exp":   time.Now().Add(time.Hour * 1).Unix(), // token expires in 1 hour
	}
	// Create the token using HS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func TestAuthMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockJWTConfig := NewMockJWTConfig(ctrl)
	mockJWTConfig.EXPECT().GetJWTSecretKey().Return("mockSecretKey").AnyTimes()
	middleware := AuthMiddleware(mockJWTConfig)
	tokenString, err := generateTestToken("mockSecretKey")
	assert.NoError(t, err)
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		login, ok := r.Context().Value(loginKey).(string)
		assert.True(t, ok, "Login should be in the context")
		assert.Equal(t, "testUser", login, "Expected login to be 'testUser'")
		w.WriteHeader(http.StatusOK)
	})
	handler := middleware(testHandler)
	req, err := http.NewRequest("GET", "/test", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	t.Log("Response status:", rr.Code)
	t.Log("Response body:", rr.Body.String())
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestAuthMiddleware_MissingAuthorizationHeader(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockJWTConfig := NewMockJWTConfig(ctrl)
	mockJWTConfig.EXPECT().GetJWTSecretKey().Return("mockSecretKey").AnyTimes()
	middleware := AuthMiddleware(mockJWTConfig)
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := middleware(testHandler)
	req, err := http.NewRequest("GET", "/test", nil)
	assert.NoError(t, err)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), ErrAuthHeaderMissing.Error())
}

func TestAuthMiddleware_InvalidAuthorizationHeaderFormat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockJWTConfig := NewMockJWTConfig(ctrl)
	mockJWTConfig.EXPECT().GetJWTSecretKey().Return("mockSecretKey").AnyTimes()
	middleware := AuthMiddleware(mockJWTConfig)
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := middleware(testHandler)
	req, err := http.NewRequest("GET", "/test", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "InvalidFormat token")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), ErrInvalidAuthHeaderFormat.Error())
}
