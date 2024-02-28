package psql

import (
	"context"
	"fmt"

	"github.com/dsbasko/yandex-go-shortener/internal/entity"
)

// Stats returns the stats of the URL.
func (s *Storage) Stats(ctx context.Context) (entity.URLStats, error) {
	row := s.conn.QueryRowContext(ctx, `SELECT COUNT(*) AS users, COUNT(DISTINCT user_id) FROM urls`)

	var resp entity.URLStats
	if err := row.Scan(&resp.URLs, &resp.Users); err != nil {
		return resp, fmt.Errorf("failed to scan response: %w", err)
	}

	return resp, nil
}
