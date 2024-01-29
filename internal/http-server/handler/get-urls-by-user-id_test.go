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
				return fmt.Sprintf("%s\n", respBytes)
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
		store.EXPECT().GetURLsByUserID(gomock.Any(), gomock.Any()).Return([]entities.URL{
			{ID: "1", OriginalURL: "https://ya1.ru", ShortURL: "1"},
			{ID: "2", OriginalURL: "https://ya2.ru", ShortURL: "2"},
			{ID: "3", OriginalURL: "https://ya3.ru", ShortURL: "3"},
			{ID: "4", OriginalURL: "https://ya4.ru", ShortURL: "4"},
			{ID: "5", OriginalURL: "https://ya5.ru", ShortURL: "5"},
			{ID: "6", OriginalURL: "https://ya6.ru", ShortURL: "6"},
			{ID: "7", OriginalURL: "https://ya7.ru", ShortURL: "7"},
			{ID: "8", OriginalURL: "https://ya8.ru", ShortURL: "8"},
			{ID: "9", OriginalURL: "https://ya9.ru", ShortURL: "9"},
			{ID: "10", OriginalURL: "https://ya10.ru", ShortURL: "10"},
			{ID: "11", OriginalURL: "https://ya11.ru", ShortURL: "11"},
			{ID: "12", OriginalURL: "https://ya12.ru", ShortURL: "12"},
			{ID: "13", OriginalURL: "https://ya13.ru", ShortURL: "13"},
			{ID: "14", OriginalURL: "https://ya14.ru", ShortURL: "14"},
			{ID: "15", OriginalURL: "https://ya15.ru", ShortURL: "15"},
		}, nil)
		b.StartTimer()

		resp, _ := test.Request(&testing.T{}, ts, &test.RequestArgs{
			Method: "GET",
			Path:   "/api/user/urls",
			Cookie: mockCookie,
		})
		resp.Body.Close()
	}
}
