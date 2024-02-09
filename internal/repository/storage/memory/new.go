package memory

import (
	"context"
	"sync"

	"github.com/dsbasko/yandex-go-shortener/internal/entity"
	"github.com/dsbasko/yandex-go-shortener/pkg/logger"
)

// Storage a memory storage.
type Storage struct {
	mu    *sync.RWMutex
	log   *logger.Logger
	store map[string]entity.URL
}

// New creates a new instance of the memory storage.
func New(_ context.Context, log *logger.Logger) (*Storage, error) {
	log.Infof("memory storage initialized")

	return &Storage{
		mu:    &sync.RWMutex{},
		log:   log,
		store: map[string]entity.URL{},
	}, nil
}
