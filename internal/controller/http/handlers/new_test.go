package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/dsbasko/url-shortener/internal/config"
	"github.com/dsbasko/url-shortener/internal/controller/http/middlewares"
	mockStorage "github.com/dsbasko/url-shortener/internal/repository/storage/mocks"
	"github.com/dsbasko/url-shortener/internal/service/jwt"
	"github.com/dsbasko/url-shortener/internal/service/urls"
	"github.com/dsbasko/url-shortener/pkg/logger"
)

type SuiteHandlers struct {
	*suite.Suite

	attr struct {
		log          *logger.Logger
		urls         urls.URLs
		storage      *mockStorage.MockStorage
		handler      Handler
		errService   error
		errNotFound  error
		ctxWithToken context.Context
		ts           *httptest.Server
		cookie       *http.Cookie
	}
}

func (s *SuiteHandlers) SetupSuite() {
	t := s.T()
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config.InitMock()
	s.attr.log = logger.NewMock()
	s.attr.storage = mockStorage.NewMockStorage(ctrl)
	s.attr.urls = urls.New(ctx, s.attr.log, s.attr.storage, s.attr.storage, s.attr.storage)
	router := chi.NewRouter()
	s.attr.handler = New(s.attr.log, s.attr.storage, s.attr.urls)
	mw := middlewares.New(s.attr.log)
	token, err := jwt.GenerateToken()
	assert.NoError(t, err)
	s.attr.cookie = &http.Cookie{Name: jwt.CookieKey, Value: token}

	// Роуты
	router.MethodNotAllowed(s.attr.handler.BadRequest)
	router.Get("/ping", s.attr.handler.Ping)
	router.Get("/{short_url}", s.attr.handler.Redirect)
	router.With(mw.JWT).Post("/api/shorten", s.attr.handler.CreateURLJSON)
	router.With(mw.JWT).Post("/", s.attr.handler.CreateURLTextPlain)
	router.With(mw.JWT).Post("/api/shorten/batch", s.attr.handler.CreateURLsJSON)
	router.With(mw.JWT).Get("/api/user/urls", s.attr.handler.GetURLsByUserID)
	router.With(mw.JWT).Delete("/api/user/urls", s.attr.handler.DeleteURLs)
	router.With(mw.JWT).Get("/api/internal/stats", s.attr.handler.Stats)

	s.attr.ts = httptest.NewServer(router)
	s.attr.errService = errors.New("service error")
}

func (s *SuiteHandlers) TearDownSuite() {
	s.attr.ts.Close()
}

func (s *SuiteHandlers) Test_New() {
	t := s.T()

	t.Run("Success", func(t *testing.T) {
		assert.NotNil(t, s.attr.handler)
		assert.Equal(
			t,
			Handler{log: s.attr.log, pinger: s.attr.storage, urls: s.attr.urls},
			s.attr.handler,
		)
	})
}

func Test_Handlers_Controller(t *testing.T) {
	suite.Run(t, &SuiteHandlers{
		Suite: new(suite.Suite),
	})
}
