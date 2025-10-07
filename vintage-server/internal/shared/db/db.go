package db

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type DBTX interface {
	sqlx.ExtContext
	sqlx.PreparerContext

	// NamedExecContext ada di *sqlx.DB dan *sqlx.Tx, jadi ini aman.
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)

	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	// Kita butuh metode-metode ini untuk pendekatan manual
	Rebind(query string) string
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
}
