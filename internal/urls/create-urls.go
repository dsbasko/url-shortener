package urls

import (
	"context"
	"fmt"
	goURL "net/url"

	"github.com/dsbasko/yandex-go-shortener/internal/config"
	"github.com/dsbasko/yandex-go-shortener/internal/entities"
	"github.com/dsbasko/yandex-go-shortener/internal/jwt"
	"github.com/dsbasko/yandex-go-shortener/pkg/api"
)

// CreateURLs creates a new short url.
func (u *URLs) CreateURLs(
	ctx context.Context,
	dto []api.CreateURLsRequest,
) ([]api.CreateURLsResponse, error) {
	token, err := jwt.GetFromContext(ctx)
	if err != nil {
		return []api.CreateURLsResponse{}, fmt.Errorf("failed to get token from context: %w", err)
	}

	userID := jwt.TokenToUserID(token)

	response := make([]api.CreateURLsResponse, 0, len(dto))
	urlEntities := make([]entities.URL, 0, len(dto))

	for _, url := range dto {
		if _, err = goURL.ParseRequestURI(url.OriginalURL); err != nil {
			return []api.CreateURLsResponse{}, fmt.Errorf("failed to parse url: %w", ErrInvalidURL)
		}

		shortURL := RandomString(config.GetShortURLLen())

		urlEntities = append(urlEntities, entities.URL{
			OriginalURL: url.OriginalURL,
			ShortURL:    shortURL,
			UserID:      userID,
		})

		response = append(response, api.CreateURLsResponse{
			CorrelationID: url.CorrelationID,
			ShortURL:      fmt.Sprintf("%s%s", config.GetBaseURL(), shortURL),
		})
	}

	if len(urlEntities) > 0 {
		if _, err = u.storage.CreateURLs(ctx, urlEntities); err != nil {
			return []api.CreateURLsResponse{}, fmt.Errorf("failed to create a url in the storage: %w", err)
		}
	}

	return response, nil
}
