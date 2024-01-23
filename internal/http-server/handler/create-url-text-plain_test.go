package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dsbasko/yandex-go-shortener/internal/config"
	"github.com/dsbasko/yandex-go-shortener/internal/entities"
	"github.com/dsbasko/yandex-go-shortener/internal/http-server/middlewares"
	"github.com/dsbasko/yandex-go-shortener/internal/storage"
	"github.com/dsbasko/yandex-go-shortener/internal/urls"
	"github.com/dsbasko/yandex-go-shortener/pkg/logger"
	"github.com/dsbasko/yandex-go-shortener/pkg/test"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_CreateURLOnceTextPlain(t *testing.T) {
	config.Init() //nolint:errcheck
	log := logger.NewMock()
	store := storage.NewMock(t)
	urlsService := urls.New(log, store)
	router := chi.NewRouter()
	h := New(log, store, urlsService)
	mw := middlewares.New(log)
	router.
		With(mw.JWT).
		Post("/", h.CreateURLTextPlain)
	ts := httptest.NewServer(router)
	defer ts.Close()

	serviceErr := errors.New("service error")

	tests := []struct {
		name           string
		body           func() []byte
		storeCfg       func()
		wantStatusCode int
		wantBody       func() string
	}{
		{
			name:           "Empty Body",
			body:           func() []byte { return []byte("") },
			storeCfg:       func() {},
			wantStatusCode: http.StatusBadRequest,
			wantBody:       func() string { return "" },
		},
		{
			name: "Service Error",
			body: func() []byte { return []byte("https://ya.ru/") },
			storeCfg: func() {
				store.EXPECT().
					CreateURL(gomock.Any(), gomock.Any()).
					Return(entities.URL{}, false, serviceErr)
			},
			wantStatusCode: http.StatusInternalServerError,
			wantBody:       func() string { return "" },
		},
		{
			name: "Success Unique",
			body: func() []byte { return []byte("https://ya.ru/") },
			storeCfg: func() {
				store.EXPECT().
					CreateURL(gomock.Any(), gomock.Any()).
					Return(entities.URL{
						ID:          "42",
						ShortURL:    "42",
						OriginalURL: "https://ya.ru/",
					}, true, nil)
			},
			wantStatusCode: http.StatusCreated,
			wantBody: func() string {
				return fmt.Sprintf("%s42", config.GetBaseURL())
			},
		},
		{
			name: "Success NotUnique",
			body: func() []byte { return []byte("https://ya.ru/") },
			storeCfg: func() {
				store.EXPECT().
					CreateURL(gomock.Any(), gomock.Any()).
					Return(entities.URL{
						ID:          "42",
						ShortURL:    "42",
						OriginalURL: "https://ya.ru/",
					}, false, nil)
			},
			wantStatusCode: http.StatusConflict,
			wantBody: func() string {
				return fmt.Sprintf("%s42", config.GetBaseURL())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.storeCfg()
			resp, body := test.Request(t, ts, &test.RequestArgs{
				Method: "POST",
				Path:   "/",
				Body:   tt.body(),
			})
			defer resp.Body.Close()

			assert.Equal(t, tt.wantStatusCode, resp.StatusCode)
			assert.Equal(t, tt.wantBody(), body)
		})
	}
}

func BenchmarkHandler_CreateURLOnceTextPlain(b *testing.B) {
	config.Init() //nolint:errcheck
	log := logger.NewMock()
	store := storage.NewMock(&testing.T{})
	urlsService := urls.New(log, store)
	router := chi.NewRouter()
	h := New(log, store, urlsService)
	mw := middlewares.New(log)
	router.
		With(mw.JWT).
		Post("/", h.CreateURLTextPlain)
	ts := httptest.NewServer(router)
	defer ts.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		if i%2 == 0 {
			store.EXPECT().CreateURL(gomock.Any(), gomock.Any()).Return(entities.URL{
				ID:          "42",
				ShortURL:    "42",
				OriginalURL: "https://ya.ru/",
				UserID:      "42",
			}, false, nil)
		} else {
			store.EXPECT().CreateURL(gomock.Any(), gomock.Any()).Return(entities.URL{
				ID:          "42",
				ShortURL:    "42",
				OriginalURL: "https://ya.ru/",
				UserID:      "42",
			}, true, nil)
		}
		b.StartTimer()

		resp, _ := test.Request(&testing.T{}, ts, &test.RequestArgs{
			Method: "POST",
			Path:   "/",
			Body:   []byte("https://ya.ru/"),
		})
		resp.Body.Close()
	}
}
