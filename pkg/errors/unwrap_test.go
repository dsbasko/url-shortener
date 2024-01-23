package errors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnwrapAll(t *testing.T) {
	err := errors.New("test")

	tests := []struct {
		name    string
		err     error
		wantErr error
	}{
		{
			name:    "Nil",
			err:     nil,
			wantErr: nil,
		},
		{
			name:    "No Wrap",
			err:     err,
			wantErr: err,
		},
		{
			name:    "Wrap Once",
			err:     fmt.Errorf("first: %w", err),
			wantErr: err,
		},
		{
			name:    "Wrap Many",
			err:     fmt.Errorf("one: %w", fmt.Errorf("two: %w", err)),
			wantErr: err,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantErr, UnwrapAll(tt.err))
		})
	}
}
