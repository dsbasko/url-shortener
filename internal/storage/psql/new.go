package psql

import (
	"context"
	"fmt"

	"github.com/dsbasko/yandex-go-shortener/internal/config"
	"github.com/dsbasko/yandex-go-shortener/internal/interfaces"
	"github.com/dsbasko/yandex-go-shortener/pkg/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	log  *logger.Logger
	conn *sqlx.DB
}

var _ interfaces.Storage = (*Storage)(nil)

func New(ctx context.Context, log *logger.Logger) (interfaces.Storage, error) {
	conn, err := sqlx.ConnectContext(ctx, "pgx", config.GetPsqlDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	conn.SetMaxOpenConns(config.GetPsqlMaxConns())

	if _, err = conn.ExecContext(ctx, `
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

		CREATE TABLE IF NOT EXISTS urls (
		    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		    user_id TEXT NOT NULL,
			short_url VARCHAR(255) UNIQUE,
			original_url TEXT NOT NULL UNIQUE,
			is_deleted BOOLEAN DEFAULT FALSE NOT NULL
		);
		
		CREATE INDEX IF NOT EXISTS short_url ON urls (short_url);
	
		CREATE INDEX IF NOT EXISTS original_url ON urls (original_url);
	`); err != nil {
		return nil, fmt.Errorf("the database query could not be executed: %w", err)
	}

	log.Infof("postgresql storage initialized")

	return &Storage{
		log:  log,
		conn: conn,
	}, nil
}
