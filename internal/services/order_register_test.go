package services

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	gomock "github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/sbilibin2017/go-gophermart/internal/repositories"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func init() {
	logger.Init(zapcore.InfoLevel)
}

func setupOrderRegisterTest(t *testing.T) (
	*gomock.Controller,
	sqlmock.Sqlmock,
	*sql.DB,
	*OrderRegisterService,
	*MockOrderRegisterOrderExistsRepository,
	*MockOrderRegisterOrderSaveRepository,
	*MockOrderRegisterRewardFilterRepository,
) {
	ctrl := gomock.NewController(t)

	mockOER := NewMockOrderRegisterOrderExistsRepository(ctrl)
	mockOSR := NewMockOrderRegisterOrderSaveRepository(ctrl)
	mockRGA := NewMockOrderRegisterRewardFilterRepository(ctrl)

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "pgx")

	svc := NewOrderRegisterService(mockOER, mockOSR, mockRGA, sqlxDB)

	return ctrl, mock, db, svc, mockOER, mockOSR, mockRGA
}

func TestOrderAlreadyRegistered(t *testing.T) {
	ctrl, mock, db, svc, mockOER, _, _ := setupOrderRegisterTest(t)
	defer ctrl.Finish()
	defer db.Close()

	mockOER.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(true, nil)

	mock.ExpectBegin()
	mock.ExpectRollback()

	err := svc.Register(context.Background(), &Order{
		Number: 12345,
	})

	assert.Equal(t, ErrOrderAlreadyRegistered, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRegister_Success(t *testing.T) {
	ctrl, mock, db, svc, mockOER, mockOSR, mockRGA := setupOrderRegisterTest(t)
	defer ctrl.Finish()
	defer db.Close()

	mockOER.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(false, nil)
	mockRGA.EXPECT().Filter(gomock.Any(), gomock.Any()).Return(&repositories.RewardFilteredDB{Reward: 10.0}, nil)
	mockOSR.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)

	mock.ExpectBegin()
	mock.ExpectCommit()

	order := &Order{
		Number: 12345,
		Goods: []Good{
			{Description: "item1", Price: 100},
		},
	}

	err := svc.Register(context.Background(), order)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRegister_FailedToSaveOrder(t *testing.T) {
	ctrl, mock, db, svc, mockOER, mockOSR, mockRGA := setupOrderRegisterTest(t)
	defer ctrl.Finish()
	defer db.Close()

	mockOER.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(false, nil)
	mockRGA.EXPECT().Filter(gomock.Any(), gomock.Any()).Return(&repositories.RewardFilteredDB{Reward: 10.0}, nil)
	mockOSR.EXPECT().Save(gomock.Any(), gomock.Any()).Return(errors.New("failed to save order"))

	mock.ExpectBegin()
	mock.ExpectRollback()

	order := &Order{
		Number: 12345,
		Goods: []Good{
			{Description: "item1", Price: 100},
		},
	}

	err := svc.Register(context.Background(), order)

	assert.Equal(t, ErrOrderRegisterDBInternal, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRegister_FailedToCheckOrderExistence(t *testing.T) {
	ctrl, mock, db, svc, mockOER, _, _ := setupOrderRegisterTest(t)
	defer ctrl.Finish()
	defer db.Close()

	mockOER.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(false, errors.New("db error"))

	mock.ExpectBegin()
	mock.ExpectRollback()

	order := &Order{
		Number: 12345,
		Goods: []Good{
			{Description: "item1", Price: 100},
		},
	}

	err := svc.Register(context.Background(), order)

	assert.Equal(t, ErrOrderRegisterDBInternal, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestApplyReward(t *testing.T) {
	tests := []struct {
		name        string
		reward      *repositories.RewardFilteredDB
		good        Good
		initialAcc  float64
		expectedAcc float64
	}{
		{
			name: "Percentage reward",
			reward: &repositories.RewardFilteredDB{
				RewardType: "%",
				Reward:     10.0,
			},
			good: Good{
				Description: "item1",
				Price:       100,
			},
			initialAcc:  0,
			expectedAcc: 10.0,
		},
		{
			name: "Points reward",
			reward: &repositories.RewardFilteredDB{
				RewardType: "pt",
				Reward:     20.0,
			},
			good: Good{
				Description: "item2",
				Price:       100,
			},
			initialAcc:  0,
			expectedAcc: 20.0,
		},
		{
			name: "Percentage reward with initial accrual",
			reward: &repositories.RewardFilteredDB{
				RewardType: "%",
				Reward:     5.0,
			},
			good: Good{
				Description: "item3",
				Price:       200,
			},
			initialAcc:  10.0,
			expectedAcc: 20.0,
		},
		{
			name: "Points reward with initial accrual",
			reward: &repositories.RewardFilteredDB{
				RewardType: "pt",
				Reward:     15.0,
			},
			good: Good{
				Description: "item4",
				Price:       150,
			},
			initialAcc:  5.0,
			expectedAcc: 20.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc := applyReward(tt.reward, tt.good, tt.initialAcc)

			assert.Equal(t, tt.expectedAcc, acc)
		})
	}
}

func TestRegister_FailedToGetRewardForGood(t *testing.T) {
	ctrl, mock, db, svc, mockOER, _, mockRGA := setupOrderRegisterTest(t)
	defer ctrl.Finish()
	defer db.Close()

	mockOER.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(false, nil)

	mockRGA.EXPECT().Filter(gomock.Any(), gomock.Any()).Return(nil, errors.New("db error"))

	mock.ExpectBegin()
	mock.ExpectRollback()

	order := &Order{
		Number: 12345,
		Goods: []Good{
			{Description: "item1", Price: 100},
		},
	}

	err := svc.Register(context.Background(), order)

	assert.Equal(t, ErrOrderRegisterDBInternal, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
