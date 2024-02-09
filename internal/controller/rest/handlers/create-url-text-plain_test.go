package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/dsbasko/yandex-go-shortener/internal/config"
	"github.com/dsbasko/yandex-go-shortener/internal/controller/rest/middlewares"
	"github.com/dsbasko/yandex-go-shortener/internal/entity"
	mockStorage "github.com/dsbasko/yandex-go-shortener/internal/repository/storage/mocks"
	"github.com/dsbasko/yandex-go-shortener/internal/service/urls"
	"github.com/dsbasko/yandex-go-shortener/pkg/logger"
	"github.com/dsbasko/yandex-go-shortener/pkg/test"
)

func (s *SuiteHandlers) Test_CreateURL_TextPlain() {
	t := s.T()
	tests := []struct {
		name           string
		body           func() []byte
		storeCfg       func()
		wantStatusCode int
		wantBody       func() string
	}{
		{
			name: "Service Error",
			body: func() []byte { return []byte("https://ya.ru/") },
			storeCfg: func() {
				s.attr.urlsMutator.EXPECT().
					CreateURL(gomock.Any(), gomock.Any()).
					Return(entity.URL{}, false, s.attr.errService)
			},
			wantStatusCode: http.StatusBadRequest,
			wantBody:       func() string { return "" },
		},
		{
			name: "Success Unique",
			body: func() []byte { return []byte("https://ya.ru/") },
			storeCfg: func() {
				s.attr.urlsMutator.EXPECT().
					CreateURL(gomock.Any(), gomock.Any()).
					Return(entity.URL{
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
				s.attr.urlsMutator.EXPECT().
					CreateURL(gomock.Any(), gomock.Any()).
					Return(entity.URL{
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
			resp, body := test.Request(t, s.attr.ts, &test.RequestArgs{
				Method: "POST",
				Path:   "/",
				Body:   tt.body(),
			})
			err := resp.Body.Close()
			assert.NoError(t, err)

			assert.Equal(t, tt.wantStatusCode, resp.StatusCode)
			assert.Equal(t, tt.wantBody(), body)
		})
	}
}

func BenchmarkHandler_CreateURLOnceTextPlain(b *testing.B) {
	t := testing.T{}
	ctrl := gomock.NewController(&t)
	defer ctrl.Finish()

	err := config.Init()
	assert.NoError(b, err)
	log := logger.NewMock()
	store := mockStorage.NewMockStorage(ctrl)
	urlsService := urls.New(log, store, store)
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
			store.EXPECT().CreateURL(gomock.Any(), gomock.Any()).Return(entity.URL{
				ID:          "42",
				ShortURL:    "42",
				OriginalURL: "https://ya.ru/",
				UserID:      "42",
			}, false, nil)
		} else {
			store.EXPECT().CreateURL(gomock.Any(), gomock.Any()).Return(entity.URL{
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
		err = resp.Body.Close()
		assert.NoError(b, err)
	}
}
