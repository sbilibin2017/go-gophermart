package contextutils

import (
	"context"
	"testing"

	"github.com/sbilibin2017/go-gophermart/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestSetClaims(t *testing.T) {
	claims := &types.Claims{Login: "testuser"}
	ctx := context.Background()
	ctx = SetClaims(ctx, claims)
	retrievedClaims, err := GetClaims(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedClaims)
	assert.Equal(t, claims.Login, retrievedClaims.Login)
}

func TestGetClaims_Error(t *testing.T) {
	ctx := context.Background()
	retrievedClaims, err := GetClaims(ctx)
	assert.Error(t, err)
	assert.Nil(t, retrievedClaims)
}

func TestSetClaims_Overwrite(t *testing.T) {
	claims1 := &types.Claims{Login: "testuser1"}
	claims2 := &types.Claims{Login: "testuser2"}
	ctx := context.Background()
	ctx = SetClaims(ctx, claims1)
	ctx = SetClaims(ctx, claims2)
	retrievedClaims, err := GetClaims(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedClaims)
	assert.Equal(t, claims2.Login, retrievedClaims.Login)
}

func TestSetClaims_NilClaims(t *testing.T) {
	ctx := context.Background()
	ctx = SetClaims(ctx, nil)
	retrievedClaims, _ := GetClaims(ctx)
	assert.Nil(t, retrievedClaims)
}
