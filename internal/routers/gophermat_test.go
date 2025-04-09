package routers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/stretchr/testify/assert"
)

func registerHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("register"))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("login"))
}

func uploadOrderHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("upload order"))
}

func getOrderHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get order"))
}

func getBalanceHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get balance"))
}

func withdrawBalanceHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("withdraw balance"))
}

func withdrawalsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("withdrawals"))
}

type claims struct {
	Login string `json:"login"`
	jwt.RegisteredClaims
}

func encodeToken(mockConfig *configs.GophermartConfig) (string, error) {
	claims := &claims{
		Login: "test",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(mockConfig.JWTExp)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(mockConfig.JWTSecretKey))
	if err != nil {
		return "", errors.New("err")
	}
	return signedToken, nil
}

func TestNewGophermartRouter(t *testing.T) {
	mockConfig := &configs.GophermartConfig{
		JWTSecretKey: "test",
		JWTExp:       365 * 24 * time.Hour,
	}
	r := NewGophermartRouter(
		mockConfig,
		registerHandler,
		loginHandler,
		uploadOrderHandler,
		getOrderHandler,
		getBalanceHandler,
		withdrawBalanceHandler,
		withdrawalsHandler,
	)
	authHeader, _ := encodeToken(mockConfig)
	authHeader = "Bearer " + authHeader
	req := httptest.NewRequest("POST", "/register", nil)
	req.Header.Set("Authorization", authHeader)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "register", w.Body.String())
	req = httptest.NewRequest("POST", "/login", nil)
	req.Header.Set("Authorization", authHeader)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "login", w.Body.String())
	req = httptest.NewRequest("POST", "/orders", nil)
	req.Header.Set("Authorization", authHeader)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "upload order", w.Body.String())
	req = httptest.NewRequest("GET", "/orders", nil)
	req.Header.Set("Authorization", authHeader)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "get order", w.Body.String())
	req = httptest.NewRequest("GET", "/balance", nil)
	req.Header.Set("Authorization", authHeader)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "get balance", w.Body.String())
	req = httptest.NewRequest("POST", "/balance/withdraw", nil)
	req.Header.Set("Authorization", authHeader)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "withdraw balance", w.Body.String())
	req = httptest.NewRequest("GET", "/withdrawals", nil)
	req.Header.Set("Authorization", authHeader)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "withdrawals", w.Body.String())
}
