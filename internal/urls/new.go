package urls

import (
	"github.com/dsbasko/yandex-go-shortener/internal/interfaces"
	"github.com/dsbasko/yandex-go-shortener/pkg/logger"
)

// URLs a service for working with URLs.
type URLs struct {
	log     *logger.Logger
	storage interfaces.Storage
}

// New creates new URLs service.
func New(log *logger.Logger, storage interfaces.Storage) URLs {
	return URLs{
		log:     log,
		storage: storage,
	}
}
