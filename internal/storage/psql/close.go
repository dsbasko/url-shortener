package psql

import "fmt"

// Close closes the storage connection.
func (s *Storage) Close() error {
	if err := s.conn.Close(); err != nil {
		return fmt.Errorf("failed to close the storage connection: %w", err)
	}
	return nil
}
