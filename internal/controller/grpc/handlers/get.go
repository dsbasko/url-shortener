package handlers

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/dsbasko/url-shortener/api/proto"
	"github.com/dsbasko/url-shortener/internal/service/jwt"
)

// GetOriginalURL returns the original URL for the given short URL
func (s *URLShortenerServer) GetOriginalURL(
	ctx context.Context,
	in *pb.GetOriginalURLRequest,
) (*pb.GetOriginalURLResponse, error) {
	var resp pb.GetOriginalURLResponse

	foundURL, err := s.urlService.GetURL(ctx, in.ShortUrl)
	if err != nil {
		s.log.Errorf("failed to get url from service layer: %v", err)
		return &resp, status.Errorf(codes.Internal, "failed to get url from service layer: %v", err)
	}

	resp.OriginalUrl = foundURL.OriginalURL
	return &resp, nil
}

// GetURLsByUserID returns the URLs for the given user
func (s *URLShortenerServer) GetURLsByUserID(
	ctx context.Context,
	_ *pb.GetURLsByUserIDRequest,
) (*pb.GetURLsByUserIDResponse, error) {
	var resp pb.GetURLsByUserIDResponse

	token, err := jwt.GetFromContextGRPC(ctx)
	if err != nil {
		s.log.Errorf("failed to get token from context: %v", err)
		return &resp, err
	}

	userID := jwt.TokenToUserID(token)
	foundURLs, err := s.urlService.GetURLsByUserID(ctx, userID)
	if err != nil {
		s.log.Errorf("failed to create url from service layer: %v", err)
		return &resp, err
	}

	result := make([]*pb.URLEntity, 0, len(foundURLs))
	for _, foundURL := range foundURLs {
		result = append(result, &pb.URLEntity{
			Id:          foundURL.ID,
			UserId:      foundURL.UserID,
			ShortUrl:    foundURL.ShortURL,
			OriginalUrl: foundURL.OriginalURL,
			IsDeleted:   foundURL.DeletedFlag,
		})
	}

	resp.Result = result
	return &resp, nil
}
