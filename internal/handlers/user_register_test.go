package handlers

import (
	"bytes"
	"encoding/json"
	e "errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sbilibin2017/go-gophermart/internal/errors"
	"github.com/sbilibin2017/go-gophermart/internal/usecases"
	"github.com/stretchr/testify/assert"
)

func TestUserRegisterHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUsecase := NewMockUserRegisterUsecase(ctrl)
	mockDecoder := NewMockDecoder(ctrl)
	reqBody := &usecases.UserRegisterRequest{
		Login:    "validuser",
		Password: "ValidPassword123!",
	}
	mockResponse := &usecases.UserRegisterResponse{
		AccessToken: "valid-access-token",
	}
	mockUsecase.EXPECT().Execute(gomock.Any(), reqBody).Return(mockResponse, nil).Times(1)
	mockDecoder.EXPECT().Decode(gomock.Any()).DoAndReturn(func(v interface{}) error {
		*v.(*usecases.UserRegisterRequest) = *reqBody
		return nil
	}).Times(1)
	reqBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqBytes))
	recorder := httptest.NewRecorder()
	handler := UserRegisterHandler(mockUsecase, mockDecoder)
	handler.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "Bearer valid-access-token", recorder.Header().Get("Authorization"))
}

func TestUserRegisterHandler_DecodeError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUsecase := NewMockUserRegisterUsecase(ctrl)
	mockDecoder := NewMockDecoder(ctrl)
	invalidReqBody := `{ "Login": "validuser", "Password": }`
	mockDecoder.EXPECT().Decode(gomock.Any()).Return(e.New("Invalid JSON format")).Times(1)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader([]byte(invalidReqBody)))
	recorder := httptest.NewRecorder()
	handler := UserRegisterHandler(mockUsecase, mockDecoder)
	handler.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Equal(t, strings.TrimSpace(errors.ErrUnprocessableJson.Error()), strings.TrimSpace(recorder.Body.String()))
}

func TestUserRegisterHandler_InvalidJSONBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDecoder := NewMockDecoder(ctrl)
	mockDecoder.EXPECT().Decode(gomock.Any()).Return(e.New("invalid JSON")).Times(1)
	mockUC := NewMockUserRegisterUsecase(ctrl)
	reqBody := []byte(`{"login": "user", "password": "pass"}`)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBody))
	w := httptest.NewRecorder()
	handler := UserRegisterHandler(mockUC, mockDecoder)
	handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserRegisterHandler_HandleError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := NewMockUserRegisterUsecase(ctrl)
	mockDecoder := NewMockDecoder(ctrl)
	tests := []struct {
		name         string
		mockExecute  func()
		mockDecode   func()
		expectedCode int
		expectedBody string
	}{
		{
			name: "UserAlreadyExists",
			mockExecute: func() {
				mockUC.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil, errors.ErrUserAlreadyExists).Times(1)
			},
			mockDecode: func() {
				mockDecoder.EXPECT().Decode(gomock.Any()).Return(nil).Times(1)
			},
			expectedCode: http.StatusConflict,
			expectedBody: errors.ErrUserAlreadyExists.Error(),
		},
		{
			name: "InvalidLogin",
			mockExecute: func() {
				mockUC.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil, errors.ErrInvalidLogin).Times(1)
			},
			mockDecode: func() {
				mockDecoder.EXPECT().Decode(gomock.Any()).Return(nil).Times(1)
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: errors.ErrInvalidLogin.Error(),
		},
		{
			name: "InvalidPassword",
			mockExecute: func() {
				mockUC.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil, errors.ErrInvalidPassword).Times(1)
			},
			mockDecode: func() {
				mockDecoder.EXPECT().Decode(gomock.Any()).Return(nil).Times(1)
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: errors.ErrInvalidPassword.Error(),
		},
		{
			name: "InternalError",
			mockExecute: func() {
				mockUC.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil, errors.ErrInternal).Times(1)
			},
			mockDecode: func() {
				mockDecoder.EXPECT().Decode(gomock.Any()).Return(nil).Times(1)
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: errors.ErrInternal.Error(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockDecode()
			tt.mockExecute()
			reqBody := []byte(`{"login": "user", "password": "pass"}`)
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBody))
			w := httptest.NewRecorder()
			handler := UserRegisterHandler(mockUC, mockDecoder)
			handler.ServeHTTP(w, req)
			actualBody := strings.TrimSpace(w.Body.String())
			assert.Equal(t, tt.expectedBody, actualBody)
		})
	}
}
