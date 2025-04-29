package contextutils

import "context"

type contextLoginKey string

const loginContextKey contextLoginKey = "login"

func SetLogin(ctx context.Context, login string) context.Context {
	return context.WithValue(ctx, loginContextKey, login)
}

func GetLogin(ctx context.Context) (string, bool) {
	login, ok := ctx.Value(loginContextKey).(string)
	return login, ok
}
