package urls

import (
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
}

// New creates new URLs service.
func New(log *logger.Logger, urlProvider URLProvider, urlMutator URLMutator) URLs {
	return URLs{
		log:         log,
		urlProvider: urlProvider,
		urlMutator:  urlMutator,
	}
}

// Generate mocks for tests.
//go:generate ../../../bin/mockgen -destination=./mocks/url-mutator.go -package=mock_urls github.com/dsbasko/yandex-go-shortener/internal/service/urls URLMutator
