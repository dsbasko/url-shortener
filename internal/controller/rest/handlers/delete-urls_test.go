package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/dsbasko/yandex-go-shortener/internal/config"
	"github.com/dsbasko/yandex-go-shortener/internal/controller/rest/middlewares"
	"github.com/dsbasko/yandex-go-shortener/internal/entities"
	"github.com/dsbasko/yandex-go-shortener/internal/repository/storage"
	"github.com/dsbasko/yandex-go-shortener/internal/service/jwt"
	"github.com/dsbasko/yandex-go-shortener/internal/service/urls"
	"github.com/dsbasko/yandex-go-shortener/pkg/logger"
	"github.com/dsbasko/yandex-go-shortener/pkg/test"
)

func (s *SuiteHandlers) Test_DeleteURLs() {
	t := s.T()

	tests := []struct {
		name           string
		contentType    string
		body           []byte
		storeCfg       func()
		cookie         *http.Cookie
		wantStatusCode int
	}{
		{
			name:           "Unauthorized",
			body:           []byte(""),
			storeCfg:       func() {},
			cookie:         nil,
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name:           "JSON Marshal Error",
			contentType:    "application/json",
			body:           []byte("42[],,"),
			storeCfg:       func() {},
			cookie:         s.attr.cookie,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:        "Success",
			body:        []byte(`["42"]`),
			contentType: "application/json",
			storeCfg: func() {
				s.attr.store.EXPECT().DeleteURLs(gomock.Any(), gomock.Any()).Return(nil, nil)
			},
			cookie:         s.attr.cookie,
			wantStatusCode: http.StatusAccepted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.storeCfg()
			resp, _ := test.Request(t, s.attr.ts, &test.RequestArgs{
				Method:      "DELETE",
				Path:        "/api/user/urls",
				ContentType: tt.contentType,
				Body:        tt.body,
				Cookie:      tt.cookie,
			})
			err := resp.Body.Close()
			assert.NoError(t, err)

			assert.Equal(t, tt.wantStatusCode, resp.StatusCode)
		})
	}
}

func BenchmarkHandler_DeleteURLs(b *testing.B) {
	err := config.Init()
	assert.NoError(b, err)
	log := logger.NewMock()
	store := storage.NewMock(&testing.T{})
	urlsService := urls.New(log, store)
	router := chi.NewRouter()
	h := New(log, store, urlsService)
	mw := middlewares.New(log)
	router.
		With(mw.JWT).
		Delete("/api/user/urls", h.DeleteURLs)
	ts := httptest.NewServer(router)
	defer ts.Close()

	token, err := jwt.GenerateToken()
	assert.NoError(b, err)
	mockCookie := &http.Cookie{Name: jwt.CookieKey, Value: token}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		store.EXPECT().DeleteURLs(gomock.Any(), gomock.Any()).Return([]entities.URL{
			{ID: "42", OriginalURL: "https://ya42.ru", ShortURL: "42"},
			{ID: "43", OriginalURL: "https://ya43.ru", ShortURL: "43"},
			{ID: "44", OriginalURL: "https://ya44.ru", ShortURL: "44"},
			{ID: "45", OriginalURL: "https://ya45.ru", ShortURL: "45"},
			{ID: "46", OriginalURL: "https://ya46.ru", ShortURL: "46"},
			{ID: "47", OriginalURL: "https://ya47.ru", ShortURL: "47"},
			{ID: "48", OriginalURL: "https://ya48.ru", ShortURL: "48"},
			{ID: "49", OriginalURL: "https://ya49.ru", ShortURL: "49"},
			{ID: "50", OriginalURL: "https://ya50.ru", ShortURL: "50"},
		}, nil)
		b.StartTimer()

		resp, _ := test.Request(&testing.T{}, ts, &test.RequestArgs{
			Method:      "DELETE",
			Path:        "/api/user/urls",
			ContentType: "application/json",
			Body:        []byte(`["42", "43", "44", "45", "46", "47", "48", "49", "50"]`),
			Cookie:      mockCookie,
		})
		err = resp.Body.Close()
		assert.NoError(b, err)
	}
}
