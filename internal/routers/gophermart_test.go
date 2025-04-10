package routers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/stretchr/testify/assert"
)

func registerHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("register"))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("login"))
}

func uploadOrderHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("order uploaded"))
}

func getOrderHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("orders"))
}

func getBalanceHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("balance"))
}

func withdrawBalanceHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("balance withdrawn"))
}

func withdrawalsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("withdrawals"))
}

type Claims struct {
	Login string `json:"login"`
	jwt.RegisteredClaims
}

func Generate(login string, secret string, expireTime time.Duration) (string, error) {
	claims := &Claims{
		Login: login,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "test",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireTime)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func TestNewGophermartRouter(t *testing.T) {
	config := &configs.GophermartConfig{
		JWTSecretKey: "secret",
	}

	tokenString, err := Generate("testuser", config.JWTSecretKey, time.Hour)
	assert.NoError(t, err)

	router := NewGophermartRouter(
		config,
		registerHandler,
		loginHandler,
		uploadOrderHandler,
		getOrderHandler,
		getBalanceHandler,
		withdrawBalanceHandler,
		withdrawalsHandler,
	)

	ts := httptest.NewServer(router)
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/api/user/register", "application/json", nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = http.Post(ts.URL+"/api/user/login", "application/json", nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = http.Post(ts.URL+"/api/user/orders", "application/json", nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	req, err := http.NewRequest(http.MethodPost, ts.URL+"/api/user/orders", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+tokenString)

	client := &http.Client{}
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	req, err = http.NewRequest(http.MethodGet, ts.URL+"/api/user/orders", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	req, err = http.NewRequest(http.MethodGet, ts.URL+"/api/user/balance", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	req, err = http.NewRequest(http.MethodPost, ts.URL+"/api/user/balance/withdraw", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	req, err = http.NewRequest(http.MethodGet, ts.URL+"/api/user/withdrawals", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
