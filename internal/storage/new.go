package storage

import (
	"context"
	"testing"

	"github.com/dsbasko/yandex-go-shortener/internal/config"
	"github.com/dsbasko/yandex-go-shortener/internal/interfaces"
	"github.com/dsbasko/yandex-go-shortener/internal/storage/file"
	"github.com/dsbasko/yandex-go-shortener/internal/storage/memory"
	"github.com/dsbasko/yandex-go-shortener/internal/storage/mock"
	"github.com/dsbasko/yandex-go-shortener/internal/storage/psql"
	"github.com/dsbasko/yandex-go-shortener/pkg/logger"
	"github.com/golang/mock/gomock"
)

func New(ctx context.Context, log *logger.Logger) (interfaces.Storage, error) {
	if len(config.GetPsqlDSN()) > 0 {
		return psql.New(ctx, log)
	}

	if len(config.GetStoragePath()) > 0 {
		return file.New(ctx, log)
	}

	return memory.New(ctx, log)
}

func NewMock(t *testing.T) *mock.MockStorage {
	controller := gomock.NewController(t)
	defer controller.Finish()
	return mock.NewMockStorage(controller)
}
