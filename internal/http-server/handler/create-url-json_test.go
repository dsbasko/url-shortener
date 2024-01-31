package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
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

func TestHandler_CreateURLJSON(t *testing.T) {
	config.Init() //nolint:errcheck
	log := logger.NewMock()
	store := storage.NewMock(t)
	urlsService := urls.New(log, store)
	router := chi.NewRouter()
	h := New(log, store, urlsService)
	mw := middlewares.New(log)
	router.
		With(mw.JWT).
		Post("/api/shorten", h.CreateURLJSON)
	ts := httptest.NewServer(router)
	defer ts.Close()

	serviceErr := errors.New("service error")

	tests := []struct {
		name           string
		contentType    string
		body           func() []byte
		storeCfg       func()
		wantStatusCode int
		wantBody       func() string
	}{
		{
			name:           "Wrong Content-Type",
			body:           func() []byte { return []byte("") },
			storeCfg:       func() {},
			wantStatusCode: http.StatusBadRequest,
			wantBody:       func() string { return "" },
		},
		{
			name:           "Empty Body",
			contentType:    "application/json",
			body:           func() []byte { return []byte("") },
			storeCfg:       func() {},
			wantStatusCode: http.StatusBadRequest,
			wantBody:       func() string { return "" },
		},
		{
			name: "Service Error",
			body: func() []byte {
				dtoBytes, _ := json.Marshal(api.CreateURLRequest{URL: "https://ya.ru/"})
				return dtoBytes
			},
			contentType: "application/json",
			storeCfg: func() {
				store.EXPECT().
					CreateURL(gomock.Any(), gomock.Any()).
					Return(entities.URL{}, false, serviceErr)
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
			storeCfg: func() {
				store.EXPECT().
					CreateURL(gomock.Any(), gomock.Any()).
					Return(entities.URL{
						ID:          "42",
						ShortURL:    "42",
						OriginalURL: "https://ya2.ru/",
					}, true, nil)
			},
			wantStatusCode: http.StatusCreated,
			wantBody: func() string {
				resBytes, _ := json.Marshal(api.CreateURLResponse{
					Result: fmt.Sprintf("%s42", config.GetBaseURL()),
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
			storeCfg: func() {
				store.EXPECT().
					CreateURL(gomock.Any(), gomock.Any()).
					Return(entities.URL{
						ID:          "42",
						ShortURL:    "42",
						OriginalURL: "https://ya3.ru/",
					}, false, nil)
			},
			wantStatusCode: http.StatusConflict,
			wantBody: func() string {
				resBytes, _ := json.Marshal(api.CreateURLResponse{
					Result: fmt.Sprintf("%s42", config.GetBaseURL()),
				})
				return fmt.Sprintf("%s\n", resBytes)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.storeCfg()
			resp, body := test.Request(t, ts, &test.RequestArgs{
				Method:      "POST",
				Path:        "/api/shorten",
				ContentType: tt.contentType,
				Body:        tt.body(),
			})
			defer resp.Body.Close()

			assert.Equal(t, tt.wantStatusCode, resp.StatusCode)
			assert.Equal(t, tt.wantBody(), body)
		})
	}
}

func BenchmarkHandler_CreateURLJSON(b *testing.B) {
	config.Init() //nolint:errcheck
	log := logger.NewMock()
	store := storage.NewMock(&testing.T{})
	urlsService := urls.New(log, store)
	router := chi.NewRouter()
	h := New(log, store, urlsService)
	mw := middlewares.New(log)
	router.
		With(mw.JWT).
		Post("/api/shorten", h.CreateURLJSON)
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
			Method:      "POST",
			Path:        "/api/shorten",
			ContentType: "application/json",
			Body:        []byte(`{"url":"https://ya.ru/"}`),
		})
		resp.Body.Close()
	}
}
