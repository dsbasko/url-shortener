package urls

import (
	"github.com/dsbasko/yandex-go-shortener/internal/interfaces"
	"github.com/dsbasko/yandex-go-shortener/pkg/logger"
)

type URLs struct {
	log     *logger.Logger
	storage interfaces.Storage
}

func New(log *logger.Logger, storage interfaces.Storage) URLs {
	return URLs{
		log:     log,
		storage: storage,
	}
}
