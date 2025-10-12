package repository

import (
	"time"
	"vintage-server/internal/domain/account"
	"vintage-server/internal/shared/db"

	"github.com/jmoiron/sqlx"
)

type sqlAccountRepository struct {
	db db.DBTX
}

const DefaultQueryTimeout = 15 * time.Second

func (r *sqlAccountRepository) WithTx(tx *sqlx.Tx) account.AccountRepository {
	return &sqlAccountRepository{
		db: tx, // override db dengan tx agar bisa commit/rollback
	}
}

// --- Constructor ---
func NewAccountRepository(db *sqlx.DB) account.AccountRepository {
	return &sqlAccountRepository{db: db}
}
