package urls

import (
	"context"

	"github.com/dsbasko/yandex-go-shortener/pkg/graceful"
	"github.com/dsbasko/yandex-go-shortener/pkg/logger"
)

// Mutator is an interface for deleting URLs.
type Mutator interface {
	Creator // Creator is a service for creating URLs.
	Deleter // Deleter is an interface for deleting URLs.
}

// URLs a service for working with URLs.
type URLs struct {
	log         *logger.Logger
	urlProvider Provider
	urlMutator  Mutator
	urlAnalyzer Analyzer
	deleteTask  chan map[string][]string // deleteTask map [userID][]shortURLs
}

// New creates new URLs service.
func New(
	ctx context.Context,
	log *logger.Logger,
	urlProvider Provider,
	urlMutator Mutator,
	urlAnalyzer Analyzer,
) URLs {
	service := URLs{
		log:         log,
		urlProvider: urlProvider,
		urlMutator:  urlMutator,
		urlAnalyzer: urlAnalyzer,
		deleteTask:  make(chan map[string][]string, 1),
	}

	graceful.Add()
	go service.deleteWorker(ctx)

	return service
}

// Generate mocks for tests.
//go:generate ../../../bin/mockgen -destination=./mocks/mutator.go -package=mock_urls github.com/dsbasko/yandex-go-shortener/internal/service/urls Mutator
//go:generate ../../../bin/mockgen -destination=./mocks/provider.go -package=mock_urls github.com/dsbasko/yandex-go-shortener/internal/service/urls Provider
//go:generate ../../../bin/mockgen -destination=./mocks/analyzer.go -package=mock_urls github.com/dsbasko/yandex-go-shortener/internal/service/urls Analyzer
