package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetLoginFromContext(t *testing.T) {
	ctx := context.Background()
	ctx = setLoginToContext(ctx, "testlogin")
	login, err := GetLoginFromContext(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "testlogin", login)

	ctx = context.Background()
	login, err = GetLoginFromContext(ctx)
	assert.Error(t, err)
	assert.Equal(t, "", login)
}

func TestGetLoginFromToken(t *testing.T) {
	jwtSecretKey := "secret"
	claims := &struct {
		jwt.RegisteredClaims
		Login string `json:"login"`
	}{Login: "testlogin"}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecretKey))
	assert.NoError(t, err)

	login, err := getLoginFromToken(tokenString, jwtSecretKey)
	assert.NoError(t, err)
	assert.Equal(t, "testlogin", login)

	invalidToken := "invalidtoken"
	login, err = getLoginFromToken(invalidToken, jwtSecretKey)
	assert.Error(t, err)
	assert.Equal(t, "", login)
}

func TestSetLoginToContext(t *testing.T) {
	ctx := context.Background()
	ctx = setLoginToContext(ctx, "testlogin")
	login, err := GetLoginFromContext(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "testlogin", login)
}

func TestGetTokenStringFromHeader(t *testing.T) {
	req, err := http.NewRequest("GET", "/test", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer testtoken")

	token, err := getTokenStringFromHeader(req)
	assert.NoError(t, err)
	assert.Equal(t, "testtoken", token)

	req.Header.Set("Authorization", "InvalidToken")
	token, err = getTokenStringFromHeader(req)
	assert.Error(t, err)
	assert.Equal(t, "", token)
}

func TestHandleUnauthorizedError(t *testing.T) {
	rr := httptest.NewRecorder()
	handleUnauthorizedError(rr)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "Unauthorized")
}

func TestAuthMiddleware_Success(t *testing.T) {
	jwtSecretKey := "secret"
	claims := &struct {
		jwt.RegisteredClaims
		Login string `json:"login"`
	}{Login: "testlogin"}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecretKey))
	assert.NoError(t, err)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		login, err := GetLoginFromContext(r.Context())
		assert.NoError(t, err)
		assert.Equal(t, "testlogin", login)
		w.WriteHeader(http.StatusOK)
	})
	middleware := AuthMiddleware(jwtSecretKey)(handler)

	req, err := http.NewRequest("GET", "/test", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+tokenString)

	rr := httptest.NewRecorder()
	middleware.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestAuthMiddleware_MissingToken(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	middleware := AuthMiddleware("secret")(handler)

	req, err := http.NewRequest("GET", "/test", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	middleware.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	middleware := AuthMiddleware("secret")(handler)

	req, err := http.NewRequest("GET", "/test", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer invalidtoken")

	rr := httptest.NewRecorder()
	middleware.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}
