package repositories

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRewardSaveRepository_Save(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockExecutor := NewMockRewardSaveExecutor(ctrl)
	repo := NewRewardSaveRepository(mockExecutor)

	ctx := context.Background()
	args := map[string]any{"match": "match_value", "reward": 100, "reward_type": "cash"}

	tests := []struct {
		name       string
		mockReturn func()
		err        error
	}{
		{
			name: "should save reward successfully",
			mockReturn: func() {
				mockExecutor.EXPECT().
					Execute(ctx, rewardSaveQuery, args).
					Return(nil)
			},
			err: nil,
		},
		{
			name: "should return error when execution fails",
			mockReturn: func() {
				mockExecutor.EXPECT().
					Execute(ctx, rewardSaveQuery, args).
					Return(errors.New("execution failed"))
			},
			err: errors.New("execution failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockReturn()
			err := repo.Save(ctx, args)

			if tt.err != nil {
				require.Error(t, err)
				assert.EqualError(t, err, tt.err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
