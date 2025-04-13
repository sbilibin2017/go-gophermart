package services

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sbilibin2017/go-gophermart/internal/repositories"
	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) (
	*gomock.Controller,
	*MockOrderRegisterOrderExistsRepository,
	*MockOrderRegisterOrderSaveRepository,
	*MockOrderRegisterRewardFilterRepository,
	*MockTransaction,
	*OrderRegisterService,
	*Order,
) {
	ctrl := gomock.NewController(t)
	mockOrderExistsRepo := NewMockOrderRegisterOrderExistsRepository(ctrl)
	mockOrderSaveRepo := NewMockOrderRegisterOrderSaveRepository(ctrl)
	mockRewardFilterRepo := NewMockOrderRegisterRewardFilterRepository(ctrl)
	mockUnitOfWork := NewMockTransaction(ctrl)

	service := NewOrderRegisterService(
		mockOrderExistsRepo,
		mockOrderSaveRepo,
		mockRewardFilterRepo,
		mockUnitOfWork,
	)

	order := &Order{
		Order: 123,
		Goods: []Good{
			{Description: "Item 1", Price: 100},
			{Description: "Item 2", Price: 200},
		},
	}

	return ctrl, mockOrderExistsRepo, mockOrderSaveRepo, mockRewardFilterRepo, mockUnitOfWork, service, order
}

func TestOrderRegisterService_Register_HappyPath(t *testing.T) {
	ctrl, mockOrderExistsRepo, mockOrderSaveRepo, mockRewardFilterRepo, mockUnitOfWork, service, order := setup(t)
	defer ctrl.Finish()

	mockOrderExistsRepo.EXPECT().Exists(gomock.Any(), gomock.Eq(&repositories.OrderExistsID{OrderID: order.Order})).Return(false, nil)

	mockRewardFilterRepo.EXPECT().Filter(gomock.Any(), gomock.Eq([]*repositories.RewardFilter{
		{Description: "Item 1"},
		{Description: "Item 2"},
	})).Return([]*repositories.RewardFiltered{
		{RewardType: "%", Reward: 10},
		{RewardType: "pt", Reward: 50},
	}, nil)

	mockOrderSaveRepo.EXPECT().Save(gomock.Any(), gomock.Eq(&repositories.OrderSave{
		OrderID: order.Order,
		Status:  repositories.StatusNew,
		Accrual: 60,
	})).Return(nil)

	mockUnitOfWork.EXPECT().Do(gomock.Any(), gomock.Any()).Do(func(ctx context.Context, op func(tx *sql.Tx) error) error {
		return op(nil)
	}).Return(nil)

	err := service.Register(context.Background(), order)

	assert.NoError(t, err)
}

func TestOrderRegisterService_Register_OrderAlreadyRegistered(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockExistsRepo := NewMockOrderRegisterOrderExistsRepository(ctrl)
	mockSaveRepo := NewMockOrderRegisterOrderSaveRepository(ctrl)
	mockFilterRepo := NewMockOrderRegisterRewardFilterRepository(ctrl)
	mockUnitOfWork := NewMockTransaction(ctrl)

	order := &Order{Order: 123, Goods: []Good{{Description: "Item1", Price: 100}}}

	mockExistsRepo.EXPECT().Exists(gomock.Any(), gomock.Eq(&repositories.OrderExistsID{OrderID: 123})).Return(true, nil).Times(1)

	mockUnitOfWork.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, operation func(tx *sql.Tx) error) error {
		return operation(nil)
	}).Times(1)

	service := &OrderRegisterService{
		oer: mockExistsRepo,
		osr: mockSaveRepo,
		rga: mockFilterRepo,
		tx:  mockUnitOfWork,
	}

	err := service.Register(context.Background(), order)
	assert.Error(t, err)
	assert.Equal(t, ErrOrderAlreadyRegistered, err)
}

func TestOrderRegisterService_Register_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockExistsRepo := NewMockOrderRegisterOrderExistsRepository(ctrl)
	mockSaveRepo := NewMockOrderRegisterOrderSaveRepository(ctrl)
	mockFilterRepo := NewMockOrderRegisterRewardFilterRepository(ctrl)
	mockUnitOfWork := NewMockTransaction(ctrl)

	order := &Order{Order: 123, Goods: []Good{{Description: "Item1", Price: 100}}}

	mockExistsRepo.EXPECT().Exists(gomock.Any(), gomock.Eq(&repositories.OrderExistsID{OrderID: 123})).Return(false, errors.New("some error")).Times(1)

	mockUnitOfWork.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, operation func(tx *sql.Tx) error) error {
		return operation(nil)
	}).Times(1)

	service := &OrderRegisterService{
		oer: mockExistsRepo,
		osr: mockSaveRepo,
		rga: mockFilterRepo,
		tx:  mockUnitOfWork,
	}

	err := service.Register(context.Background(), order)

	assert.Error(t, err)
	assert.Equal(t, ErrOrderRegisterInternal, err)
}

func TestOrderRegisterService_Register_FilterError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockExistsRepo := NewMockOrderRegisterOrderExistsRepository(ctrl)
	mockSaveRepo := NewMockOrderRegisterOrderSaveRepository(ctrl)
	mockFilterRepo := NewMockOrderRegisterRewardFilterRepository(ctrl)
	mockUnitOfWork := NewMockTransaction(ctrl)

	order := &Order{
		Order: 123,
		Goods: []Good{{Description: "Item1", Price: 100}},
	}

	mockExistsRepo.EXPECT().Exists(gomock.Any(), gomock.Eq(&repositories.OrderExistsID{OrderID: 123})).Return(false, nil).Times(1)

	mockFilterRepo.EXPECT().Filter(gomock.Any(), gomock.Any()).Return(nil, errors.New("filter error")).Times(1)

	mockUnitOfWork.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, operation func(tx *sql.Tx) error) error {
		return operation(nil)
	}).Times(1)

	service := &OrderRegisterService{
		oer: mockExistsRepo,
		osr: mockSaveRepo,
		rga: mockFilterRepo,
		tx:  mockUnitOfWork,
	}

	err := service.Register(context.Background(), order)

	assert.Error(t, err)
	assert.Equal(t, ErrOrderRegisterInternal, err)
}

func TestOrderRegisterService_Register_SaveError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockExistsRepo := NewMockOrderRegisterOrderExistsRepository(ctrl)
	mockSaveRepo := NewMockOrderRegisterOrderSaveRepository(ctrl)
	mockFilterRepo := NewMockOrderRegisterRewardFilterRepository(ctrl)
	mockUnitOfWork := NewMockTransaction(ctrl)

	order := &Order{
		Order: 123,
		Goods: []Good{{Description: "Item1", Price: 100}},
	}

	mockExistsRepo.EXPECT().Exists(gomock.Any(), gomock.Eq(&repositories.OrderExistsID{OrderID: 123})).Return(false, nil).Times(1)

	mockFilterRepo.EXPECT().Filter(gomock.Any(), gomock.Any()).Return([]*repositories.RewardFiltered{
		{Reward: 10, RewardType: "%"},
	}, nil).Times(1)

	mockUnitOfWork.EXPECT().Do(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, operation func(tx *sql.Tx) error) error {
		return operation(nil)
	}).Times(1)

	mockSaveRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(errors.New("save error")).Times(1)

	service := &OrderRegisterService{
		oer: mockExistsRepo,
		osr: mockSaveRepo,
		rga: mockFilterRepo,
		tx:  mockUnitOfWork,
	}

	err := service.Register(context.Background(), order)

	assert.Error(t, err)
	assert.Equal(t, ErrOrderRegisterInternal, err)
}
