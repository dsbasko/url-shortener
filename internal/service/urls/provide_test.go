package urls

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/dsbasko/yandex-go-shortener/internal/config"
	"github.com/dsbasko/yandex-go-shortener/internal/entity"
	"github.com/dsbasko/yandex-go-shortener/pkg/errors"
)

func (s *SuiteURLs) Test_GetURL() {
	t := s.T()

	type want struct {
		resp entity.URL
		err  error
	}

	tests := []struct {
		name       string
		shortURL   string
		storageCfg func()
		want       want
	}{
		{
			name:     "Not Found",
			shortURL: "42",
			storageCfg: func() {
				s.attr.urlProvider.EXPECT().
					GetURLByShortURL(gomock.Any(), gomock.Any()).
					Return(entity.URL{}, s.attr.errNotFound)
			},
			want: want{
				resp: entity.URL{},
				err:  s.attr.errNotFound,
			},
		},
		{
			name:     "Found",
			shortURL: "42",
			storageCfg: func() {
				s.attr.urlProvider.EXPECT().
					GetURLByShortURL(gomock.Any(), gomock.Any()).
					Return(entity.URL{
						ID:          "42",
						ShortURL:    "42",
						OriginalURL: "https://ya.ru/",
					}, nil)
			},
			want: want{
				resp: entity.URL{
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
			tt.storageCfg()
			resp, err := s.attr.service.GetURL(context.Background(), tt.shortURL)

			assert.Equal(t, tt.want.resp, resp)
			assert.Equal(t, tt.want.err, errors.UnwrapAll(err))
		})
	}
}

func (s *SuiteURLs) Test_GetURLsByUserID() {
	t := s.T()

	type want struct {
		resp []entity.URL
		err  error
	}

	tests := []struct {
		name       string
		userID     string
		storageCfg func()
		want       want
	}{
		{
			name:   "Not Found",
			userID: "42",
			storageCfg: func() {
				s.attr.urlProvider.EXPECT().
					GetURLsByUserID(gomock.Any(), gomock.Any()).
					Return([]entity.URL{}, s.attr.errNotFound)
			},
			want: want{
				resp: []entity.URL{},
				err:  s.attr.errNotFound,
			},
		},
		{
			name:   "Found",
			userID: "42",
			storageCfg: func() {
				s.attr.urlProvider.EXPECT().
					GetURLsByUserID(gomock.Any(), gomock.Any()).
					Return([]entity.URL{
						{
							ID:          "42",
							ShortURL:    "42",
							OriginalURL: "https://ya.ru/",
						},
					}, nil)
			},
			want: want{
				resp: []entity.URL{
					{
						ID:          "42",
						ShortURL:    fmt.Sprintf("%s42", config.BaseURL()),
						OriginalURL: "https://ya.ru/",
					},
				},
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.storageCfg()
			resp, err := s.attr.service.GetURLsByUserID(context.Background(), tt.userID)

			assert.Equal(t, tt.want.resp, resp)
			assert.Equal(t, tt.want.err, errors.UnwrapAll(err))
		})
	}
}
