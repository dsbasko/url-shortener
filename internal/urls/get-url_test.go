package urls

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/dsbasko/yandex-go-shortener/internal/entities"
	"github.com/dsbasko/yandex-go-shortener/pkg/errors"
)

func (s *SuiteURLs) Test_GetURL() {
	t := s.T()

	type want struct {
		resp entities.URL
		err  error
	}

	tests := []struct {
		name     string
		shortURL string
		storeCfg func()
		want     want
	}{
		{
			name:     "Not Found",
			shortURL: "42",
			storeCfg: func() {
				s.attr.store.EXPECT().
					GetURLByShortURL(gomock.Any(), gomock.Any()).
					Return(entities.URL{}, s.attr.errNotFound)
			},
			want: want{
				resp: entities.URL{},
				err:  s.attr.errNotFound,
			},
		},
		{
			name:     "Found",
			shortURL: "42",
			storeCfg: func() {
				s.attr.store.EXPECT().
					GetURLByShortURL(gomock.Any(), gomock.Any()).
					Return(entities.URL{
						ID:          "42",
						ShortURL:    "42",
						OriginalURL: "https://ya.ru/",
					}, nil)
			},
			want: want{
				resp: entities.URL{
					ID:          "42",
					ShortURL:    "42",
					OriginalURL: "https://ya.ru/",
				},
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.storeCfg()
			resp, err := s.attr.service.GetURL(context.Background(), tt.shortURL)

			assert.Equal(t, tt.want.resp, resp)
			assert.Equal(t, tt.want.err, errors.UnwrapAll(err))
		})
	}
}
