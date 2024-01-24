package middlewares

import (
	"github.com/dsbasko/yandex-go-shortener/pkg/logger"
)

// Middlewares a collection of middlewares.
type Middlewares struct {
	log *logger.Logger
}

// New creates a new middlewares constructor.
func New(log *logger.Logger) *Middlewares {
	return &Middlewares{
		log: log,
	}
}
