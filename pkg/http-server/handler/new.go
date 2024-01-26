package handler

import (
	"github.com/dsbasko/yandex-go-shortener/internal/interfaces"
	"github.com/dsbasko/yandex-go-shortener/internal/urls"
	"github.com/dsbasko/yandex-go-shortener/pkg/logger"
)

// Handler a collection of handlers.
type Handler struct {
	log     *logger.Logger
	storage interfaces.Storage
	urls    urls.URLs
}

// New creates a new handler constructor.
func New(log *logger.Logger, storage interfaces.Storage, urlService urls.URLs) *Handler {
	return &Handler{
		log:     log,
		storage: storage,
		urls:    urlService,
	}
}
