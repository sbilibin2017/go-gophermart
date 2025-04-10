package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/sbilibin2017/go-gophermart/internal/domain"
	"github.com/sbilibin2017/go-gophermart/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newRequest(t *testing.T, method, url string, body interface{}) *http.Request {
	bodyData, err := json.Marshal(body)
	require.NoError(t, err)
	reqReader := bytes.NewReader(bodyData)
	req, err := http.NewRequest(method, url, reqReader)
	require.NoError(t, err)
	return req
}

func sendRequest(t *testing.T, router *httprouter.Router, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func checkSuccessResponse(t *testing.T, rr *httptest.ResponseRecorder, expectedToken string) {
	assert.Equal(t, http.StatusOK, rr.Code)
	var response UserRegisterResponse
	err := json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, expectedToken, response.AccessToken)
}

func TestUserRegisterHandler_SuccessfulRegistration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := NewMockUserRegisterService(ctrl)
	router := httprouter.New()
	router.POST("/register", UserRegisterHandler(mockService))
	req := UserRegisterRequest{
		Login:    "validLogin123",
		Password: "Valid1Password@",
	}
	expectedToken := &domain.UserToken{Access: "validAccessToken"}
	mockService.EXPECT().Register(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, u *domain.User) (*domain.UserToken, error) {
		assert.Equal(t, req.Login, u.Login)
		assert.Equal(t, req.Password, u.Password)
		return expectedToken, nil
	}).Times(1)
	request := newRequest(t, http.MethodPost, "/register", req)
	rr := sendRequest(t, router, request)
	checkSuccessResponse(t, rr, expectedToken.Access)
}

func TestUserRegisterHandler_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := NewMockUserRegisterService(ctrl)
	router := httprouter.New()
	router.POST("/register", UserRegisterHandler(mockService))
	invalidJSON := `{"login": "validLogin123", "password": "Valid1Password@"`
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader([]byte(invalidJSON)))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

}

func TestUserRegisterHandler_InvalidRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockUserRegisterService(ctrl)
	router := httprouter.New()
	router.POST("/register", UserRegisterHandler(mockService))

	req := UserRegisterRequest{
		Login:    "short",
		Password: "123", // Некорректный пароль
	}

	body, err := json.Marshal(req)
	require.NoError(t, err)

	request := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, request)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

}

func TestUserRegisterHandler_UserAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := NewMockUserRegisterService(ctrl)
	router := httprouter.New()
	router.POST("/register", UserRegisterHandler(mockService))
	req := UserRegisterRequest{
		Login:    "existinguser",
		Password: "Valid1Password@",
	}
	mockService.EXPECT().
		Register(gomock.Any(), gomock.Any()).
		Return(nil, services.ErrUserAlreadyExists).
		Times(1)
	request := newRequest(t, http.MethodPost, "/register", req)
	rr := sendRequest(t, router, request)
	assert.Equal(t, http.StatusConflict, rr.Code)
	assert.Contains(t, rr.Body.String(), services.ErrUserAlreadyExists.Error())
}

func TestHandleUserRegisterError(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		expectedCode int
	}{
		{
			name:         "User already exists",
			err:          services.ErrUserAlreadyExists,
			expectedCode: http.StatusConflict,
		},
		{
			name:         "Internal server error",
			err:          assert.AnError,
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			handleUserRegisterError(rr, tc.err)
			assert.Equal(t, tc.expectedCode, rr.Code)
			assert.Contains(t, rr.Body.String(), tc.err.Error())
		})
	}
}
