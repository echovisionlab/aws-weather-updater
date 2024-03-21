package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Database interface {
	BindNamed(query string, arg interface{}) (string, []interface{}, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Rebind(query string) string
	Close() error
	Query(query string, args ...any) (*sql.Rows, error)
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

func NewDatabase() (Database, error) {
	if cfg, err := readDatabaseConfig(); err != nil {
		return nil, fmt.Errorf("failed to retrieve database connection: %w", err)
	} else if db, err := sqlx.Connect("postgres", cfg.ConnStr()); err != nil {
		return nil, fmt.Errorf("failed to retrieve database connection: %w", err)
	} else {
		return db, nil
	}
}
