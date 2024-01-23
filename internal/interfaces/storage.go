package interfaces

import (
	"context"

	"github.com/dsbasko/yandex-go-shortener/internal/entities"
)

type Storage interface {
	Ping(ctx context.Context) error
	Close() error

	GetURLByOriginalURL(ctx context.Context, originalURL string) (resp entities.URL, err error)
	GetURLByShortURL(ctx context.Context, shortURL string) (resp entities.URL, err error)
	GetURLsByUserID(ctx context.Context, userID string) (resp []entities.URL, err error)

	CreateURL(ctx context.Context, dto entities.URL) (resp entities.URL, unique bool, err error)
	CreateURLs(ctx context.Context, dto []entities.URL) (resp []entities.URL, err error)

	DeleteURLs(ctx context.Context, dto []entities.URL) (resp []entities.URL, err error)
}

// Generate mocks for tests.
//go:generate ../../bin/mockgen -destination=../storage/mock/mock.go -package=mock github.com/dsbasko/yandex-go-shortener/internal/interfaces Storage
