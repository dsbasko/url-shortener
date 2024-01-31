package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/dsbasko/yandex-go-shortener/internal/config"
	"github.com/dsbasko/yandex-go-shortener/internal/entities"
	"github.com/dsbasko/yandex-go-shortener/internal/http-server/middlewares"
	"github.com/dsbasko/yandex-go-shortener/internal/storage"
	"github.com/dsbasko/yandex-go-shortener/internal/urls"
	"github.com/dsbasko/yandex-go-shortener/pkg/api"
	"github.com/dsbasko/yandex-go-shortener/pkg/logger"
	"github.com/dsbasko/yandex-go-shortener/pkg/test"
)

func TestHandler_CreateURLManyJSON(t *testing.T) {
	config.Init() //nolint:errcheck
	log := logger.NewMock()
	store := storage.NewMock(t)
	urlsService := urls.New(log, store)
	router := chi.NewRouter()
	h := New(log, store, urlsService)
	mw := middlewares.New(log)
	router.
		With(mw.JWT).
		Post("/api/shorten/batch", h.CreateURLsJSON)
	ts := httptest.NewServer(router)
	defer ts.Close()

	serviceErr := errors.New("service error")

	tests := []struct {
		name           string
		contentType    string
		body           func() []byte
		storeCfg       func()
		wantStatusCode int
		wantBody       []api.CreateURLsResponse
	}{
		{
			name:           "Wrong Content-Type",
			body:           func() []byte { return []byte("") },
			storeCfg:       func() {},
			wantStatusCode: http.StatusBadRequest,
			wantBody:       nil,
		},
		{
			name:           "Empty Body",
			contentType:    "application/json",
			body:           func() []byte { return []byte("") },
			storeCfg:       func() {},
			wantStatusCode: http.StatusBadRequest,
			wantBody:       nil,
		},
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
			storeCfg: func() {
				store.EXPECT().
					CreateURLs(gomock.Any(), gomock.Any()).
					Return(nil, serviceErr)
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
			storeCfg: func() {
				store.EXPECT().
					CreateURLs(gomock.Any(), gomock.Any()).
					Return([]entities.URL{}, nil)
			},
			wantStatusCode: http.StatusCreated,
			wantBody: []api.CreateURLsResponse{
				{
					ShortURL:      config.GetBaseURL(),
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
			storeCfg: func() {
				store.EXPECT().
					CreateURLs(gomock.Any(), gomock.Any()).
					Return([]entities.URL{}, nil)
			},
			wantStatusCode: http.StatusCreated,
			wantBody: []api.CreateURLsResponse{
				{
					ShortURL:      config.GetBaseURL(),
					CorrelationID: "1",
				},
				{
					ShortURL:      config.GetBaseURL(),
					CorrelationID: "2",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.storeCfg()
			resp, body := test.Request(t, ts, &test.RequestArgs{
				Method:      "POST",
				Path:        "/api/shorten/batch",
				ContentType: tt.contentType,
				Body:        tt.body(),
			})
			defer resp.Body.Close()

			if resp.ContentLength > 4 || tt.wantBody != nil {
				var bodyStruct []api.CreateURLsResponse
				err := json.Unmarshal([]byte(body), &bodyStruct)
				assert.NoError(t, err)
				assert.Equal(t, len(tt.wantBody), len(bodyStruct))
				assert.True(t, strings.Contains(bodyStruct[0].ShortURL, config.GetBaseURL()))
			}

			assert.Equal(t, tt.wantStatusCode, resp.StatusCode)
		})
	}
}

func BenchmarkHandler_CreateURLManyJSON(b *testing.B) {
	config.Init() //nolint:errcheck
	log := logger.NewMock()
	store := storage.NewMock(&testing.T{})
	urlsService := urls.New(log, store)
	router := chi.NewRouter()
	h := New(log, store, urlsService)
	mw := middlewares.New(log)
	router.
		With(mw.JWT).
		Post("/api/shorten/batch", h.CreateURLsJSON)
	ts := httptest.NewServer(router)
	defer ts.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		store.EXPECT().CreateURLs(gomock.Any(), gomock.Any()).Return([]entities.URL{}, nil)
		b.StartTimer()

		resp, _ := test.Request(&testing.T{}, ts, &test.RequestArgs{
			Method:      "POST",
			Path:        "/api/shorten/batch",
			ContentType: "application/json",
			Body: []byte(`[
				{"url":"https://ya1.ru/"},
				{"url":"https://ya2.ru/"},
				{"url":"https://ya3.ru/"},
				{"url":"https://ya4.ru/"},
				{"url":"https://ya5.ru/"},
				{"url":"https://ya6.ru/"},
				{"url":"https://ya7.ru/"},
				{"url":"https://ya8.ru/"},
				{"url":"https://ya9.ru/"},
				{"url":"https://ya10.ru/"},
				{"url":"https://ya11.ru/"},
				{"url":"https://ya12.ru/"},
				{"url":"https://ya13.ru/"},
				{"url":"https://ya14.ru/"},
				{"url":"https://ya15.ru/"},
			]`),
		})
		resp.Body.Close()
	}
}
