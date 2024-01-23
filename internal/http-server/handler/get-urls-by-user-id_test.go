package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dsbasko/yandex-go-shortener/internal/config"
	"github.com/dsbasko/yandex-go-shortener/internal/entities"
	"github.com/dsbasko/yandex-go-shortener/internal/http-server/middlewares"
	"github.com/dsbasko/yandex-go-shortener/internal/jwt"
	"github.com/dsbasko/yandex-go-shortener/internal/storage"
	"github.com/dsbasko/yandex-go-shortener/internal/urls"
	"github.com/dsbasko/yandex-go-shortener/pkg/logger"
	"github.com/dsbasko/yandex-go-shortener/pkg/test"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_GetURLsByUserID(t *testing.T) {
	config.Init() //nolint:errcheck
	log := logger.NewMock()
	store := storage.NewMock(t)
	urlsService := urls.New(log, store)
	router := chi.NewRouter()
	h := New(log, store, urlsService)
	mw := middlewares.New(log)
	router.
		With(mw.JWT).
		Get("/api/user/urls", h.GetURLsByUserID)
	ts := httptest.NewServer(router)
	defer ts.Close()

	token, err := jwt.GenerateToken()
	assert.NoError(t, err)
	mockCookie := &http.Cookie{Name: jwt.CookieKey, Value: token}
	serviceError := errors.New("service error")

	tests := []struct {
		name           string
		storeCfg       func()
		cookie         *http.Cookie
		wantStatusCode int
		wantBody       func() string
	}{
		{
			name:           "Unauthorized",
			storeCfg:       func() {},
			cookie:         nil,
			wantStatusCode: http.StatusUnauthorized,
			wantBody: func() string {
				return ""
			},
		},
		{
			name: "Service Error",
			storeCfg: func() {
				store.EXPECT().
					GetURLsByUserID(gomock.Any(), gomock.Any()).
					Return(nil, serviceError)
			},
			cookie:         mockCookie,
			wantStatusCode: http.StatusBadRequest,
			wantBody: func() string {
				return ""
			},
		},
		{
			name: "Not Found",
			storeCfg: func() {
				store.EXPECT().
					GetURLsByUserID(gomock.Any(), gomock.Any()).
					Return([]entities.URL{}, nil)
			},
			cookie:         mockCookie,
			wantStatusCode: http.StatusNoContent,
			wantBody: func() string {
				return ""
			},
		},
		{
			name: "Success",
			storeCfg: func() {
				store.EXPECT().
					GetURLsByUserID(gomock.Any(), gomock.Any()).
					Return([]entities.URL{
						{
							ID:          "42",
							UserID:      "42",
							ShortURL:    "42",
							OriginalURL: "https://ya.ru",
						},
					}, nil)
			},
			cookie:         mockCookie,
			wantStatusCode: http.StatusOK,
			wantBody: func() string {
				respBytes, _ := json.Marshal([]entities.URL{
					{
						ID:          "42",
						UserID:      "42",
						ShortURL:    fmt.Sprintf("%s42", config.GetBaseURL()),
						OriginalURL: "https://ya.ru",
					},
				})
				return string(respBytes)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.storeCfg()
			resp, body := test.Request(t, ts, &test.RequestArgs{
				Method: "GET",
				Path:   "/api/user/urls",
				Cookie: tt.cookie,
			})
			defer resp.Body.Close()

			assert.Equal(t, tt.wantStatusCode, resp.StatusCode)
			assert.Equal(t, tt.wantBody(), body)
		})
	}
}

func BenchmarkHandler_GetURLsByUserID(b *testing.B) {
	config.Init() //nolint:errcheck
	log := logger.NewMock()
	store := storage.NewMock(&testing.T{})
	urlsService := urls.New(log, store)
	router := chi.NewRouter()
	h := New(log, store, urlsService)
	mw := middlewares.New(log)
	router.
		With(mw.JWT).
		Get("/api/user/urls", h.GetURLsByUserID)
	ts := httptest.NewServer(router)
	defer ts.Close()

	token, err := jwt.GenerateToken()
	assert.NoError(b, err)
	mockCookie := &http.Cookie{Name: jwt.CookieKey, Value: token}

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		store.EXPECT().GetURLsByUserID(gomock.Any(), gomock.Any()).Return([]entities.URL{}, nil)
		b.StartTimer()

		resp, _ := test.Request(&testing.T{}, ts, &test.RequestArgs{
			Method: "GET",
			Path:   "/api/user/urls",
			Cookie: mockCookie,
		})
		resp.Body.Close()
	}
}
