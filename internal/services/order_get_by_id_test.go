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

func TestOrderGetService_GetByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockOrderRepository(ctrl)
	mockValidator := NewMockValidator(ctrl)

	service := NewOrderGetService(mockValidator, mockRepo)

	req := &types.OrderGetByIDRequest{Number: "12345678903"}
	expectedOrder := &types.OrderDB{
		OrderID: "12345678903",
		Status:  "PROCESSED",
		Accrual: ptrInt64(100),
	}

	mockValidator.
		EXPECT().
		Struct(req).
		Return(nil)

	mockRepo.
		EXPECT().
		GetByID(gomock.Any(), req.Number, []string{"order_id", "status", "accrual"}).
		Return(expectedOrder, nil)

	resp, status, err := service.GetByID(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, status.Status)
	assert.Equal(t, "12345678903", resp.Order)
	assert.Equal(t, types.OrderStatus("PROCESSED"), resp.Status)
	assert.Equal(t, int64(100), *resp.Accrual)
}

func TestOrderGetService_GetByID_InvalidRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockOrderRepository(ctrl)
	mockValidator := NewMockValidator(ctrl)

	service := NewOrderGetService(mockValidator, mockRepo)

	req := &types.OrderGetByIDRequest{Number: ""}

	mockValidator.
		EXPECT().
		Struct(req).
		Return(errors.New("validation failed"))

	resp, status, err := service.GetByID(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, http.StatusBadRequest, status.Status)
	assert.Equal(t, ErrInvalidOrderNumber, status.Message)
}

func TestOrderGetService_GetByID_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockOrderRepository(ctrl)
	mockValidator := NewMockValidator(ctrl)

	service := NewOrderGetService(mockValidator, mockRepo)

	req := &types.OrderGetByIDRequest{Number: "12345678903"}

	mockValidator.
		EXPECT().
		Struct(req).
		Return(nil)

	mockRepo.
		EXPECT().
		GetByID(gomock.Any(), req.Number, []string{"order_id", "status", "accrual"}).
		Return(nil, errors.New("db error"))

	resp, status, err := service.GetByID(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, http.StatusInternalServerError, status.Status)
	assert.Equal(t, ErrInternalServerErrorGet, status.Message)
}

func TestOrderGetService_GetByID_OrderNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockOrderRepository(ctrl)
	mockValidator := NewMockValidator(ctrl)

	service := NewOrderGetService(mockValidator, mockRepo)

	req := &types.OrderGetByIDRequest{Number: "12345678903"}

	mockValidator.
		EXPECT().
		Struct(req).
		Return(nil)

	mockRepo.
		EXPECT().
		GetByID(gomock.Any(), req.Number, []string{"order_id", "status", "accrual"}).
		Return(nil, nil)

	resp, status, err := service.GetByID(context.Background(), req)

	assert.NoError(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, http.StatusNoContent, status.Status)
	assert.Equal(t, ErrOrderNotRegistered, status.Message)
}

func ptrInt64(i int64) *int64 {
	return &i
}
