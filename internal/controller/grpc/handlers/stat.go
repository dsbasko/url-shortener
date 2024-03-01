package handlers

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	pb "github.com/dsbasko/url-shortener/api/proto"
	"github.com/dsbasko/url-shortener/internal/config"
)

// Stat returns the stats of the service
func (s *URLShortenerServer) Stat(
	ctx context.Context,
	_ *pb.StatRequest,
) (*pb.StatResponse, error) {
	var resp pb.StatResponse

	p, ok := peer.FromContext(ctx)
	if !ok {
		s.log.Error("failed to get peer from context")
		return nil, status.Error(codes.Internal, "failed to get peer from context")
	}

	isTrustedSubnet, err := config.IsTrustedSubnet(p.Addr.String())
	if err != nil {
		s.log.Errorf("failed to check trusted subnet: %v", err)
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	if !isTrustedSubnet {
		s.log.Warnf("untrusted subnet: %s", p.Addr.String())
		return nil, status.Errorf(codes.PermissionDenied, "untrusted subnet: %s", p.Addr.String())
	}

	result, err := s.urlService.Stats(ctx)
	if err != nil {
		s.log.Errorf("failed to get stats: %v", err)
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	resp.Urls = result.URLs
	resp.Users = result.Users

	return &resp, nil
}
