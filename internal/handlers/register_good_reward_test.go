package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	e "github.com/sbilibin2017/go-gophermart/internal/errors"
	"github.com/sbilibin2017/go-gophermart/internal/requests"
	"github.com/stretchr/testify/assert"
)

func TestSuccessfulRewardRegistration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRewardService := NewMockRewardService(ctrl)

	handler := RegisterGoodRewardHandler(
		mockRewardService,
		func(w http.ResponseWriter, r *http.Request, v *requests.RewardRequest) error {
			return json.NewDecoder(r.Body).Decode(v)
		},
		func(w http.ResponseWriter, validate *validator.Validate, v interface{}) error {
			return validate.Struct(v)
		},
		func(w http.ResponseWriter, err error, status int) {
			http.Error(w, err.Error(), status)
		},
	)

	mockRewardService.EXPECT().Register(gomock.Any(), gomock.Any()).Return(nil)

	reqBody := `{"match": "match1", "reward": 100, "reward_type": "%"}`
	r := httptest.NewRequest("POST", "/register", strings.NewReader(reqBody))
	recorder := httptest.NewRecorder()

	handler(recorder, r)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "Good reward registered successfully", recorder.Body.String())
}

func TestInvalidRequestBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRewardService := NewMockRewardService(ctrl)

	handler := RegisterGoodRewardHandler(
		mockRewardService,
		func(w http.ResponseWriter, r *http.Request, v *requests.RewardRequest) error {
			return json.NewDecoder(r.Body).Decode(v)
		},
		func(w http.ResponseWriter, validate *validator.Validate, v interface{}) error {
			return validate.Struct(v)
		},
		func(w http.ResponseWriter, err error, status int) {
			http.Error(w, err.Error(), status)
		},
	)

	reqBody := `{"match": "match1", "reward": "invalid", "reward_type": "cash"}`
	r := httptest.NewRequest("POST", "/register", strings.NewReader(reqBody))
	recorder := httptest.NewRecorder()

	handler(recorder, r)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "json: cannot unmarshal string into Go struct field")
}

func TestValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRewardService := NewMockRewardService(ctrl)

	handler := RegisterGoodRewardHandler(
		mockRewardService,
		func(w http.ResponseWriter, r *http.Request, v *requests.RewardRequest) error {
			return json.NewDecoder(r.Body).Decode(v)
		},
		func(w http.ResponseWriter, validate *validator.Validate, v interface{}) error {
			return validate.Struct(v)
		},
		func(w http.ResponseWriter, err error, status int) {
			http.Error(w, err.Error(), status)
		},
	)

	reqBody := `{"match": "match1", "reward": 100}`
	r := httptest.NewRequest("POST", "/register", strings.NewReader(reqBody))
	recorder := httptest.NewRecorder()

	handler(recorder, r)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Key: 'RewardRequest.RewardType' Error:")
}

func TestRewardAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRewardService := NewMockRewardService(ctrl)

	handler := RegisterGoodRewardHandler(
		mockRewardService,
		func(w http.ResponseWriter, r *http.Request, v *requests.RewardRequest) error {
			return json.NewDecoder(r.Body).Decode(v)
		},
		func(w http.ResponseWriter, validate *validator.Validate, v interface{}) error {
			return validate.Struct(v)
		},
		func(w http.ResponseWriter, err error, status int) {
			http.Error(w, err.Error(), status)
		},
	)

	mockRewardService.EXPECT().Register(gomock.Any(), gomock.Any()).Return(e.ErrGoodRewardAlreadyExists)

	reqBody := `{"match": "match1", "reward": 100, "reward_type": "%"}`
	r := httptest.NewRequest("POST", "/register", strings.NewReader(reqBody))
	recorder := httptest.NewRecorder()

	handler(recorder, r)

	assert.Equal(t, http.StatusConflict, recorder.Code)
	assert.Equal(t, e.ErrGoodRewardAlreadyExists.Error(), strings.TrimSpace(recorder.Body.String()))
}

func TestInternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRewardService := NewMockRewardService(ctrl)

	handler := RegisterGoodRewardHandler(
		mockRewardService,
		func(w http.ResponseWriter, r *http.Request, v *requests.RewardRequest) error {
			return json.NewDecoder(r.Body).Decode(v)
		},
		func(w http.ResponseWriter, validate *validator.Validate, v interface{}) error {
			return validate.Struct(v)
		},
		func(w http.ResponseWriter, err error, status int) {
			http.Error(w, err.Error(), status)
		},
	)

	mockRewardService.EXPECT().Register(gomock.Any(), gomock.Any()).Return(errors.New("internal server error"))

	reqBody := `{"match": "match1", "reward": 100, "reward_type": "%"}`
	r := httptest.NewRequest("POST", "/register", strings.NewReader(reqBody))
	recorder := httptest.NewRecorder()

	handler(recorder, r)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Equal(t, "internal server error", strings.TrimSpace(recorder.Body.String()))
}
