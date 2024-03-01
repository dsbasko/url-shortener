package handlers

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/dsbasko/url-shortener/api/proto"
)

// Pinger checks connection to storage.
type Pinger interface {
	// Ping checks connection to storage.
	Ping(ctx context.Context) error
}

// URLShortenerServer is the server that provides the URLShortener service.
func (s *URLShortenerServer) Ping(ctx context.Context, _ *pb.PingRequest) (*pb.PingResponse, error) {
	var resp pb.PingResponse

	if err := s.pinger.Ping(ctx); err != nil {
		return nil, status.Error(codes.Unavailable, "")
	}

	resp.Message = "pong"
	return &resp, nil
}
