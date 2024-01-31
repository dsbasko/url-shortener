package urls

import (
	"context"
	"strings"
	"testing"

	"github.com/dsbasko/yandex-go-shortener/internal/config"
	"github.com/dsbasko/yandex-go-shortener/internal/entities"
	"github.com/dsbasko/yandex-go-shortener/internal/jwt"
	"github.com/dsbasko/yandex-go-shortener/pkg/api"
	"github.com/dsbasko/yandex-go-shortener/pkg/errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func (s *SuiteURLs) Test_CreateURLs() {
	t := s.T()

	type args struct {
		ctx context.Context
		dto []api.CreateURLsRequest
	}
	type want struct {
		resp []api.CreateURLsResponse
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
			name: "Empty Token",
			args: args{
				ctx: context.Background(),
				dto: []api.CreateURLsRequest{
					{OriginalURL: "invalid-url", CorrelationID: "1"},
				},
			},
			storeCfg: func() {},
			want: want{
				resp: []api.CreateURLsResponse{},
				err:  jwt.ErrNotFoundFromContext,
			},
		},
		{
			name: "Invalid URL",
			args: args{
				ctx: s.attr.ctxWithToken,
				dto: []api.CreateURLsRequest{
					{OriginalURL: "invalid-url", CorrelationID: "1"},
				},
			},
			storeCfg: func() {},
			want: want{
				resp: []api.CreateURLsResponse{},
				err:  ErrInvalidURL,
			},
		},
		{
			name: "Storage Error",
			args: args{
				ctx: s.attr.ctxWithToken,
				dto: []api.CreateURLsRequest{
					{OriginalURL: "https://ya.ru/", CorrelationID: "1"},
				},
			},
			storeCfg: func() {
				s.attr.store.EXPECT().
					CreateURLs(gomock.Any(), gomock.Any()).
					Return(nil, s.attr.errStore)
			},
			want: want{
				resp: []api.CreateURLsResponse{},
				err:  s.attr.errStore,
			},
		},
		{
			name: "Success",
			args: args{
				ctx: s.attr.ctxWithToken,
				dto: []api.CreateURLsRequest{
					{OriginalURL: "https://ya.ru/", CorrelationID: "1"},
				},
			},
			storeCfg: func() {
				s.attr.store.EXPECT().
					CreateURLs(gomock.Any(), gomock.Any()).
					Return([]entities.URL{
						{
							OriginalURL: "https://ya.ru/",
							ShortURL:    "42",
						},
					}, nil)
			},
			want: want{
				resp: []api.CreateURLsResponse{
					{
						ShortURL:      "42",
						CorrelationID: "1",
					},
				},
				uniq: true,
				err:  nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.storeCfg()
			resp, err := s.attr.service.CreateURLs(tt.args.ctx, tt.args.dto)

			assert.Equal(t, tt.want.err, errors.UnwrapAll(err))
			assert.Equal(t, len(tt.want.resp), len(resp))
			if len(resp) > 0 {
				assert.True(t, strings.Contains(resp[0].ShortURL, config.GetBaseURL()))
				assert.Equal(t, tt.want.resp[0].CorrelationID, resp[0].CorrelationID)
			}
		})
	}
}
