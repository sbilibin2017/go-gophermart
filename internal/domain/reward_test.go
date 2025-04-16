package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNeeReward_Success(t *testing.T) {
	r := NewReward("ProductA", 20, string(Percentage))
	require.NotNil(t, r)

}

func TestApplyReward_Percentage(t *testing.T) {
	r := &Reward{
		Match:      "ProductA",
		Reward:     20,
		RewardType: Percentage,
	}
	price := uint64(1000)
	expectedReward := uint64(200)
	actualReward := r.ApplyReward(price)
	require.NotNil(t, actualReward)
	assert.Equal(t, expectedReward, *actualReward)
}

func TestApplyReward_Points(t *testing.T) {
	r := &Reward{
		Match:      "ProductB",
		Reward:     150,
		RewardType: Points,
	}
	price := uint64(1000)
	expectedReward := uint64(150)
	actualReward := r.ApplyReward(price)
	require.NotNil(t, actualReward)
	assert.Equal(t, expectedReward, *actualReward)
}

func TestApplyReward_InvalidRewardType(t *testing.T) {
	r := &Reward{
		Match:      "ProductC",
		Reward:     100,
		RewardType: "invalid",
	}
	price := uint64(1000)
	actualReward := r.ApplyReward(price)
	assert.Nil(t, actualReward)
}

func TestApplyReward_ZeroPrice(t *testing.T) {
	r := &Reward{
		Match:      "ProductD",
		Reward:     20,
		RewardType: Percentage,
	}
	price := uint64(0)
	expectedReward := uint64(0)
	actualReward := r.ApplyReward(price)
	require.NotNil(t, actualReward)
	assert.Equal(t, expectedReward, *actualReward)
}
