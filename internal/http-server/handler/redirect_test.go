package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/dsbasko/yandex-go-shortener/internal/config"
	"github.com/dsbasko/yandex-go-shortener/internal/entities"
	"github.com/dsbasko/yandex-go-shortener/internal/storage"
	"github.com/dsbasko/yandex-go-shortener/internal/urls"
	"github.com/dsbasko/yandex-go-shortener/pkg/logger"
	"github.com/dsbasko/yandex-go-shortener/pkg/test"
)

func TestHandler_Redirect(t *testing.T) {
	config.Init() //nolint:errcheck
	log := logger.NewMock()
	store := storage.NewMock(t)
	urlsService := urls.New(log, store)
	router := chi.NewRouter()
	h := New(log, store, urlsService)
	router.Get("/{short_url}", h.Redirect)
	ts := httptest.NewServer(router)
	defer ts.Close()

	tests := []struct {
		name           string
		shortURL       string
		storeCfg       func()
		wantStatusCode int
	}{
		{
			name:     "Not Found",
			shortURL: "42",
			storeCfg: func() {
				store.EXPECT().
					GetURLByShortURL(gomock.Any(), gomock.Any()).
					Return(entities.URL{}, errors.New("not found"))
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:     "Found",
			shortURL: "42",
			storeCfg: func() {
				store.EXPECT().
					GetURLByShortURL(gomock.Any(), gomock.Any()).
					Return(entities.URL{
						ID:          "42",
						ShortURL:    "42",
						OriginalURL: "https://ya.ru/",
					}, nil)
			},
			wantStatusCode: http.StatusTemporaryRedirect,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.storeCfg()
			resp, _ := test.Request(t, ts, &test.RequestArgs{
				Method: "GET",
				Path:   fmt.Sprintf("/%s", tt.shortURL),
			})
			defer resp.Body.Close()

			assert.Equal(t, tt.wantStatusCode, resp.StatusCode)
		})
	}
}

func BenchmarkHandler_Redirect(b *testing.B) {
	config.Init() //nolint:errcheck
	log := logger.NewMock()
	store := storage.NewMock(&testing.T{})
	urlsService := urls.New(log, store)
	router := chi.NewRouter()
	h := New(log, store, urlsService)
	router.Get("/{short_url}", h.Redirect)
	ts := httptest.NewServer(router)
	defer ts.Close()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		store.EXPECT().GetURLByShortURL(gomock.Any(), gomock.Any()).Return(entities.URL{
			ID:          "42",
			ShortURL:    "42",
			OriginalURL: "https://ya.ru/",
			UserID:      "42",
		}, nil)
		b.StartTimer()

		resp, _ := test.Request(&testing.T{}, ts, &test.RequestArgs{
			Method: "GET",
			Path:   "/42",
		})
		resp.Body.Close()
	}
}
