package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/sbilibin2017/go-gophermart/internal/domain"
	"github.com/sbilibin2017/go-gophermart/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestUserRegisterHandler_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockUserRegisterService(ctrl)
	handler := UserRegisterHandler(mockService)

	tests := []struct {
		name           string
		body           []byte
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "too short login",
			body:           []byte(`{"Login":"ab", "Password":"StrongP@ss1"}`),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   ErrLoginValidation.Error(),
		},
		{
			name:           "invalid login characters",
			body:           []byte(`{"Login":"abc$", "Password":"StrongP@ss1"}`),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   ErrLoginValidation.Error(),
		},
		{
			name:           "password too short",
			body:           []byte(`{"Login":"validlogin", "Password":"Short1!"}`),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   ErrPasswordValidation.Error(),
		},
		{
			name:           "password missing uppercase",
			body:           []byte(`{"Login":"validlogin", "Password":"lowercase1!"}`),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   ErrPasswordValidation.Error(),
		},
		{
			name:           "password missing lowercase",
			body:           []byte(`{"Login":"validlogin", "Password":"UPPERCASE1!"}`),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   ErrPasswordValidation.Error(),
		},
		{
			name:           "password missing digit",
			body:           []byte(`{"Login":"validlogin", "Password":"NoDigits!"}`),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   ErrPasswordValidation.Error(),
		},
		{
			name:           "password missing special char",
			body:           []byte(`{"Login":"validlogin", "Password":"NoSpecial1"}`),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   ErrPasswordValidation.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBuffer(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler(rr, req, httprouter.Params{})

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedBody)
		})
	}
}

func TestUserRegisterHandler_DecodeJSONError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockUserRegisterService(ctrl)
	handler := UserRegisterHandler(mockService)

	body := []byte(`{"Login": "testuser", "Password":`) // malformed JSON
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler(rr, req, httprouter.Params{})

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), ErrUserRegisterInvalidJSONPayload.Error())
}

func TestUserRegisterHandler_CallsRegisterWithCorrectUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockUserRegisterService(ctrl)

	expectedUser := &domain.User{
		Login:    "testuser",
		Password: "ValidPass1!",
	}

	mockService.
		EXPECT().
		Register(gomock.Any(), expectedUser).
		Return("token123", nil)

	handler := UserRegisterHandler(mockService)

	body := []byte(`{"Login":"testuser","Password":"ValidPass1!"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler(rr, req, httprouter.Params{})

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "token123")
}

func TestUserRegisterHandler_HandleRegisterError_UserAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockUserRegisterService(ctrl)

	user := &domain.User{Login: "existing", Password: "Valid1@pass"}

	mockService.
		EXPECT().
		Register(gomock.Any(), user).
		Return("", services.ErrUserAlreadyExists)

	handler := UserRegisterHandler(mockService)

	body := []byte(`{"Login":"existing","Password":"Valid1@pass"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler(rr, req, httprouter.Params{})

	assert.Equal(t, http.StatusConflict, rr.Code)
	assert.Contains(t, rr.Body.String(), services.ErrUserAlreadyExists.Error())
}

func TestUserRegisterHandler_HandleRegisterError_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockUserRegisterService(ctrl)

	user := &domain.User{Login: "someuser", Password: "Valid1@pass"}

	mockService.
		EXPECT().
		Register(gomock.Any(), user).
		Return("", errors.New("db failure"))

	handler := UserRegisterHandler(mockService)

	body := []byte(`{"Login":"someuser","Password":"Valid1@pass"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler(rr, req, httprouter.Params{})

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "db failure")
}
