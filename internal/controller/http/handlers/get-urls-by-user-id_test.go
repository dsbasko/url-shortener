package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/dsbasko/url-shortener/internal/config"
	"github.com/dsbasko/url-shortener/internal/controller/http/middlewares"
	"github.com/dsbasko/url-shortener/internal/entity"
	mockStorage "github.com/dsbasko/url-shortener/internal/repository/storage/mocks"
	"github.com/dsbasko/url-shortener/internal/service/jwt"
	"github.com/dsbasko/url-shortener/internal/service/urls"
	"github.com/dsbasko/url-shortener/pkg/logger"
	"github.com/dsbasko/url-shortener/pkg/test"
)

func (s *SuiteHandlers) Test_GetURLsByUserID() {
	t := s.T()

	tests := []struct {
		name           string
		storageCfg     func()
		cookie         *http.Cookie
		wantStatusCode int
		wantBody       func() string
	}{
		{
			name:           "Unauthorized",
			storageCfg:     func() {},
			cookie:         nil,
			wantStatusCode: http.StatusUnauthorized,
			wantBody: func() string {
				return ""
			},
		},
		{
			name: "Service Error",
			storageCfg: func() {
				s.attr.storage.EXPECT().
					GetURLsByUserID(gomock.Any(), gomock.Any()).
					Return(nil, s.attr.errService)
			},
			cookie:         s.attr.cookie,
			wantStatusCode: http.StatusBadRequest,
			wantBody: func() string {
				return ""
			},
		},
		{
			name: "Not Found",
			storageCfg: func() {
				s.attr.storage.EXPECT().
					GetURLsByUserID(gomock.Any(), gomock.Any()).
					Return([]entity.URL{}, nil)
			},
			cookie:         s.attr.cookie,
			wantStatusCode: http.StatusNoContent,
			wantBody: func() string {
				return ""
			},
		},
		{
			name: "Success",
			storageCfg: func() {
				s.attr.storage.EXPECT().
					GetURLsByUserID(gomock.Any(), gomock.Any()).
					Return([]entity.URL{
						{
							ID:          "42",
							UserID:      "42",
							ShortURL:    "42",
							OriginalURL: "https://ya.ru",
						},
					}, nil)
			},
			cookie:         s.attr.cookie,
			wantStatusCode: http.StatusOK,
			wantBody: func() string {
				respBytes, _ := json.Marshal([]entity.URL{
					{
						ID:          "42",
						UserID:      "42",
						ShortURL:    fmt.Sprintf("%s42", config.BaseURL()),
						OriginalURL: "https://ya.ru",
					},
				})
				return fmt.Sprintf("%s\n", respBytes)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.storageCfg()
			resp, body := test.Request(t, s.attr.ts, &test.RequestArgs{
				Method: "GET",
				Path:   "/api/user/urls",
				Cookie: tt.cookie,
			})
			err := resp.Body.Close()
			assert.NoError(t, err)

			assert.Equal(t, tt.wantStatusCode, resp.StatusCode)
			assert.Equal(t, tt.wantBody(), body)
		})
	}
}

func Benchmark_Handler_GetURLsByUserID(b *testing.B) {
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
		storage.EXPECT().GetURLsByUserID(gomock.Any(), gomock.Any()).Return([]entity.URL{
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
		err = resp.Body.Close()
		assert.NoError(b, err)
	}
}
