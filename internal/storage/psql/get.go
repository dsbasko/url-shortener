package psql

import (
	"context"
	"fmt"

	"github.com/dsbasko/yandex-go-shortener/internal/entities"
)

func (s *Storage) GetURLByOriginalURL(
	ctx context.Context,
	originalURL string,
) (resp entities.URL, err error) {
	row := s.conn.QueryRowxContext(
		ctx,
		`SELECT * FROM urls WHERE original_url = $1`,
		originalURL,
	)

	if err = row.StructScan(&resp); err != nil {
		return entities.URL{}, fmt.Errorf("failed to scan response: %w", err)
	}

	return resp, nil
}

func (s *Storage) GetURLByShortURL(
	ctx context.Context,
	shortURL string,
) (resp entities.URL, err error) {
	row := s.conn.QueryRowxContext(
		ctx,
		`SELECT * FROM urls WHERE short_url = $1`,
		shortURL,
	)

	if err = row.StructScan(&resp); err != nil {
		return entities.URL{}, fmt.Errorf("failed to scan response: %w", err)
	}

	return resp, nil
}

func (s *Storage) GetURLsByUserID(
	ctx context.Context,
	userID string,
) (resp []entities.URL, err error) {
	rows, err := s.conn.QueryxContext(
		ctx,
		`SELECT * FROM urls WHERE user_id = $1`,
		userID,
	)
	if err != nil {
		return []entities.URL{}, fmt.Errorf("failed to get rows: %w", err)
	}

	if rows.Err() != nil {
		return []entities.URL{}, fmt.Errorf("failed to get rows: %w", rows.Err())
	}

	for rows.Next() {
		var foundRow entities.URL

		if err = rows.StructScan(&foundRow); err != nil {
			return []entities.URL{}, fmt.Errorf("failed to scan response: %w", err)
		}

		resp = append(resp, foundRow)
	}

	return resp, nil
}
