package handlers

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/dsbasko/url-shortener/api/proto"
	"github.com/dsbasko/url-shortener/internal/service/jwt"
)

// DeleteURLs deletes the given URLs
func (s *URLShortenerServer) DeleteURLs(
	ctx context.Context,
	in *pb.DeleteURLsRequest,
) (*pb.DeleteURLsResponse, error) {
	token, err := jwt.GetFromContextGRPC(ctx)
	if err != nil {
		s.log.Errorf("failed to get token from context: %v", err)
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	userID := jwt.TokenToUserID(token)
	s.urlService.DeleteURLs(userID, in.Urls)
	return nil, nil
}
