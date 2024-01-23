package handler

import (
	"testing"

	"github.com/dsbasko/yandex-go-shortener/internal/storage"
	"github.com/dsbasko/yandex-go-shortener/internal/urls"
	"github.com/dsbasko/yandex-go-shortener/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	log := logger.NewMock()
	store := storage.NewMock(t)
	urlService := urls.New(log, store)

	t.Run("Success", func(t *testing.T) {
		handler := New(log, store, urlService)
		assert.NotNil(t, handler)
		mockServiceHandler := &Handler{log: log, storage: store, urls: urlService}
		assert.Equal(t, mockServiceHandler, handler)
	})
}
