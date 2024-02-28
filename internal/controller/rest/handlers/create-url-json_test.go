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

	"github.com/dsbasko/yandex-go-shortener/internal/config"
	"github.com/dsbasko/yandex-go-shortener/internal/controller/rest/middlewares"
	"github.com/dsbasko/yandex-go-shortener/internal/entity"
	mockStorage "github.com/dsbasko/yandex-go-shortener/internal/repository/storage/mocks"
	"github.com/dsbasko/yandex-go-shortener/internal/service/urls"
	"github.com/dsbasko/yandex-go-shortener/pkg/api"
	"github.com/dsbasko/yandex-go-shortener/pkg/logger"
	"github.com/dsbasko/yandex-go-shortener/pkg/test"
)

func (s *SuiteHandlers) Test_CreateURL_JSON() {
	t := s.T()

	tests := []struct {
		name           string
		contentType    string
		body           func() []byte
		storageCfg     func()
		wantStatusCode int
		wantBody       func() string
	}{
		{
			name: "Service Error",
			body: func() []byte {
				dtoBytes, _ := json.Marshal(api.CreateURLRequest{URL: "https://ya.ru/"})
				return dtoBytes
			},
			contentType: "application/json",
			storageCfg: func() {
				s.attr.urlsMutator.EXPECT().
					CreateURL(gomock.Any(), gomock.Any()).
					Return(entity.URL{}, false, s.attr.errService)
			},
			wantStatusCode: http.StatusBadRequest,
			wantBody:       func() string { return "" },
		},
		{
			name: "Success Unique",
			body: func() []byte {
				dtoBytes, _ := json.Marshal(api.CreateURLRequest{URL: "https://ya.ru/"})
				return dtoBytes
			},
			contentType: "application/json",
			storageCfg: func() {
				s.attr.urlsMutator.EXPECT().
					CreateURL(gomock.Any(), gomock.Any()).
					Return(entity.URL{
						ID:          "42",
						ShortURL:    "42",
						OriginalURL: "https://ya2.ru/",
					}, true, nil)
			},
			wantStatusCode: http.StatusCreated,
			wantBody: func() string {
				resBytes, _ := json.Marshal(api.CreateURLResponse{
					Result: fmt.Sprintf("%s42", config.BaseURL()),
				})
				return fmt.Sprintf("%s\n", resBytes)
			},
		},
		{
			name: "Success NotUnique",
			body: func() []byte {
				dtoBytes, _ := json.Marshal(api.CreateURLRequest{URL: "https://ya.ru/"})
				return dtoBytes
			},
			contentType: "application/json",
			storageCfg: func() {
				s.attr.urlsMutator.EXPECT().
					CreateURL(gomock.Any(), gomock.Any()).
					Return(entity.URL{
						ID:          "42",
						ShortURL:    "42",
						OriginalURL: "https://ya3.ru/",
					}, false, nil)
			},
			wantStatusCode: http.StatusConflict,
			wantBody: func() string {
				resBytes, _ := json.Marshal(api.CreateURLResponse{
					Result: fmt.Sprintf("%s42", config.BaseURL()),
				})
				return fmt.Sprintf("%s\n", resBytes)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.storageCfg()
			resp, body := test.Request(t, s.attr.ts, &test.RequestArgs{
				Method:      "POST",
				Path:        "/api/shorten",
				ContentType: tt.contentType,
				Body:        tt.body(),
			})
			err := resp.Body.Close()
			assert.NoError(t, err)

			assert.Equal(t, tt.wantStatusCode, resp.StatusCode)
			assert.Equal(t, tt.wantBody(), body)
		})
	}
}

func Benchmark_Handler_CreateURLJSON(b *testing.B) {
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
		Post("/api/shorten", h.CreateURLJSON)
	ts := httptest.NewServer(router)
	defer ts.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		storage.EXPECT().CreateURL(gomock.Any(), gomock.Any()).Return(entity.URL{
			ID:          "42",
			ShortURL:    "42",
			OriginalURL: "https://ya.ru/",
			UserID:      "42",
		}, false, nil)
		b.StartTimer()

		resp, _ := test.Request(&testing.T{}, ts, &test.RequestArgs{
			Method:      "POST",
			Path:        "/api/shorten",
			ContentType: "application/json",
			Body:        []byte(`{"url":"https://ya.ru/"}`),
		})
		err = resp.Body.Close()
		assert.NoError(b, err)
	}
}
