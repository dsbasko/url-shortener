package interfaces

import (
	"context"

	"github.com/dsbasko/yandex-go-shortener/internal/entities"
)

// Storage is an interface for storage.
type Storage interface {
	// Ping checks connection to storage.
	Ping(ctx context.Context) error

	// Close closes connection to storage.
	Close() error

	// GetURLByOriginalURL gets URL by original URL.
	GetURLByOriginalURL(ctx context.Context, originalURL string) (resp entities.URL, err error)

	// GetURLByShortURL gets URL by short URL.
	GetURLByShortURL(ctx context.Context, shortURL string) (resp entities.URL, err error)

	// GetURLsByUserID gets URLs by user ID.
	GetURLsByUserID(ctx context.Context, userID string) (resp []entities.URL, err error)

	// CreateURL creates a new URL.
	CreateURL(ctx context.Context, dto entities.URL) (resp entities.URL, unique bool, err error)

	// CreateURLs creates URLs.
	CreateURLs(ctx context.Context, dto []entities.URL) (resp []entities.URL, err error)

	// DeleteURLs deletes URLs.
	DeleteURLs(ctx context.Context, dto []entities.URL) (resp []entities.URL, err error)
}

// Generate mocks for tests.
//go:generate ../../bin/mockgen -destination=../storage/mock/mock.go -package=mock github.com/dsbasko/yandex-go-shortener/internal/interfaces Storage
