package interceptors

import (
	"context"
	"time"

	"google.golang.org/grpc"

	"github.com/dsbasko/url-shortener/pkg/logger"
)

// LoggerInterceptor logs the request and response of the grpc server
func LoggerInterceptor(log *logger.Logger) grpc.UnaryServerInterceptor {
	log.Debug("logger interceptor enabled")

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		args := make([]any, 0, 3)
		timeStart := time.Now()

		// Вызов обработчика
		resp, err := handler(ctx, req)

		args = append(args, []any{
			"response_duration", time.Since(timeStart).String(),
		}...)

		log.Infow("request is done", args...)
		return resp, err
	}
}
