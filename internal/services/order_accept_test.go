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

func TestOrderAcceptService_Accept_ExistsReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockValidator := NewMockOrderAcceptValidator(ctrl)
	mockOrderExistsRepo := NewMockOrderAcceptOrderExistsByIDRepository(ctrl)
	mockOrderSaveRepo := NewMockOrderAcceptOrderSaveRepository(ctrl)
	mockRewardFilterRepo := NewMockOrderAcceptGoodRewardFilterILikeRepository(ctrl)

	svc := NewOrderAcceptService(mockValidator, mockOrderExistsRepo, mockOrderSaveRepo, mockRewardFilterRepo)

	req := &types.OrderAcceptRequest{
		Order: "123456",
		Goods: []types.Good{
			{Description: "test good", Price: 100},
		},
	}

	mockValidator.EXPECT().Struct(req).Return(nil).Times(1)
	mockOrderExistsRepo.EXPECT().Exists(context.Background(), req.Order).Return(false, assert.AnError).Times(1)

	resp, err := svc.Accept(context.Background(), req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.Status)
	assert.Equal(t, ErrInternalServerError, resp.Message)
}

func TestOrderAcceptService_Accept_SaveError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockValidator := NewMockOrderAcceptValidator(ctrl)
	mockOrderExistsRepo := NewMockOrderAcceptOrderExistsByIDRepository(ctrl)
	mockOrderSaveRepo := NewMockOrderAcceptOrderSaveRepository(ctrl)
	mockRewardFilterRepo := NewMockOrderAcceptGoodRewardFilterILikeRepository(ctrl)

	svc := NewOrderAcceptService(mockValidator, mockOrderExistsRepo, mockOrderSaveRepo, mockRewardFilterRepo)

	req := &types.OrderAcceptRequest{
		Order: "123456",
		Goods: []types.Good{
			{Description: "test good", Price: 100},
		},
	}

	mockValidator.EXPECT().Struct(req).Return(nil).Times(1)
	mockOrderExistsRepo.EXPECT().Exists(context.Background(), req.Order).Return(false, nil).Times(1)
	mockRewardFilterRepo.EXPECT().FilterILike(context.Background(), "test good", []string{"reward_type", "reward"}).Return(&types.RewardDB{
		Reward:     10,
		RewardType: types.RewardTypePercent,
	}, nil).Times(1)
	mockOrderSaveRepo.EXPECT().Save(context.Background(), req.Order, string(types.OrderStatusRegistered), int64(10)).Return(assert.AnError).Times(1)

	resp, err := svc.Accept(context.Background(), req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.Status)
	assert.Equal(t, ErrFailedToSaveOrder, resp.Message)
}

func TestOrderAcceptService_Accept_InvalidRewardRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockValidator := NewMockOrderAcceptValidator(ctrl)
	mockOrderExistsRepo := NewMockOrderAcceptOrderExistsByIDRepository(ctrl)
	mockOrderSaveRepo := NewMockOrderAcceptOrderSaveRepository(ctrl)
	mockRewardFilterRepo := NewMockOrderAcceptGoodRewardFilterILikeRepository(ctrl)
	svc := NewOrderAcceptService(mockValidator, mockOrderExistsRepo, mockOrderSaveRepo, mockRewardFilterRepo)
	req := &types.OrderAcceptRequest{
		Order: "123456",
		Goods: []types.Good{
			{Description: "test good", Price: 100},
		},
	}
	mockValidator.EXPECT().Struct(req).Return(errors.New("invalid request")).Times(1)
	mockOrderExistsRepo.EXPECT().Exists(gomock.Any(), gomock.Any()).Times(0)
	resp, err := svc.Accept(context.Background(), req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.Status)
	assert.Equal(t, ErrInvalidRewardRequest, resp.Message)
}

func TestOrderAcceptService_Accept_OrderAlreadyProcessed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockValidator := NewMockOrderAcceptValidator(ctrl)
	mockOrderExistsRepo := NewMockOrderAcceptOrderExistsByIDRepository(ctrl)
	mockOrderSaveRepo := NewMockOrderAcceptOrderSaveRepository(ctrl)
	mockRewardFilterRepo := NewMockOrderAcceptGoodRewardFilterILikeRepository(ctrl)
	svc := NewOrderAcceptService(mockValidator, mockOrderExistsRepo, mockOrderSaveRepo, mockRewardFilterRepo)
	req := &types.OrderAcceptRequest{
		Order: "123456",
		Goods: []types.Good{
			{Description: "test good", Price: 100},
		},
	}
	mockValidator.EXPECT().Struct(req).Return(nil).Times(1)
	mockOrderExistsRepo.EXPECT().Exists(context.Background(), req.Order).Return(true, nil).Times(1)
	resp, err := svc.Accept(context.Background(), req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusConflict, resp.Status)
	assert.Equal(t, ErrOrderAlreadyProcessed, resp.Message)
}

func TestOrderAcceptService_Accept_InternalServerErrorFilter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockValidator := NewMockOrderAcceptValidator(ctrl)
	mockOrderExistsRepo := NewMockOrderAcceptOrderExistsByIDRepository(ctrl)
	mockOrderSaveRepo := NewMockOrderAcceptOrderSaveRepository(ctrl)
	mockRewardFilterRepo := NewMockOrderAcceptGoodRewardFilterILikeRepository(ctrl)
	svc := NewOrderAcceptService(mockValidator, mockOrderExistsRepo, mockOrderSaveRepo, mockRewardFilterRepo)
	req := &types.OrderAcceptRequest{
		Order: "123456",
		Goods: []types.Good{
			{Description: "test good", Price: 100},
		},
	}
	mockValidator.EXPECT().Struct(req).Return(nil).Times(1)
	mockOrderExistsRepo.EXPECT().Exists(context.Background(), req.Order).Return(false, nil).Times(1)
	mockRewardFilterRepo.EXPECT().FilterILike(context.Background(), "test good", []string{"reward_type", "reward"}).Return(nil, assert.AnError).Times(1)
	resp, err := svc.Accept(context.Background(), req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.Status)
	assert.Equal(t, ErrInternalServerErrorFilter, resp.Message)
}

func TestOrderAcceptService_Accept_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockValidator := NewMockOrderAcceptValidator(ctrl)
	mockOrderExistsRepo := NewMockOrderAcceptOrderExistsByIDRepository(ctrl)
	mockOrderSaveRepo := NewMockOrderAcceptOrderSaveRepository(ctrl)
	mockRewardFilterRepo := NewMockOrderAcceptGoodRewardFilterILikeRepository(ctrl)
	svc := NewOrderAcceptService(mockValidator, mockOrderExistsRepo, mockOrderSaveRepo, mockRewardFilterRepo)
	req := &types.OrderAcceptRequest{
		Order: "123456",
		Goods: []types.Good{
			{Description: "test good", Price: 100},
		},
	}
	mockValidator.EXPECT().Struct(req).Return(nil).Times(1)
	mockOrderExistsRepo.EXPECT().Exists(context.Background(), req.Order).Return(false, nil).Times(1)
	mockRewardFilterRepo.EXPECT().FilterILike(context.Background(), "test good", []string{"reward_type", "reward"}).Return(&types.RewardDB{
		Reward:     10,
		RewardType: types.RewardTypePercent,
	}, nil).Times(1)
	mockOrderSaveRepo.EXPECT().Save(context.Background(), req.Order, string(types.OrderStatusRegistered), int64(10)).Return(nil).Times(1)
	resp, err := svc.Accept(context.Background(), req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusAccepted, resp.Status)
	assert.Equal(t, SuccessOrderAcceptedForProcessing, resp.Message)
}

func TestCalcAccrual(t *testing.T) {
	tests := []struct {
		price      int64
		reward     int64
		rewardType types.RewardType
		expected   int64
	}{
		{1000, 10, types.RewardTypePercent, 100},   // 10% от 1000 = 100
		{1000, 10, types.RewardTypePoint, 10},      // В случае типа балла возвращается сама сумма награды
		{1000, 10, types.RewardType("unknown"), 0}, // Для неизвестного типа награды должно быть 0
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			result := calcAccrual(test.price, test.reward, test.rewardType)
			assert.Equal(t, test.expected, result)
		})
	}
}
