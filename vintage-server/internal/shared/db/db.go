package db

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type DBTX interface {
	sqlx.ExtContext
	sqlx.PreparerContext

	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)

	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	Rebind(query string) string
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}
