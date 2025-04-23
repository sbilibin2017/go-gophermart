package services

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sbilibin2017/go-gophermart/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestGoodRewardService_Register_ValidRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockValidator := NewMockGoodRewardValidator(ctrl)
	mockExistsRepo := NewMockGoodRewardExistsRepository(ctrl)
	mockSaveRepo := NewMockGoodRewardSaveRepository(ctrl)
	svc := NewGoodRewardService(mockValidator, mockExistsRepo, mockSaveRepo)
	mockValidator.EXPECT().Struct(gomock.Any()).Return(nil)
	mockExistsRepo.EXPECT().Exists(gomock.Any(), "match_123").Return(false, nil)
	mockSaveRepo.EXPECT().Save(gomock.Any(), "match_123", int64(100), "type_a").Return(nil)
	req := &types.GoodRewardRegisterRequest{
		Match:      "match_123",
		Reward:     100,
		RewardType: "type_a",
	}
	status, err := svc.Register(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, status.Status)
	assert.Equal(t, RewardRegisteredSuccessMessage, status.Message)
}

func TestGoodRewardService_Register_InvalidRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockValidator := NewMockGoodRewardValidator(ctrl)
	mockExistsRepo := NewMockGoodRewardExistsRepository(ctrl)
	mockSaveRepo := NewMockGoodRewardSaveRepository(ctrl)
	svc := NewGoodRewardService(mockValidator, mockExistsRepo, mockSaveRepo)
	mockValidator.EXPECT().Struct(gomock.Any()).Return(errors.New("invalid struct"))
	req := &types.GoodRewardRegisterRequest{
		Match:      "match_123",
		Reward:     100,
		RewardType: "type_a",
	}
	status, err := svc.Register(context.Background(), req)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, status.Status)
	assert.Equal(t, InvalidRewardRequestMessage, status.Message)
}

func TestGoodRewardService_Register_RewardAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockValidator := NewMockGoodRewardValidator(ctrl)
	mockExistsRepo := NewMockGoodRewardExistsRepository(ctrl)
	mockSaveRepo := NewMockGoodRewardSaveRepository(ctrl)
	svc := NewGoodRewardService(mockValidator, mockExistsRepo, mockSaveRepo)
	mockValidator.EXPECT().Struct(gomock.Any()).Return(nil)
	mockExistsRepo.EXPECT().Exists(gomock.Any(), "match_123").Return(true, nil)
	req := &types.GoodRewardRegisterRequest{
		Match:      "match_123",
		Reward:     100,
		RewardType: "type_a",
	}
	status, err := svc.Register(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusConflict, status.Status)
	assert.Equal(t, RewardExistsMessage, status.Message)
}

func TestGoodRewardService_Register_FailedToSaveReward(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockValidator := NewMockGoodRewardValidator(ctrl)
	mockExistsRepo := NewMockGoodRewardExistsRepository(ctrl)
	mockSaveRepo := NewMockGoodRewardSaveRepository(ctrl)
	svc := NewGoodRewardService(mockValidator, mockExistsRepo, mockSaveRepo)
	mockValidator.EXPECT().Struct(gomock.Any()).Return(nil)
	mockExistsRepo.EXPECT().Exists(gomock.Any(), "match_123").Return(false, nil)
	mockSaveRepo.EXPECT().Save(gomock.Any(), "match_123", int64(100), "type_a").Return(errors.New("failed to save"))
	req := &types.GoodRewardRegisterRequest{
		Match:      "match_123",
		Reward:     100,
		RewardType: "type_a",
	}
	status, err := svc.Register(context.Background(), req)
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, status.Status)
	assert.Equal(t, FailedToSaveRewardMessage, status.Message)
}

func TestGoodRewardService_Register_InternalServerErrorWhenCheckingExistence(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockValidator := NewMockGoodRewardValidator(ctrl)
	mockExistsRepo := NewMockGoodRewardExistsRepository(ctrl)
	mockSaveRepo := NewMockGoodRewardSaveRepository(ctrl)
	svc := NewGoodRewardService(mockValidator, mockExistsRepo, mockSaveRepo)
	mockValidator.EXPECT().Struct(gomock.Any()).Return(nil)
	mockExistsRepo.EXPECT().Exists(gomock.Any(), "match_123").Return(false, errors.New("db error"))
	req := &types.GoodRewardRegisterRequest{
		Match:      "match_123",
		Reward:     100,
		RewardType: "type_a",
	}
	status, err := svc.Register(context.Background(), req)
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, status.Status)
	assert.Equal(t, InternalServerErrorMessage, status.Message)
}
