package urls

import (
	"context"
	"fmt"
	goURL "net/url"

	"github.com/dsbasko/url-shortener/internal/config"
	"github.com/dsbasko/url-shortener/internal/entity"
	"github.com/dsbasko/url-shortener/internal/service/jwt"
	"github.com/dsbasko/url-shortener/pkg/api"
)

// Creator is a service for creating URLs.
type Creator interface {
	// CreateURL creates a new URL.
	CreateURL(ctx context.Context, dto entity.URL) (resp entity.URL, unique bool, err error)

	// CreateURLs creates URLs.
	CreateURLs(ctx context.Context, dto []entity.URL) (resp []entity.URL, err error)
}

// CreateURL creates a new short url.
func (u *URLs) CreateURL(
	ctx context.Context,
	originalURL string,
) (resp entity.URL, unique bool, err error) {
	var dto entity.URL
	var urlCreator Creator = u.urlMutator

	if _, err = goURL.ParseRequestURI(originalURL); err != nil {
		return entity.URL{}, false, fmt.Errorf("failed to parse url: %w", ErrInvalidURL)
	}

	token, err := jwt.GetFromContext(ctx)
	if err != nil {
		return entity.URL{}, false, fmt.Errorf("failed to get token from context: %w", err)
	}

	userID := jwt.TokenToUserID(token)

	dto.OriginalURL = originalURL
	dto.UserID = userID
	dto.ShortURL = RandomString(config.ShortURLLen())

	resp, uniq, err := urlCreator.CreateURL(ctx, dto)
	if err != nil {
		return entity.URL{}, false, fmt.Errorf("failed to create a url in the storage: %w", err)
	}

	resp.ShortURL = fmt.Sprintf("%s%s", config.BaseURL(), resp.ShortURL)
	return resp, uniq, nil
}

// CreateURLs creates a new short url.
func (u *URLs) CreateURLs(
	ctx context.Context,
	dto []api.CreateURLsRequest,
) ([]api.CreateURLsResponse, error) {
	var urlCreator Creator = u.urlMutator

	token, err := jwt.GetFromContext(ctx)
	if err != nil {
		return []api.CreateURLsResponse{}, fmt.Errorf("failed to get token from context: %w", err)
	}

	userID := jwt.TokenToUserID(token)

	response := make([]api.CreateURLsResponse, 0, len(dto))
	urlEntities := make([]entity.URL, 0, len(dto))

	for _, url := range dto {
		if _, err = goURL.ParseRequestURI(url.OriginalURL); err != nil {
			return []api.CreateURLsResponse{}, fmt.Errorf("failed to parse url: %w", ErrInvalidURL)
		}

		shortURL := RandomString(config.ShortURLLen())

		urlEntities = append(urlEntities, entity.URL{
			OriginalURL: url.OriginalURL,
			ShortURL:    shortURL,
			UserID:      userID,
		})

		response = append(response, api.CreateURLsResponse{
			CorrelationID: url.CorrelationID,
			ShortURL:      fmt.Sprintf("%s%s", config.BaseURL(), shortURL),
		})
	}

	if len(urlEntities) > 0 {
		if _, err = urlCreator.CreateURLs(ctx, urlEntities); err != nil {
			return []api.CreateURLsResponse{}, fmt.Errorf("failed to create a url in the storage: %w", err)
		}
	}

	return response, nil
}
