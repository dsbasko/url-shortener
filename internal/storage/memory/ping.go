package memory

import (
	"context"
)

func (s *Storage) Ping(_ context.Context) error {
	return nil
}
