package urls

import (
	"context"

	"github.com/dsbasko/yandex-go-shortener/pkg/logger"
)

// URLMutator is an interface for deleting URLs.
type URLMutator interface {
	URLCreator // URLMutator is a service for creating URLs.
	URLDeleter // URLMutator is an interface for deleting URLs.
}

// URLs a service for working with URLs.
type URLs struct {
	log         *logger.Logger
	urlProvider URLProvider
	urlMutator  URLMutator
	deleteTask  chan map[string][]string // deleteTask map [userID][]shortURLs
}

// New creates new URLs service.
func New(
	ctx context.Context,
	log *logger.Logger,
	urlProvider URLProvider,
	urlMutator URLMutator,
) URLs {
	service := URLs{
		log:         log,
		urlProvider: urlProvider,
		urlMutator:  urlMutator,
		deleteTask:  make(chan map[string][]string, 1),
	}

	go service.deleteWorker(ctx)

	return service
}

// Generate mocks for tests.
//go:generate ../../../bin/mockgen -destination=./mocks/url-mutator.go -package=mock_urls github.com/dsbasko/yandex-go-shortener/internal/service/urls URLMutator
