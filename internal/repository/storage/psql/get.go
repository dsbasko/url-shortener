package psql

import (
	"context"
	"fmt"

	"github.com/dsbasko/yandex-go-shortener/internal/entity"
)

// GetURLByOriginalURL returns a URL by original URL.
func (s *Storage) GetURLByOriginalURL(
	ctx context.Context,
	originalURL string,
) (resp entity.URL, err error) {
	row := s.conn.QueryRowxContext(
		ctx,
		`SELECT * FROM urls WHERE original_url = $1`,
		originalURL,
	)

	if err = row.StructScan(&resp); err != nil {
		return entity.URL{}, fmt.Errorf("failed to scan response: %w", err)
	}

	return resp, nil
}

// GetURLByShortURL returns a URL by short URL.
func (s *Storage) GetURLByShortURL(
	ctx context.Context,
	shortURL string,
) (resp entity.URL, err error) {
	row := s.conn.QueryRowxContext(
		ctx,
		`SELECT * FROM urls WHERE short_url = $1`,
		shortURL,
	)

	if err = row.StructScan(&resp); err != nil {
		return entity.URL{}, fmt.Errorf("failed to scan response: %w", err)
	}

	return resp, nil
}

// GetURLsByUserID returns URLs by user ID.
func (s *Storage) GetURLsByUserID(
	ctx context.Context,
	userID string,
) (resp []entity.URL, err error) {
	rows, err := s.conn.QueryxContext(
		ctx,
		`SELECT * FROM urls WHERE user_id = $1`,
		userID,
	)
	if err != nil {
		return []entity.URL{}, fmt.Errorf("failed to get rows: %w", err)
	}

	if rows.Err() != nil {
		return []entity.URL{}, fmt.Errorf("failed to get rows: %w", rows.Err())
	}

	for rows.Next() {
		var foundRow entity.URL

		if err = rows.StructScan(&foundRow); err != nil {
			return []entity.URL{}, fmt.Errorf("failed to scan response: %w", err)
		}

		resp = append(resp, foundRow)
	}

	return resp, nil
}
