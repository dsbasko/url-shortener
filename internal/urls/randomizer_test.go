package urls

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_URLs_RandomString(t *testing.T) {
	tests := []struct {
		name    string
		length  int
		wantLen int
	}{
		{
			name:    "Empty",
			length:  0,
			wantLen: 0,
		},
		{
			name:    "Length 1",
			length:  1,
			wantLen: 1,
		},
		{
			name:    "Length 10",
			length:  10,
			wantLen: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantLen, len(RandomString(tt.length)))
		})
	}
}
