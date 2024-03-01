package grpc

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/dsbasko/url-shortener/api/proto"
	"github.com/dsbasko/url-shortener/internal/config"
	"github.com/dsbasko/url-shortener/internal/controller/grpc/handlers"
	"github.com/dsbasko/url-shortener/internal/controller/grpc/interceptors"
	"github.com/dsbasko/url-shortener/internal/service/urls"
	"github.com/dsbasko/url-shortener/pkg/graceful"
	"github.com/dsbasko/url-shortener/pkg/logger"
)

// Run starts the grpc server
func Run(ctx context.Context, log *logger.Logger, pinger handlers.Pinger, urlService urls.URLs) {
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{
		Port: 8081,
	})
	if err != nil {
		log.Errorf("failed to listen: %v", err)
		return
	}

	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(interceptors.LoggerInterceptor(log)),
		grpc.ChainUnaryInterceptor(interceptors.JWTInterceptor(log)),
	)
	h := handlers.New(log, pinger, urlService)
	pb.RegisterURLShortenerV1Server(srv, &h)

	if config.Env() != "prod" {
		reflection.Register(srv)
	}

	graceful.Add()
	go gracefulShutdown(ctx, log, srv)

	graceful.Add()
	go runServer(ctx, log, listen, srv)
}

func runServer(ctx context.Context, log *logger.Logger, listen *net.TCPListener, srv *grpc.Server) {
	defer graceful.Done()

	_, cancel := context.WithCancel(ctx)
	defer cancel()

	log.Infof("start grpc server on %s", listen.Addr().String())
	if err := srv.Serve(listen); err != nil {
		log.Errorf("failed to serve grpc: %v", err)
		cancel()
	}
}

func gracefulShutdown(ctx context.Context, log *logger.Logger, srv *grpc.Server) {
	defer graceful.Done()

	<-ctx.Done()
	log.Infof("shutdown grpc server by signal")

	srv.GracefulStop()
}
