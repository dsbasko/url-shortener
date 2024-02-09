package handlers

import (
	"errors"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/dsbasko/yandex-go-shortener/pkg/test"
)

func (s *SuiteHandlers) Test_Ping() {
	t := s.T()

	tests := []struct {
		name           string
		storeCfg       func()
		wantStatusCode int
		wantBody       func() string
	}{
		{
			name: "Error",
			storeCfg: func() {
				s.attr.pinger.EXPECT().Ping(gomock.Any()).Return(errors.New(""))
			},
			wantStatusCode: http.StatusBadRequest,
			wantBody:       func() string { return "" },
		},
		{
			name: "Success",
			storeCfg: func() {
				s.attr.pinger.EXPECT().Ping(gomock.Any()).Return(nil)
			},
			wantStatusCode: http.StatusOK,
			wantBody:       func() string { return "pong" },
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.storeCfg()
			resp, body := test.Request(t, s.attr.ts, &test.RequestArgs{
				Method: "GET",
				Path:   "/ping",
			})
			err := resp.Body.Close()
			assert.NoError(t, err)

			assert.Equal(t, tt.wantStatusCode, resp.StatusCode)
			assert.Equal(t, tt.wantBody(), body)
		})
	}
}
