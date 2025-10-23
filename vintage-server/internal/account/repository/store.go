package repository

import (
	"context"
	"fmt"
	"vintage-server/internal/shared/db"

	"github.com/jmoiron/sqlx"
)

// AccountStore adalah entry point untuk semua repo terkait account
type AccountStore interface {
	// jalankan transaksi
	ExecTx(ctx context.Context, fn func(AccountStore) error) error

	// getter untuk akses repo
	GetWishlistRepo() WishlistRepository
	GetAccountRepo() AccountRepository
	GetAddressRepo() AddressRepository
}

// accountSqlStore implementasi AccountStore
type accountSqlStore struct {
	db db.DBTX

	WishlistRepo WishlistRepository
	AccountRepo  AccountRepository
	AddressRepo  AddressRepository
}

// constructor
func NewAccountStore(db db.DBTX) AccountStore {
	return &accountSqlStore{
		db:           db,
		WishlistRepo: NewWishlistRepositoryPostgres(db),
		AccountRepo:  NewAccountRepositoryPostgres(db),
		AddressRepo:  NewAddressRepositoryPostgres(db),
	}
}

// ExecTx menjalankan transaksi dan inject ke semua repo
func (s *accountSqlStore) ExecTx(ctx context.Context, fn func(AccountStore) error) error {
	sqlDB, ok := s.db.(*sqlx.DB)
	if !ok {
		return fmt.Errorf("db is not *sqlx.DB, cannot begin transaction")
	}

	txx, err := sqlDB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	// buat store sementara dengan semua repo memakai transaction
	txStore := &accountSqlStore{
		db:           txx,
		WishlistRepo: s.WishlistRepo.WithTx(txx),
		AccountRepo:  s.AccountRepo.WithTx(txx),
	}

	if err := fn(txStore); err != nil {
		_ = txx.Rollback()
		return err
	}

	return txx.Commit()
}

// getter WishlistRepo
func (s *accountSqlStore) GetWishlistRepo() WishlistRepository {
	return s.WishlistRepo
}

// getter AccountRepo
func (s *accountSqlStore) GetAccountRepo() AccountRepository {
	return s.AccountRepo
}

func (s *accountSqlStore) GetAddressRepo() AddressRepository {
	return s.AddressRepo
}
