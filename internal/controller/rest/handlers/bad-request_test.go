package handlers

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dsbasko/yandex-go-shortener/pkg/test"
)

func (s *SuiteHandlers) Test_BadRequest() {
	t := s.T()

	tests := []struct {
		name string
		path string
	}{
		{
			name: "Success",
			path: "/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, _ := test.Request(t, s.attr.ts, &test.RequestArgs{
				Method: "POST",
				Path:   "/ping",
			})
			err := resp.Body.Close()
			assert.NoError(t, err)

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}
}
