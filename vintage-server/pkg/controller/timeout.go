package controller

import (
	"context"
	"time"
)

const (
	DefaultTransactionTimeout = 30 * time.Second
	DefaultQueryTimeout       = 3 * time.Second
)

func WithTxTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, DefaultTransactionTimeout)
}

func WithQueryTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, DefaultQueryTimeout)
}
