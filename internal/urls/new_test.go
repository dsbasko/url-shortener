package urls

import (
	"context"
	"fmt"
	"testing"

	"github.com/dsbasko/yandex-go-shortener/internal/config"
	"github.com/dsbasko/yandex-go-shortener/internal/jwt"
	"github.com/dsbasko/yandex-go-shortener/internal/storage/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/dsbasko/yandex-go-shortener/internal/storage"
	"github.com/dsbasko/yandex-go-shortener/pkg/logger"
)

type SuiteURLs struct {
	*suite.Suite

	attr struct {
		log          *logger.Logger
		store        *mock.MockStorage
		service      URLs
		errStore     error
		errNotFound  error
		ctxWithToken context.Context
	}
}

func (s *SuiteURLs) SetupSuite() {
	t := s.T()

	config.Init() //nolint:errcheck
	s.attr.log = logger.NewMock()
	s.attr.store = storage.NewMock(t)
	s.attr.service = New(s.attr.log, s.attr.store)
	s.attr.errStore = fmt.Errorf("storage error")
	s.attr.errNotFound = fmt.Errorf("not found")

	token, err := jwt.GenerateToken()
	assert.NoError(t, err)
	s.attr.ctxWithToken = context.WithValue(context.Background(), jwt.ContextKey, token)
}

func (s *SuiteURLs) Test_New() {
	t := s.T()

	t.Run("Success", func(t *testing.T) {
		assert.NotNil(t, s.attr.service)
		assert.Equal(
			t,
			URLs{log: s.attr.log, storage: s.attr.store},
			s.attr.service,
		)
	})
}

func Test_URLs_Service(t *testing.T) {
	suite.Run(t, &SuiteURLs{
		Suite: new(suite.Suite),
	})
}
