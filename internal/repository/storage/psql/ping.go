package psql

import (
	"context"
	"fmt"
)

// Ping checks the connection to the database.
func (s *Storage) Ping(ctx context.Context) error {
	if err := s.conn.PingContext(ctx); err != nil {
		return fmt.Errorf("the database query could not be executed: %w", err)
	}
	return nil
}
