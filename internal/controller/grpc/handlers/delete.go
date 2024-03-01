package handlers

import (
	"context"

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
		return nil, err
	}

	userID := jwt.TokenToUserID(token)
	err = s.urlService.DeleteURLs(userID, in.Urls)
	if err != nil {
		s.log.Errorf("failed to delete urls: %v", err)
		return nil, err
	}

	return nil, nil
}
