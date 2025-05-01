package contextutils

import (
	"context"
	"errors"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type contextLoginKey struct{}

func SetClaims(ctx context.Context, claims *types.Claims) context.Context {
	return context.WithValue(ctx, contextLoginKey{}, claims)
}

func GetClaims(ctx context.Context) (*types.Claims, error) {
	claims, _ := ctx.Value(contextLoginKey{}).(*types.Claims)
	if claims == nil {
		return nil, errors.New("claims are not in context")
	}
	return claims, nil
}
