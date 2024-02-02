package urls

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/dsbasko/yandex-go-shortener/internal/config"
	"github.com/dsbasko/yandex-go-shortener/internal/entities"
	"github.com/dsbasko/yandex-go-shortener/internal/jwt"
	"github.com/dsbasko/yandex-go-shortener/pkg/errors"
)

func (s *SuiteURLs) Test_CreateURL() {
	t := s.T()

	type args struct {
		ctx         context.Context
		originalURL string
	}
	type want struct {
		resp entities.URL
		uniq bool
		err  error
	}

	tests := []struct {
		name     string
		args     args
		storeCfg func()
		want     want
	}{
		{
			name: "Invalid URL",
			args: args{
				ctx:         context.Background(),
				originalURL: "invalid-url",
			},
			storeCfg: func() {},
			want: want{
				resp: entities.URL{},
				err:  ErrInvalidURL,
			},
		},
		{
			name: "Storage Error",
			args: args{
				ctx:         s.attr.ctxWithToken,
				originalURL: "https://ya.ru/",
			},
			storeCfg: func() {
				s.attr.store.EXPECT().
					CreateURL(gomock.Any(), gomock.Any()).
					Return(entities.URL{}, false, s.attr.errStore)
			},
			want: want{
				resp: entities.URL{},
				err:  s.attr.errStore,
			},
		},
		{
			name: "Success Unique",
			args: args{
				ctx:         s.attr.ctxWithToken,
				originalURL: "https://ya.ru/",
			},
			storeCfg: func() {
				s.attr.store.EXPECT().
					CreateURL(gomock.Any(), gomock.Any()).
					Return(entities.URL{
						ID:          "42",
						ShortURL:    "42",
						OriginalURL: "https://ya.ru/",
					}, true, nil)
			},
			want: want{
				resp: entities.URL{
					ID:          "42",
					ShortURL:    fmt.Sprintf("%s42", config.GetBaseURL()),
					OriginalURL: "https://ya.ru/",
				},
				uniq: true,
				err:  nil,
			},
		},
		{
			name: "Success NotUnique",
			args: args{
				ctx:         s.attr.ctxWithToken,
				originalURL: "https://ya.ru/",
			},
			storeCfg: func() {
				s.attr.store.EXPECT().
					CreateURL(gomock.Any(), gomock.Any()).
					Return(entities.URL{
						ID:          "42",
						ShortURL:    "42",
						OriginalURL: "https://ya.ru/",
					}, false, nil)
			},
			want: want{
				resp: entities.URL{
					ID:          "42",
					ShortURL:    fmt.Sprintf("%s42", config.GetBaseURL()),
					OriginalURL: "https://ya.ru/",
				},
				uniq: false,
				err:  nil,
			},
		},
		{
			name: "Empty Token",
			args: args{
				ctx:         context.Background(),
				originalURL: "https://ya.ru/",
			},
			storeCfg: func() {
				s.attr.store.EXPECT().
					CreateURL(gomock.Any(), gomock.Any()).
					Return(entities.URL{}, false, s.attr.errStore)
			},
			want: want{
				resp: entities.URL{},
				err:  jwt.ErrNotFoundFromContext,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.storeCfg()
			resp, uniq, err := s.attr.service.CreateURL(tt.args.ctx, tt.args.originalURL)

			assert.Equal(t, tt.want.resp, resp)
			assert.Equal(t, tt.want.uniq, uniq)
			assert.Equal(t, tt.want.err, errors.UnwrapAll(err))
		})
	}
}
