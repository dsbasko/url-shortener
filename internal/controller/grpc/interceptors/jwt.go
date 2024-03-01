package interceptors

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/dsbasko/url-shortener/internal/service/jwt"
	"github.com/dsbasko/url-shortener/pkg/logger"
)

// JWTInterceptor checks if the jwt token is valid
func JWTInterceptor(log *logger.Logger) grpc.UnaryServerInterceptor {
	log.Debug("jwt interceptor is enabled")

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		token, err := jwt.GetFromContextGRPC(ctx)
		if err != nil {
			return generateJWT(ctx, log, handler, req)
		}

		if valid := jwt.TokenValidate(token); !valid {
			return generateJWT(ctx, log, handler, req)
		}

		newCTX := context.WithValue(ctx, jwt.ContextKey, token)
		resp, err := handler(newCTX, req)

		return resp, err
	}
}

func generateJWT(
	ctx context.Context,
	log *logger.Logger,
	handler grpc.UnaryHandler,
	req any,
) (any, error) {
	token, err := jwt.GenerateToken()
	if err != nil {
		log.Debugf("failed to generate jwt token: %s", err.Error())
		return handler(ctx, req)
	}

	newCTX := context.WithValue(ctx, jwt.ContextKey, token)
	md := metadata.Pairs(jwt.ContextString, token)
	if err = grpc.SendHeader(newCTX, md); err != nil {
		log.Errorf("failed to send token in header: %s", err.Error())
	}

	return handler(newCTX, req)
}
