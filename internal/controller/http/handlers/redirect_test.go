package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/dsbasko/url-shortener/internal/config"
	"github.com/dsbasko/url-shortener/internal/entity"
	mockStorage "github.com/dsbasko/url-shortener/internal/repository/storage/mocks"
	"github.com/dsbasko/url-shortener/internal/service/urls"
	"github.com/dsbasko/url-shortener/pkg/logger"
	"github.com/dsbasko/url-shortener/pkg/test"
)

func (s *SuiteHandlers) Test_Redirect() {
	t := s.T()

	tests := []struct {
		name           string
		shortURL       string
		storageCfg     func()
		wantStatusCode int
	}{
		{
			name:     "Not Found",
			shortURL: "42",
			storageCfg: func() {
				s.attr.storage.EXPECT().
					GetURLByShortURL(gomock.Any(), gomock.Any()).
					Return(entity.URL{}, errors.New("not found"))
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:     "Is Deleted",
			shortURL: "42",
			storageCfg: func() {
				s.attr.storage.EXPECT().
					GetURLByShortURL(gomock.Any(), gomock.Any()).
					Return(entity.URL{
						ID:          "42",
						ShortURL:    "42",
						OriginalURL: "https://ya.ru/",
						DeletedFlag: true,
					}, nil)
			},
			wantStatusCode: http.StatusGone,
		},
		{
			name:     "Empty Original URL",
			shortURL: "42",
			storageCfg: func() {
				s.attr.storage.EXPECT().
					GetURLByShortURL(gomock.Any(), gomock.Any()).
					Return(entity.URL{}, nil)
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:     "Found",
			shortURL: "42",
			storageCfg: func() {
				s.attr.storage.EXPECT().
					GetURLByShortURL(gomock.Any(), gomock.Any()).
					Return(entity.URL{
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
			tt.storageCfg()
			resp, _ := test.Request(t, s.attr.ts, &test.RequestArgs{
				Method: "GET",
				Path:   fmt.Sprintf("/%s", tt.shortURL),
			})
			err := resp.Body.Close()
			assert.NoError(t, err)

			assert.Equal(t, tt.wantStatusCode, resp.StatusCode)
		})
	}
}

func Benchmark_Handler_Redirect(b *testing.B) {
	t := testing.T{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(&t)
	defer ctrl.Finish()

	err := config.Init()
	assert.NoError(b, err)
	log := logger.NewMock()
	storage := mockStorage.NewMockStorage(ctrl)
	urlsService := urls.New(ctx, log, storage, storage, storage)
	router := chi.NewRouter()
	h := New(log, storage, urlsService)
	router.Get("/{short_url}", h.Redirect)
	ts := httptest.NewServer(router)
	defer ts.Close()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		storage.EXPECT().GetURLByShortURL(gomock.Any(), gomock.Any()).Return(entity.URL{
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
		err = resp.Body.Close()
		assert.NoError(b, err)
	}
}
