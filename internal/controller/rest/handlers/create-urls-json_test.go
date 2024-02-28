package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/dsbasko/yandex-go-shortener/internal/config"
	"github.com/dsbasko/yandex-go-shortener/internal/controller/rest/middlewares"
	"github.com/dsbasko/yandex-go-shortener/internal/entity"
	mockStorage "github.com/dsbasko/yandex-go-shortener/internal/repository/storage/mocks"
	"github.com/dsbasko/yandex-go-shortener/internal/service/jwt"
	"github.com/dsbasko/yandex-go-shortener/internal/service/urls"
	"github.com/dsbasko/yandex-go-shortener/pkg/api"
	"github.com/dsbasko/yandex-go-shortener/pkg/logger"
	"github.com/dsbasko/yandex-go-shortener/pkg/test"
)

func (s *SuiteHandlers) Test_CreateURLs_JSON() {
	t := s.T()

	tests := []struct {
		name           string
		contentType    string
		body           func() []byte
		storageCfg     func()
		wantStatusCode int
		wantBody       []api.CreateURLsResponse
	}{
		{
			name: "Service Error",
			body: func() []byte {
				dtoBytes, _ := json.Marshal([]api.CreateURLsRequest{
					{
						OriginalURL:   "https://ya.ru/",
						CorrelationID: "1",
					},
				})
				return dtoBytes
			},
			contentType: "application/json",
			storageCfg: func() {
				s.attr.urlsMutator.EXPECT().
					CreateURLs(gomock.Any(), gomock.Any()).
					Return(nil, s.attr.errService)
			},
			wantStatusCode: http.StatusBadRequest,
			wantBody:       nil,
		},
		{
			name: "Success Once",
			body: func() []byte {
				dtoBytes, _ := json.Marshal([]api.CreateURLsRequest{
					{
						OriginalURL:   "https://ya.ru/",
						CorrelationID: "1",
					},
				})
				return dtoBytes
			},
			contentType: "application/json",
			storageCfg: func() {
				s.attr.urlsMutator.EXPECT().
					CreateURLs(gomock.Any(), gomock.Any()).
					Return([]entity.URL{}, nil)
			},
			wantStatusCode: http.StatusCreated,
			wantBody: []api.CreateURLsResponse{
				{
					ShortURL:      config.BaseURL(),
					CorrelationID: "1",
				},
			},
		},
		{
			name: "Success Many",
			body: func() []byte {
				dtoBytes, _ := json.Marshal([]api.CreateURLsRequest{
					{
						OriginalURL:   "https://ya.ru/",
						CorrelationID: "1",
					},
					{
						OriginalURL:   "https://yandex.ru/",
						CorrelationID: "2",
					},
				})
				return dtoBytes
			},
			contentType: "application/json",
			storageCfg: func() {
				s.attr.urlsMutator.EXPECT().
					CreateURLs(gomock.Any(), gomock.Any()).
					Return([]entity.URL{}, nil)
			},
			wantStatusCode: http.StatusCreated,
			wantBody: []api.CreateURLsResponse{
				{
					ShortURL:      config.BaseURL(),
					CorrelationID: "1",
				},
				{
					ShortURL:      config.BaseURL(),
					CorrelationID: "2",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.storageCfg()
			resp, body := test.Request(t, s.attr.ts, &test.RequestArgs{
				Method:      "POST",
				Path:        "/api/shorten/batch",
				ContentType: tt.contentType,
				Body:        tt.body(),
			})
			err := resp.Body.Close()
			assert.NoError(t, err)

			if resp.ContentLength > 4 || tt.wantBody != nil {
				var bodyStruct []api.CreateURLsResponse
				err = json.Unmarshal([]byte(body), &bodyStruct)
				assert.NoError(t, err)
				assert.Equal(t, len(tt.wantBody), len(bodyStruct))
				assert.True(t, strings.Contains(bodyStruct[0].ShortURL, config.BaseURL()))
			}

			assert.Equal(t, tt.wantStatusCode, resp.StatusCode)
		})
	}
}

func Benchmark_Handler_CreateURLsJSON(b *testing.B) {
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
	router.With(mw.JWT).Post("/api/shorten/batch", h.CreateURLsJSON)
	ts := httptest.NewServer(router)
	defer ts.Close()

	token, err := jwt.GenerateToken()
	assert.NoError(b, err)
	mockCookie := &http.Cookie{Name: jwt.CookieKey, Value: token}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		storage.EXPECT().CreateURLs(gomock.Any(), gomock.Any()).Return([]entity.URL{}, nil)
		body := func() []byte {
			dto := []api.CreateURLsRequest{
				{CorrelationID: "1", OriginalURL: "http://ya1.ru/"},
				{CorrelationID: "2", OriginalURL: "http://ya2.ru/"},
				{CorrelationID: "3", OriginalURL: "http://ya3.ru/"},
				{CorrelationID: "4", OriginalURL: "http://ya4.ru/"},
				{CorrelationID: "5", OriginalURL: "http://ya5.ru/"},
				{CorrelationID: "6", OriginalURL: "http://ya6.ru/"},
				{CorrelationID: "7", OriginalURL: "http://ya7.ru/"},
				{CorrelationID: "8", OriginalURL: "http://ya8.ru/"},
				{CorrelationID: "9", OriginalURL: "http://ya9.ru/"},
			}
			dtoBytes, _ := json.Marshal(dto)
			return dtoBytes
		}()
		b.StartTimer()

		resp, _ := test.Request(&t, ts, &test.RequestArgs{
			Method:      "POST",
			Path:        "/api/shorten/batch",
			ContentType: "application/json",
			Cookie:      mockCookie,
			Body:        body,
		})
		err = resp.Body.Close()
		assert.NoError(b, err)
	}
}
