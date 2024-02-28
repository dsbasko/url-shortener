package urls

import (
	"context"
	"fmt"

	"github.com/dsbasko/yandex-go-shortener/internal/config"
	"github.com/dsbasko/yandex-go-shortener/internal/entity"
)

// Provider is an interface for providing URLs.
type Provider interface {
	// GetURLByOriginalURL gets URL by original URL.
	GetURLByOriginalURL(ctx context.Context, originalURL string) (resp entity.URL, err error)

	// GetURLByShortURL gets URL by short URL.
	GetURLByShortURL(ctx context.Context, shortURL string) (resp entity.URL, err error)

	// GetURLsByUserID gets URLs by user ID.
	GetURLsByUserID(ctx context.Context, userID string) (resp []entity.URL, err error)
}

// GetURL returns URL by short URL.
func (u *URLs) GetURL(ctx context.Context, shortURL string) (entity.URL, error) {
	if shortURL == "" {
		return entity.URL{}, ErrInvalidURL
	}

	storageResp, err := u.urlProvider.GetURLByShortURL(ctx, shortURL)
	if err != nil {
		return entity.URL{}, fmt.Errorf("error getting url from storage: %w", err)
	}

	return storageResp, nil
}

// GetURLsByUserID returns all URLs by user ID.
func (u *URLs) GetURLsByUserID(ctx context.Context, userID string) ([]entity.URL, error) {
	urls, err := u.urlProvider.GetURLsByUserID(ctx, userID)
	if err != nil {
		return []entity.URL{}, fmt.Errorf("error getting url from storage: %w", err)
	}

	resp := make([]entity.URL, 0, len(urls))
	for _, url := range urls {
		url.ShortURL = fmt.Sprintf("%s%s", config.BaseURL(), url.ShortURL)
		resp = append(resp, url)
	}

	return resp, nil
}
