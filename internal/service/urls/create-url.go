package urls

import (
	"context"
	"fmt"
	"net/url"

	"github.com/dsbasko/yandex-go-shortener/internal/config"
	"github.com/dsbasko/yandex-go-shortener/internal/entities"
	"github.com/dsbasko/yandex-go-shortener/internal/service/jwt"
)

// CreateURL creates a new short url.
func (u *URLs) CreateURL(
	ctx context.Context,
	originalURL string,
) (resp entities.URL, unique bool, err error) {
	var dto entities.URL

	if _, err = url.ParseRequestURI(originalURL); err != nil {
		return entities.URL{}, false, fmt.Errorf("failed to parse url: %w", ErrInvalidURL)
	}

	token, err := jwt.GetFromContext(ctx)
	if err != nil {
		return entities.URL{}, false, fmt.Errorf("failed to get token from context: %w", err)
	}

	userID := jwt.TokenToUserID(token)

	dto.OriginalURL = originalURL
	dto.UserID = userID
	dto.ShortURL = RandomString(config.GetShortURLLen())

	resp, uniq, err := u.storage.CreateURL(ctx, dto)
	if err != nil {
		return entities.URL{}, false, fmt.Errorf("failed to create a url in the storage: %w", err)
	}

	resp.ShortURL = fmt.Sprintf("%s%s", config.GetBaseURL(), resp.ShortURL)
	return resp, uniq, nil
}
