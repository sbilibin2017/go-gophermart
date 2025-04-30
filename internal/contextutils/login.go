package contextutils

import (
	"context"
	"errors"
)

type contextLoginKey struct{}

func SetLogin(ctx context.Context, login string) context.Context {
	return context.WithValue(ctx, contextLoginKey{}, login)
}

func GetLogin(ctx context.Context) (string, error) {
	login, ok := ctx.Value(contextLoginKey{}).(string)
	if !ok {
		return "", errors.New("login is not in context")
	}
	return login, nil
}
