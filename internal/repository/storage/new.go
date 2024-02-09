package storage

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/dsbasko/yandex-go-shortener/internal/config"
	"github.com/dsbasko/yandex-go-shortener/internal/interfaces"
	"github.com/dsbasko/yandex-go-shortener/internal/repository/storage/file"
	"github.com/dsbasko/yandex-go-shortener/internal/repository/storage/memory"
	"github.com/dsbasko/yandex-go-shortener/internal/repository/storage/mock"
	"github.com/dsbasko/yandex-go-shortener/internal/repository/storage/psql"
	"github.com/dsbasko/yandex-go-shortener/pkg/logger"
)

// New creates a new instance of the storage.
func New(ctx context.Context, log *logger.Logger) (interfaces.Storage, error) {
	if len(config.GetPsqlDSN()) > 0 {
		return psql.New(ctx, log)
	}

	if len(config.GetStoragePath()) > 0 {
		return file.New(ctx, log)
	}

	return memory.New(ctx, log)
}

// NewMock creates a new instance of the mock storage.
func NewMock(t *testing.T) *mock.MockStorage {
	controller := gomock.NewController(t)
	defer controller.Finish()
	return mock.NewMockStorage(controller)
}
