package handlers

import (
	"github.com/dsbasko/yandex-go-shortener/internal/service/urls"
	"github.com/dsbasko/yandex-go-shortener/pkg/logger"
)

// Handler a collection of handlers.
type Handler struct {
	log    *logger.Logger
	pinger Pinger
	urls   urls.URLs
}

// New creates a new handler constructor.
func New(log *logger.Logger, pinger Pinger, urlService urls.URLs) Handler {
	return Handler{
		log:    log,
		pinger: pinger,
		urls:   urlService,
	}
}
