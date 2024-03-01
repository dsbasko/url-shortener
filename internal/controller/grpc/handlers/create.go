package handlers

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api "github.com/dsbasko/url-shortener/api/http"
	pb "github.com/dsbasko/url-shortener/api/proto"
)

// CreateURL creates a short URL from the given URL
func (s *URLShortenerServer) CreateURL(
	ctx context.Context,
	in *pb.CreateURLRequest,
) (*pb.CreateURLResponse, error) {
	var resp pb.CreateURLResponse

	createdURL, _, err := s.urlService.CreateURL(ctx, in.Url)
	if err != nil {
		s.log.Errorf("failed to create url from service layer: %v", err)
		return &resp, status.Errorf(codes.Internal, "%v", err)
	}

	resp.Result = createdURL.ShortURL
	return &resp, nil
}

// CreateURLs creates short URLs from the given URLs
func (s *URLShortenerServer) CreateURLs(
	ctx context.Context,
	in *pb.CreateURLsRequest,
) (*pb.CreateURLsResponse, error) {
	var resp pb.CreateURLsResponse

	dto := make([]api.CreateURLsRequest, len(in.Dto))
	for i, d := range in.Dto {
		dto[i] = api.CreateURLsRequest{
			CorrelationID: d.CorrelationId,
			OriginalURL:   d.OriginalUrl,
		}
	}

	createdURLs, err := s.urlService.CreateURLs(ctx, dto)
	if err != nil {
		s.log.Errorf("failed to create url from service layer: %v", err)
		return &resp, status.Errorf(codes.Internal, "%v", err)
	}

	result := make([]*pb.CreateURLsResponse_Data, 0, len(createdURLs))
	for _, createdURL := range createdURLs {
		result = append(result, &pb.CreateURLsResponse_Data{
			CorrelationId: createdURL.CorrelationID,
			ShortUrl:      createdURL.ShortURL,
		})
	}

	resp.Result = result
	return &resp, nil
}
