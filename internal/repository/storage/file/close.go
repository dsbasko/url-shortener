package file

import "fmt"

// Close closes the storage connection.
func (s *Storage) Close() error {
	if err := s.file.Close(); err != nil {
		return fmt.Errorf("failed to close the storage connection: %w", err)
	}
	return nil
}
