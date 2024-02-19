package jwt

import (
	"context"
)

// GetFromContext gets jwt token from context.
func GetFromContext(ctx context.Context) (string, error) {
	if token, ok := ctx.Value(ContextKey).(string); ok {
		return token, nil
	}

	return "", ErrNotFoundFromContext
}
