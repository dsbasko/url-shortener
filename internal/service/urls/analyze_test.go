package urls

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/dsbasko/url-shortener/internal/entity"
	"github.com/dsbasko/url-shortener/pkg/errors"
)

func (s *SuiteURLs) Test_Stats() {
	t := s.T()

	type want struct {
		resp entity.URLStats
		err  error
	}

	tests := []struct {
		name       string
		storageCfg func()
		want       want
	}{
		{
			name: "Storage Error",
			storageCfg: func() {
				s.attr.urlAnalyzer.EXPECT().
					Stats(gomock.Any()).
					Return(entity.URLStats{}, s.attr.errStorage)
			},
			want: want{
				resp: entity.URLStats{},
				err:  s.attr.errStorage,
			},
		},

		{
			name: "Success",
			storageCfg: func() {
				s.attr.urlAnalyzer.EXPECT().
					Stats(gomock.Any()).
					Return(entity.URLStats{
						Users: "42",
						URLs:  "42",
					}, nil)
			},
			want: want{
				resp: entity.URLStats{
					Users: "42",
					URLs:  "42",
				},
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.storageCfg()
			resp, err := s.attr.service.Stats(context.Background())

			assert.Equal(t, tt.want.resp, resp)
			assert.Equal(t, tt.want.err, errors.UnwrapAll(err))
		})
	}
}
