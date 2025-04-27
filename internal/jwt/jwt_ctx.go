package jwt

import (
	"context"
	"fmt"
)

type contextJWTPayloadKey string

const contextJWTPayload contextJWTPayloadKey = "jwt"

func GetJWTPayload(ctx context.Context) (map[string]any, error) {
	payload, ok := ctx.Value(contextJWTPayload).(map[string]any)
	if !ok {
		return nil, fmt.Errorf("payload not found or has invalid type")
	}
	return payload, nil
}

func SetJWT(ctx context.Context, payload map[string]any) context.Context {
	return context.WithValue(ctx, contextJWTPayload, payload)
}
