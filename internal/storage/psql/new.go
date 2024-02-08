package psql

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/dsbasko/yandex-go-shortener/internal/config"
	"github.com/dsbasko/yandex-go-shortener/internal/interfaces"
	"github.com/dsbasko/yandex-go-shortener/pkg/logger"
)

// Storage a postgresql storage.
type Storage struct {
	log  *logger.Logger
	conn *sqlx.DB
}

// Ensure that Storage implements the Storage interface.
var _ interfaces.Storage = (*Storage)(nil)

// New creates a new instance of the postgresql storage.
func New(ctx context.Context, log *logger.Logger) (interfaces.Storage, error) {
	conn, err := sqlx.ConnectContext(ctx, "pgx", config.GetPsqlDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	conn.SetMaxOpenConns(config.GetPsqlMaxConns())

	if err = createTable(ctx, conn); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	log.Infof("postgresql storage initialized")

	return &Storage{
		log:  log,
		conn: conn,
	}, nil
}

// createTable creates the table if it does not exist and creates the necessary indexes.
func createTable(ctx context.Context, conn *sqlx.DB) error {
	if _, err := conn.ExecContext(ctx, `
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
		return fmt.Errorf("the database query could not be executed: %w", err)
	}

	return nil
}
