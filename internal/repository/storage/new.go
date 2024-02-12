package storage

import (
	"context"

	"github.com/dsbasko/yandex-go-shortener/internal/config"
	"github.com/dsbasko/yandex-go-shortener/internal/entity"
	"github.com/dsbasko/yandex-go-shortener/internal/repository/storage/file"
	"github.com/dsbasko/yandex-go-shortener/internal/repository/storage/memory"
	"github.com/dsbasko/yandex-go-shortener/internal/repository/storage/psql"
	"github.com/dsbasko/yandex-go-shortener/pkg/logger"
)

// Storage is an interface for storage.
type Storage interface {
	// Ping checks connection to storage.
	Ping(ctx context.Context) error

	// Close closes connection to storage.
	Close() error

	// GetURLByOriginalURL gets URL by original URL.
	GetURLByOriginalURL(ctx context.Context, originalURL string) (resp entity.URL, err error)

	// GetURLByShortURL gets URL by short URL.
	GetURLByShortURL(ctx context.Context, shortURL string) (resp entity.URL, err error)

	// GetURLsByUserID gets URLs by user ID.
	GetURLsByUserID(ctx context.Context, userID string) (resp []entity.URL, err error)

	// CreateURL creates a new URL.
	CreateURL(ctx context.Context, dto entity.URL) (resp entity.URL, unique bool, err error)

	// CreateURLs creates URLs.
	CreateURLs(ctx context.Context, dto []entity.URL) (resp []entity.URL, err error)

	// DeleteURLs deletes URLs.
	DeleteURLs(ctx context.Context, dto []entity.URL) (resp []entity.URL, err error)
}

// New creates a new instance of the storage.
func New(ctx context.Context, log *logger.Logger) (Storage, error) {
	if len(config.PsqlDSN()) > 0 {
		return psql.New(ctx, log)
	}

	if len(config.StoragePath()) > 0 {
		return file.New(ctx, log)
	}

	return memory.New(ctx, log)
}

// Generate mocks for tests.
//go:generate ../../../bin/mockgen -destination=./mocks/storage.go -package=mock_storage github.com/dsbasko/yandex-go-shortener/internal/repository/storage Storage
