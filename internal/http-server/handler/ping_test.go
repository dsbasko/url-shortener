package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dsbasko/yandex-go-shortener/internal/storage"
	"github.com/dsbasko/yandex-go-shortener/internal/urls"
	"github.com/dsbasko/yandex-go-shortener/pkg/logger"
	"github.com/dsbasko/yandex-go-shortener/pkg/test"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Ping(t *testing.T) {
	log := logger.NewMock()
	store := storage.NewMock(t)
	urlsService := urls.New(log, store)
	router := chi.NewRouter()
	h := New(log, store, urlsService)
	router.Get("/", h.Ping)
	ts := httptest.NewServer(router)
	defer ts.Close()

	tests := []struct {
		name           string
		storeCfg       func()
		wantStatusCode int
		wantBody       func() string
	}{
		{
			name: "Error",
			storeCfg: func() {
				store.EXPECT().Ping(gomock.Any()).Return(errors.New(""))
			},
			wantStatusCode: http.StatusInternalServerError,
			wantBody:       func() string { return "" },
		},
		{
			name: "Success",
			storeCfg: func() {
				store.EXPECT().Ping(gomock.Any()).Return(nil)
			},
			wantStatusCode: http.StatusOK,
			wantBody:       func() string { return "pong" },
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.storeCfg()
			resp, body := test.Request(t, ts, &test.RequestArgs{
				Method: "GET",
				Path:   "/",
			})
			defer resp.Body.Close()

			assert.Equal(t, tt.wantStatusCode, resp.StatusCode)
			assert.Equal(t, tt.wantBody(), body)
		})
	}
}
