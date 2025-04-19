package repositories

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRewardExistsRepository_Exists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuerier := NewMockRewardExistsQuerier(ctrl)
	repo := NewRewardExistsRepository(mockQuerier)

	ctx := context.Background()
	match := map[string]any{"key": "value"}

	tests := []struct {
		name       string
		mockReturn func()
		expected   bool
		err        error
	}{
		{
			name: "should return true when reward exists",
			mockReturn: func() {
				mockQuerier.EXPECT().
					Query(ctx, gomock.Any(), rewardExistsQuery, match).
					SetArg(1, true).
					Return(nil)
			},
			expected: true,
			err:      nil,
		},
		{
			name: "should return false when reward does not exist",
			mockReturn: func() {
				mockQuerier.EXPECT().
					Query(ctx, gomock.Any(), rewardExistsQuery, match).
					SetArg(1, false).
					Return(nil)
			},
			expected: false,
			err:      nil,
		},
		{
			name: "should return error when query fails",
			mockReturn: func() {
				mockQuerier.EXPECT().
					Query(ctx, gomock.Any(), rewardExistsQuery, match).
					Return(errors.New("query failed"))
			},
			expected: false,
			err:      errors.New("query failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockReturn()
			exists, err := repo.Exists(ctx, match)

			if tt.err != nil {
				require.Error(t, err)
				assert.EqualError(t, err, tt.err.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, exists)
			}
		})
	}
}
