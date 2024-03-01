package jwt

import (
	"context"

	"google.golang.org/grpc/metadata"
)

// GetFromContext gets jwt token from context.
func GetFromContext(ctx context.Context) (string, error) {
	if token, ok := ctx.Value(ContextKey).(string); ok {
		return token, nil
	}

	return "", ErrNotFoundFromContext
}

// GetFromContext gets jwt token from context.
func GetFromContextGRPC(ctx context.Context) (string, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get(ContextString)
		if len(values) > 0 {
			return values[0], nil
		}
	}

	return "", ErrNotFoundFromContext
}
