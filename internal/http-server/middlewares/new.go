package middlewares

import (
	"github.com/dsbasko/yandex-go-shortener/pkg/logger"
)

type Middlewares struct {
	log *logger.Logger
}

func New(log *logger.Logger) *Middlewares {
	return &Middlewares{
		log: log,
	}
}
