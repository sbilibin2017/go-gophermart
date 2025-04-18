package queries

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildRewardExistsFilters(t *testing.T) {
	tests := []struct {
		name         string
		filter       map[string]any
		expectedSQL  string
		expectedArgs []any
	}{
		{
			name:         "Single Filter",
			filter:       map[string]any{"reward_id": 123},
			expectedSQL:  `"reward_id" = $1`,
			expectedArgs: []any{123},
		},
		{
			name:         "Multiple Filters",
			filter:       map[string]any{"reward_id": 123, "user_id": 456},
			expectedSQL:  `"reward_id" = $1 AND "user_id" = $2`,
			expectedArgs: []any{123, 456},
		},
		{
			name:         "Empty Filter",
			filter:       map[string]any{},
			expectedSQL:  "",
			expectedArgs: nil, // Change to nil to compare against nil
		},
		{
			name:         "Special Character in Filter",
			filter:       map[string]any{"reward_name": "Summer Special!"},
			expectedSQL:  `"reward_name" = $1`,
			expectedArgs: []any{"Summer Special!"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, args := BuildRewardExistsFilters(tt.filter)
			assert.Equal(t, tt.expectedSQL, sql)
			assert.Equal(t, tt.expectedArgs, args) // This will now pass with nil
		})
	}
}
