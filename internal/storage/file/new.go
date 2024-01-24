package file

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/dsbasko/yandex-go-shortener/internal/config"
	"github.com/dsbasko/yandex-go-shortener/internal/interfaces"
	"github.com/dsbasko/yandex-go-shortener/pkg/logger"
)

// Storage a file storage.
type Storage struct {
	mu     sync.RWMutex
	log    *logger.Logger
	file   *os.File
	writer *bufio.Writer
}

// Ensure that Storage implements the Storage interface.
var _ interfaces.Storage = (*Storage)(nil)

// New creates a new instance of the file storage.
func New(_ context.Context, log *logger.Logger) (interfaces.Storage, error) {
	file, err := os.OpenFile(config.GetStoragePath(), os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_SYNC, 0666) //nolint:gomnd
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	log.Infof("file storage initialized")

	return &Storage{
		log:    log,
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}
