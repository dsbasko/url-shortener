package psql

import (
	"context"
	"fmt"

	"github.com/dsbasko/yandex-go-shortener/internal/entities"
)

func (s *Storage) DeleteURLs(ctx context.Context, dto []entities.URL) (resp []entities.URL, err error) {
	var userID string
	var shortURLs []string

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
		return []entities.URL{}, fmt.Errorf("failed to execute query: %w", err)
	}

	if rows.Err() != nil {
		return []entities.URL{}, fmt.Errorf("failed to execute query: %w", rows.Err())
	}

	for rows.Next() {
		var url entities.URL
		if err = rows.StructScan(&url); err != nil {
			return []entities.URL{}, fmt.Errorf("failed to scan response: %w", err)
		}
		resp = append(resp, url)
	}

	return resp, nil
}
