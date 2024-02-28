package urls

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/dsbasko/yandex-go-shortener/internal/config"
	"github.com/dsbasko/yandex-go-shortener/internal/service/jwt"
	mockUrls "github.com/dsbasko/yandex-go-shortener/internal/service/urls/mocks"
	"github.com/dsbasko/yandex-go-shortener/pkg/logger"
)

type SuiteURLs struct {
	*suite.Suite

	attr struct {
		log          *logger.Logger
		urlProvider  *mockUrls.MockProvider
		urlMutator   *mockUrls.MockMutator
		urlAnalyzer  *mockUrls.MockAnalyzer
		service      URLs
		errStorage   error
		errNotFound  error
		ctxWithToken context.Context
		deleteTasks  chan map[string][]string
	}
}

func (s *SuiteURLs) SetupSuite() {
	t := s.T()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	err := config.Init()
	assert.NoError(t, err)
	s.attr.log = logger.NewMock()
	s.attr.urlProvider = mockUrls.NewMockProvider(ctrl)
	s.attr.urlMutator = mockUrls.NewMockMutator(ctrl)
	s.attr.urlAnalyzer = mockUrls.NewMockAnalyzer(ctrl)
	s.attr.service = New(ctx, s.attr.log, s.attr.urlProvider, s.attr.urlMutator, s.attr.urlAnalyzer)
	s.attr.errStorage = fmt.Errorf("storage error")
	s.attr.errNotFound = fmt.Errorf("not found")

	token, err := jwt.GenerateToken()
	assert.NoError(t, err)
	s.attr.ctxWithToken = context.WithValue(context.Background(), jwt.ContextKey, token)
}

func (s *SuiteURLs) Test_New() {
	t := s.T()

	t.Run("Success", func(t *testing.T) {
		mockService := URLs{
			log:         s.attr.log,
			urlProvider: s.attr.urlProvider,
			urlMutator:  s.attr.urlMutator,
		}

		assert.NotNil(t, s.attr.service)
		assert.EqualExportedValues(t, mockService, s.attr.service)
	})
}

func Test_URLs_Service(t *testing.T) {
	suite.Run(t, &SuiteURLs{
		Suite: new(suite.Suite),
	})
}
