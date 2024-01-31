package urls

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dsbasko/yandex-go-shortener/internal/storage"
	"github.com/dsbasko/yandex-go-shortener/pkg/logger"
)

func TestNew(t *testing.T) {
	log := logger.NewMock()
	store := storage.NewMock(t)
	t.Run("Success", func(t *testing.T) {
		service := New(log, store)
		assert.NotNil(t, service)
		mockService := URLs{log: log, storage: store}
		assert.Equal(t, mockService, service)
	})
}
