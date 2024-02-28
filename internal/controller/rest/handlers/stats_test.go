package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/dsbasko/yandex-go-shortener/internal/config"
	"github.com/dsbasko/yandex-go-shortener/internal/entity"
	"github.com/dsbasko/yandex-go-shortener/pkg/test"
)

func (s *SuiteHandlers) Test_Stats() {
	t := s.T()

	tests := []struct {
		name     string
		cfg      func()
		wantCode int
		wantBody func() string
	}{
		{
			name: "Check Trusted Subnet Error",
			cfg: func() {
				config.SetTrustedSubnet("42")
			},
			wantCode: http.StatusBadRequest,
			wantBody: func() string { return "" },
		},
		{
			name: "IP Not Trusted",
			cfg: func() {
				config.SetTrustedSubnet("127.0.0.0/32")
			},
			wantCode: http.StatusForbidden,
			wantBody: func() string { return "" },
		},
		{
			name: "Service Error",
			cfg: func() {
				config.SetTrustedSubnet("127.0.0.0/24")
				s.attr.urlsAnalyzer.EXPECT().
					Stats(gomock.Any()).
					Return(entity.URLStats{}, s.attr.errService)
			},
			wantCode: http.StatusInternalServerError,
			wantBody: func() string { return "" },
		},
		{
			name: "Success",
			cfg: func() {
				config.SetTrustedSubnet("127.0.0.0/24")
				s.attr.urlsAnalyzer.EXPECT().
					Stats(gomock.Any()).
					Return(entity.URLStats{
						Users: "42",
						URLs:  "42",
					}, nil)
			},
			wantCode: http.StatusOK,
			wantBody: func() string {
				resp, _ := json.Marshal(entity.URLStats{
					Users: "42",
					URLs:  "42",
				})
				return string(resp) + "\n"
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cfg()

			resp, respBody := test.Request(t, s.attr.ts, &test.RequestArgs{
				Method: "GET",
				Path:   "/api/internal/stats",
			})
			err := resp.Body.Close()
			assert.NoError(t, err)

			assert.Equal(t, tt.wantCode, resp.StatusCode)
			assert.Equal(t, tt.wantBody(), respBody)
		})
	}
}
