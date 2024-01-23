package jwt

import (
	"context"
)

func GetFromContext(ctx context.Context) (string, error) {
	if token, ok := ctx.Value(ContextKey).(string); ok {
		return token, nil
	}

	return "", ErrNotFoundFromContext
}
