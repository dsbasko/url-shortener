package psql

import (
	"context"
	"errors"
	"fmt"

	"github.com/dsbasko/yandex-go-shortener/internal/entities"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Storage) CreateURL(
	ctx context.Context,
	dto entities.URL,
) (resp entities.URL, unique bool, err error) {
	row := s.conn.QueryRowxContext(
		ctx,
		`INSERT INTO urls (short_url, original_url, user_id) VALUES ($1, $2, $3) RETURNING *`,
		dto.ShortURL, dto.OriginalURL, dto.UserID,
	)

	if err = row.StructScan(&resp); err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == pgerrcode.UniqueViolation {
			foundDuplicate, errDuplicate := s.GetURLByOriginalURL(ctx, dto.OriginalURL)
			if errDuplicate != nil {
				return entities.URL{},
					false,
					fmt.Errorf("couldn't get the url from the original url: %w", errDuplicate)
			}
			return foundDuplicate, false, nil
		}
		return entities.URL{}, false, fmt.Errorf("failed to scan response: %w", err)
	}

	return resp, true, nil
}

func (s *Storage) CreateURLs(
	ctx context.Context,
	dto []entities.URL,
) (resp []entities.URL, err error) {
	rows, err := s.conn.NamedQueryContext(
		ctx,
		`INSERT INTO urls (short_url, original_url, user_id) VALUES (:short_url, :original_url, :user_id) RETURNING *`,
		dto,
	)
	if err != nil {
		return []entities.URL{}, fmt.Errorf("failed to insert urls: %w", err)
	}

	if rows.Err() != nil {
		return []entities.URL{}, fmt.Errorf("failed to insert urls: %w", rows.Err())
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
