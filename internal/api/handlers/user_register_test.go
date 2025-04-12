package handlers_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sbilibin2017/go-gophermart/internal/api/handlers"
	"github.com/sbilibin2017/go-gophermart/internal/services"
	"github.com/sbilibin2017/go-gophermart/internal/usecases"
	"github.com/sbilibin2017/go-gophermart/internal/usecases/validators"
	"github.com/stretchr/testify/assert"
)

func TestUserRegisterHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := handlers.NewMockUserRegisterUsecase(ctrl)
	mockDecoder := handlers.NewMockDecoder(ctrl)

	reqBody := `{"login": "user1", "password": "password123"}`
	httpReq := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(reqBody))
	httpRes := httptest.NewRecorder()

	expectedReq := &usecases.UserRegisterRequest{
		Login:    "user1",
		Password: "password123",
	}

	mockDecoder.EXPECT().
		Decode(httpReq, gomock.Any()).
		DoAndReturn(func(r *http.Request, v any) error {
			*v.(*usecases.UserRegisterRequest) = *expectedReq
			return nil
		})

	mockUC.EXPECT().
		Execute(gomock.Any(), expectedReq).
		Return(&usecases.UserRegisterResponse{AccessToken: "token123"}, nil)

	handler := handlers.UserRegisterHandler(mockUC, mockDecoder)
	handler.ServeHTTP(httpRes, httpReq)

	assert.Equal(t, http.StatusOK, httpRes.Code)
	assert.Equal(t, "Bearer token123", httpRes.Header().Get("Authorization"))
}

func TestUserRegisterHandler_DecodeError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := handlers.NewMockUserRegisterUsecase(ctrl)
	mockDecoder := handlers.NewMockDecoder(ctrl)

	httpReq := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(`invalid-json`))
	httpRes := httptest.NewRecorder()

	mockDecoder.EXPECT().
		Decode(httpReq, gomock.Any()).
		Return(errors.New("decode error"))

	handler := handlers.UserRegisterHandler(mockUC, mockDecoder)
	handler.ServeHTTP(httpRes, httpReq)

	assert.Equal(t, http.StatusBadRequest, httpRes.Code)
	assert.Contains(t, httpRes.Body.String(), "unprocessable json")
}

func TestUserRegisterHandler_UserAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := handlers.NewMockUserRegisterUsecase(ctrl)
	mockDecoder := handlers.NewMockDecoder(ctrl)

	req := &usecases.UserRegisterRequest{Login: "user1", Password: "pass"}
	httpReq := httptest.NewRequest(http.MethodPost, "/register", nil)
	httpRes := httptest.NewRecorder()

	mockDecoder.EXPECT().
		Decode(httpReq, gomock.Any()).
		DoAndReturn(func(r *http.Request, v any) error {
			*v.(*usecases.UserRegisterRequest) = *req
			return nil
		})

	mockUC.EXPECT().
		Execute(gomock.Any(), req).
		Return(nil, services.ErrUserAlreadyExists)

	handler := handlers.UserRegisterHandler(mockUC, mockDecoder)
	handler.ServeHTTP(httpRes, httpReq)

	assert.Equal(t, http.StatusConflict, httpRes.Code)
	assert.Equal(t, services.ErrUserAlreadyExists.Error()+"\n", httpRes.Body.String())
}

func TestUserRegisterHandler_InvalidInput(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := handlers.NewMockUserRegisterUsecase(ctrl)
	mockDecoder := handlers.NewMockDecoder(ctrl)

	req := &usecases.UserRegisterRequest{Login: "bad login", Password: "123"}
	httpReq := httptest.NewRequest(http.MethodPost, "/register", nil)
	httpRes := httptest.NewRecorder()

	mockDecoder.EXPECT().
		Decode(httpReq, gomock.Any()).
		DoAndReturn(func(r *http.Request, v any) error {
			*v.(*usecases.UserRegisterRequest) = *req
			return nil
		})

	mockUC.EXPECT().
		Execute(gomock.Any(), req).
		Return(nil, validators.ErrInvalidLogin)

	handler := handlers.UserRegisterHandler(mockUC, mockDecoder)
	handler.ServeHTTP(httpRes, httpReq)

	assert.Equal(t, http.StatusBadRequest, httpRes.Code)
	assert.Equal(t, validators.ErrInvalidLogin.Error()+"\n", httpRes.Body.String())
}

func TestUserRegisterHandler_InternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := handlers.NewMockUserRegisterUsecase(ctrl)
	mockDecoder := handlers.NewMockDecoder(ctrl)

	req := &usecases.UserRegisterRequest{Login: "user1", Password: "password123"}
	httpReq := httptest.NewRequest(http.MethodPost, "/register", nil)
	httpRes := httptest.NewRecorder()

	mockDecoder.EXPECT().
		Decode(httpReq, gomock.Any()).
		DoAndReturn(func(r *http.Request, v any) error {
			*v.(*usecases.UserRegisterRequest) = *req
			return nil
		})

	mockUC.EXPECT().
		Execute(gomock.Any(), req).
		Return(nil, errors.New("some unexpected error"))

	handler := handlers.UserRegisterHandler(mockUC, mockDecoder)
	handler.ServeHTTP(httpRes, httpReq)

	assert.Equal(t, http.StatusInternalServerError, httpRes.Code)
	assert.Equal(t, "internal error\n", httpRes.Body.String())
}
