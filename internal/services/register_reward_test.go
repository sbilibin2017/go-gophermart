package services

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	gomock "github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestRewardService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRewardExistsRepo := NewMockRewardExistsRepository(ctrl)
	mockRewardSaveRepo := NewMockRewardSaveRepository(ctrl)

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	service := NewRewardService(mockRewardExistsRepo, mockRewardSaveRepo, sqlxDB)

	reward := &domain.Reward{
		Match:      "product-123",
		Reward:     100,
		RewardType: domain.RewardTypePercent,
	}

	tests := []struct {
		name              string
		mockExistsReturn  bool
		mockExistsError   error
		mockSaveError     error
		expectedError     error
		mockExpectActions func()
	}{
		{
			name:             "successful registration",
			mockExistsReturn: false,
			mockExistsError:  nil,
			mockSaveError:    nil,
			expectedError:    nil,
			mockExpectActions: func() {
				mock.ExpectBegin()
				mockRewardExistsRepo.EXPECT().Exists(context.Background(), gomock.Any(), gomock.Any()).Return(false, nil).Times(1)
				mockRewardSaveRepo.EXPECT().Save(context.Background(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
				mock.ExpectCommit()
			},
		},
		{
			name:             "reward already exists",
			mockExistsReturn: true,
			mockExistsError:  nil,
			mockSaveError:    nil,
			expectedError:    ErrGoodRewardAlreadyExists,
			mockExpectActions: func() {
				mock.ExpectBegin()
				mockRewardExistsRepo.EXPECT().Exists(context.Background(), gomock.Any(), gomock.Any()).Return(true, nil).Times(1)
				mock.ExpectRollback()
			},
		},
		{
			name:             "database error on save",
			mockExistsReturn: false,
			mockExistsError:  nil,
			mockSaveError:    errors.New("db error"),
			expectedError:    ErrRegisterRewardInternal,
			mockExpectActions: func() {
				mock.ExpectBegin()
				mockRewardExistsRepo.EXPECT().Exists(context.Background(), gomock.Any(), gomock.Any()).Return(false, nil).Times(1)
				mockRewardSaveRepo.EXPECT().Save(context.Background(), gomock.Any(), gomock.Any()).Return(errors.New("db error")).Times(1)
				mock.ExpectRollback()
			},
		},
		{
			name:             "error in transaction begin",
			mockExistsReturn: false,
			mockExistsError:  nil,
			mockSaveError:    nil,
			expectedError:    errors.New("transaction begin error"),
			mockExpectActions: func() {
				mock.ExpectBegin().WillReturnError(errors.New("transaction begin error"))
			},
		},
		{
			name:             "error during transaction commit",
			mockExistsReturn: false,
			mockExistsError:  nil,
			mockSaveError:    nil,
			expectedError:    errors.New("commit error"),
			mockExpectActions: func() {
				mock.ExpectBegin()
				mockRewardExistsRepo.EXPECT().Exists(context.Background(), gomock.Any(), gomock.Any()).Return(false, nil).Times(1)
				mockRewardSaveRepo.EXPECT().Save(context.Background(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
				mock.ExpectCommit().WillReturnError(errors.New("commit error"))
			},
		},
		{
			name:             "error in Exists method",
			mockExistsReturn: false,
			mockExistsError:  errors.New("exists method error"),
			mockSaveError:    nil,
			expectedError:    ErrRegisterRewardInternal,
			mockExpectActions: func() {
				mock.ExpectBegin()
				mockRewardExistsRepo.EXPECT().Exists(context.Background(), gomock.Any(), gomock.Any()).Return(false, errors.New("exists method error")).Times(1)
				mock.ExpectRollback()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExpectActions()

			err := service.Register(context.Background(), reward)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
