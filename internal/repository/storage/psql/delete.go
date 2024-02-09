package psql

import (
	"context"
	"fmt"

	"github.com/dsbasko/yandex-go-shortener/internal/entity"
)

// DeleteURLs deletes URLs.
func (s *Storage) DeleteURLs(ctx context.Context, dto []entity.URL) (resp []entity.URL, err error) {
	var userID string
	shortURLs := make([]string, 0, len(dto))

	for _, url := range dto {
		if userID == "" {
			userID = url.UserID
		}
		shortURLs = append(shortURLs, url.ShortURL)
	}

	rows, err := s.conn.QueryxContext(
		ctx,
		`UPDATE urls SET is_deleted = true WHERE short_url = ANY($1) AND user_id = $2 RETURNING *`,
		shortURLs, userID,
	)
	if err != nil {
		return []entity.URL{}, fmt.Errorf("failed to execute query: %w", err)
	}

	if rows.Err() != nil {
		return []entity.URL{}, fmt.Errorf("failed to execute query: %w", rows.Err())
	}

	for rows.Next() {
		var url entity.URL
		if err = rows.StructScan(&url); err != nil {
			return []entity.URL{}, fmt.Errorf("failed to scan response: %w", err)
		}
		resp = append(resp, url)
	}

	return resp, nil
}
