package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sbilibin2017/go-gophermart/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestGoodRewardHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := NewMockGoodRewardService(ctrl)
	handler := GoodRewardHandler(mockService)
	body := `{"match":"item-123","reward":10,"reward_type":"pt"}`
	req := httptest.NewRequest(http.MethodPost, "/rewards", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	expectedReq := &types.GoodRewardRegisterRequest{
		Match:      "item-123",
		Reward:     10,
		RewardType: types.RewardTypePoint,
	}
	mockService.EXPECT().
		Register(gomock.Any(), expectedReq).
		Return(&types.APIStatus{
			Status:  http.StatusOK,
			Message: "Good reward registered successfully",
		}, nil, nil)
	handler(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/plain", w.Header().Get("Content-Type"))
	assert.Equal(t, "Good reward registered successfully", w.Body.String())
}

func TestGoodRewardHandler_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := NewMockGoodRewardService(ctrl)
	handler := GoodRewardHandler(mockService)
	body := `{"match":"bad-match","reward":5,"reward_type":"%"}`
	req := httptest.NewRequest(http.MethodPost, "/rewards", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mockService.EXPECT().
		Register(gomock.Any(), &types.GoodRewardRegisterRequest{
			Match:      "bad-match",
			Reward:     5,
			RewardType: types.RewardTypePercent,
		}).
		Return(nil, &types.APIStatus{
			Status:  http.StatusBadRequest,
			Message: "Invalid match pattern",
		}, errors.New("validation failed"))
	handler(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid match pattern")
}

func TestGoodRewardHandler_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := NewMockGoodRewardService(ctrl)
	handler := GoodRewardHandler(mockService)
	body := `{"match": "item-123", "reward": "oops"`
	req := httptest.NewRequest(http.MethodPost, "/rewards", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request body")
}
