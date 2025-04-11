package transaction

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestTransaction_Do_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDB := NewMockDB(ctrl)
	mockTx := NewMockTx(ctrl)
	mockDB.EXPECT().
		BeginTx(gomock.Any(), gomock.Nil()).
		Return(mockTx, nil).Times(1)
	mockTx.EXPECT().Commit().Return(nil).Times(1)
	transaction := NewTransaction(mockDB)
	operation := func(tx Tx) error {
		return nil
	}
	err := transaction.Do(context.Background(), operation)
	assert.Nil(t, err)
}

func TestTransaction_Do_BeginTx_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDB := NewMockDB(ctrl)
	mockDB.EXPECT().
		BeginTx(gomock.Any(), gomock.Nil()).
		Return(nil, errors.New("db error")).Times(1)
	transaction := NewTransaction(mockDB)
	operation := func(tx Tx) error {
		return nil
	}
	err := transaction.Do(context.Background(), operation)
	assert.EqualError(t, err, ErrFailedToBeginTransaction.Error())
}

func TestTransaction_Do_Rollback_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDB := NewMockDB(ctrl)
	mockTx := NewMockTx(ctrl)
	mockDB.EXPECT().
		BeginTx(gomock.Any(), gomock.Nil()).
		Return(mockTx, nil).Times(1)
	operation := func(tx Tx) error {
		return errors.New("operation failed")
	}
	mockTx.EXPECT().Rollback().Return(errors.New("rollback error")).Times(1)
	transaction := NewTransaction(mockDB)
	err := transaction.Do(context.Background(), operation)
	assert.EqualError(t, err, ErrFailedToRollbackTransaction.Error())
}

func TestTransaction_Do_Commit_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDB := NewMockDB(ctrl)
	mockTx := NewMockTx(ctrl)
	mockDB.EXPECT().
		BeginTx(gomock.Any(), gomock.Nil()).
		Return(mockTx, nil).Times(1)
	operation := func(tx Tx) error {
		return nil
	}
	mockTx.EXPECT().Commit().Return(errors.New("commit error")).Times(1)
	transaction := NewTransaction(mockDB)
	err := transaction.Do(context.Background(), operation)
	assert.EqualError(t, err, ErrFailedToCommitTransaction.Error())
}
