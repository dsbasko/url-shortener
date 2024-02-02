package memory

import (
	"context"
	"sync"

	"github.com/dsbasko/yandex-go-shortener/internal/entities"
	"github.com/dsbasko/yandex-go-shortener/internal/interfaces"
	"github.com/dsbasko/yandex-go-shortener/pkg/logger"
)

// Storage a memory storage.
type Storage struct {
	mu    *sync.RWMutex
	log   *logger.Logger
	store map[string]entities.URL
}

// Ensure that Storage implements the Storage interface.
var _ interfaces.Storage = (*Storage)(nil)

// New creates a new instance of the memory storage.
func New(_ context.Context, log *logger.Logger) (interfaces.Storage, error) {
	log.Infof("memory storage initialized")

	return &Storage{
		mu:    &sync.RWMutex{},
		log:   log,
		store: map[string]entities.URL{},
	}, nil
}
