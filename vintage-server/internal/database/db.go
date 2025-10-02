package database

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

func NewPostgres(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Tuning conn pool — sesuaikan dengan workload dan hardware
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(30 * time.Minute)

	// Ping dengan timeout untuk memastikan koneksi valid
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}