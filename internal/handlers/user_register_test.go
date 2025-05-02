package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sbilibin2017/go-gophermart/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRegisterHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockUserRegisterService(ctrl)
	mockService.EXPECT().
		Register(gomock.Any(), gomock.Any()).
		Return("mockToken", &types.APIStatus{Message: "Success", StatusCode: http.StatusOK}, nil).
		Times(1)

	reqBody := &types.UserRegisterRequest{
		Login:    "testuser",
		Password: "password123",
	}

	reqBodyBytes, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqBodyBytes))
	rr := httptest.NewRecorder()

	handler := UserRegisterHandler(mockService)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Header().Get("Authorization"), "mockToken")
	assert.Equal(t, "Success", rr.Body.String())
}

func TestUserRegisterHandler_ErrorOnRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockUserRegisterService(ctrl)
	mockService.EXPECT().
		Register(gomock.Any(), gomock.Any()).
		Return("", nil, &types.APIStatus{Message: "Registration failed", StatusCode: http.StatusBadRequest}).
		Times(1)

	reqBody := &types.UserRegisterRequest{
		Login:    "testuser",
		Password: "password123",
	}

	reqBodyBytes, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqBodyBytes))
	rr := httptest.NewRecorder()

	handler := UserRegisterHandler(mockService)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "Registration failed", rr.Body.String())
}

func TestUserRegisterHandler_DecodeRequestError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockUserRegisterService(ctrl)

	// No actual call to Register will happen as we simulate a request decoding error
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader([]byte("invalid data")))
	rr := httptest.NewRecorder()

	handler := UserRegisterHandler(mockService)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUserRegisterHandler_RegisterError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockUserRegisterService(ctrl)
	mockService.EXPECT().
		Register(gomock.Any(), gomock.Any()).
		Return("", nil, &types.APIStatus{Message: "Service error", StatusCode: http.StatusInternalServerError}).
		Times(1)

	reqBody := &types.UserRegisterRequest{
		Login:    "testuser",
		Password: "password123",
	}

	reqBodyBytes, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqBodyBytes))
	rr := httptest.NewRecorder()

	handler := UserRegisterHandler(mockService)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "Service error", rr.Body.String())
}
