package urls

import (
	"context"
	"fmt"

	"github.com/dsbasko/yandex-go-shortener/internal/entities"
)

func (u *URLs) GetURL(ctx context.Context, shortURL string) (entities.URL, error) {
	if shortURL == "" {
		return entities.URL{}, ErrInvalidURL
	}

	storeResp, err := u.storage.GetURLByShortURL(ctx, shortURL)
	if err != nil {
		return entities.URL{}, fmt.Errorf("error getting url from storage: %w", err)
	}

	return storeResp, nil
}
