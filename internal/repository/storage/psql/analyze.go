package psql

import (
	"context"
	"fmt"

	"github.com/dsbasko/yandex-go-shortener/internal/entity"
)

// Stats returns the stats of the URL.
func (s *Storage) Stats(ctx context.Context) (entity.URLStats, error) {
	query := `
		SELECT
			COUNT(*) AS users,
			(SELECT COUNT(DISTINCT user_id) FROM urls) as users
		FROM urls`

	row := s.conn.QueryRowContext(ctx, query)

	var resp entity.URLStats
	if err := row.Scan(&resp.URLs, &resp.Users); err != nil {
		return resp, fmt.Errorf("failed to scan response: %w", err)
	}

	return resp, nil
}
