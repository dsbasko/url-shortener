package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dsbasko/yandex-go-shortener/internal/config"
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

func TestHandler_DeleteURLs(t *testing.T) {
	config.Init() //nolint:errcheck
	log := logger.NewMock()
	store := storage.NewMock(t)
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
	assert.NoError(t, err)
	mockCookie := &http.Cookie{Name: jwt.CookieKey, Value: token}

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
			name:           "Empty Body",
			contentType:    "application/json",
			body:           []byte(""),
			storeCfg:       func() {},
			cookie:         mockCookie,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "JSON Marshal Error",
			contentType:    "application/json",
			body:           []byte("42[],,"),
			storeCfg:       func() {},
			cookie:         mockCookie,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:        "Success",
			body:        []byte(`["42"]`),
			contentType: "application/json",
			storeCfg: func() {
				store.EXPECT().DeleteURLs(gomock.Any(), gomock.Any()).Return(nil, nil)
			},
			cookie:         mockCookie,
			wantStatusCode: http.StatusAccepted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.storeCfg()
			resp, _ := test.Request(t, ts, &test.RequestArgs{
				Method:      "DELETE",
				Path:        "/api/user/urls",
				ContentType: tt.contentType,
				Body:        tt.body,
				Cookie:      tt.cookie,
			})
			defer resp.Body.Close()
			assert.Equal(t, tt.wantStatusCode, resp.StatusCode)
		})
	}
}
