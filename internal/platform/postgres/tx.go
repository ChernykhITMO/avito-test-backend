package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBTX interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, arguments ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, arguments ...any) pgx.Row
}

type Transactor interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type txManager struct {
	db *pgxpool.Pool
}

type txContextKey struct{}

func NewTransactor(db *pgxpool.Pool) Transactor {
	return &txManager{db: db}
}

func (t *txManager) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	const op = "internal.platform.postgres.txManager.WithinTransaction"

	if _, ok := txFromContext(ctx); ok {
		return fn(ctx)
	}

	tx, err := t.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return fmt.Errorf("%s: begin tx: %w", op, err)
	}

	txCtx := context.WithValue(ctx, txContextKey{}, tx)

	if err := fn(txCtx); err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed) {
			return fmt.Errorf("%s: %w", op, errors.Join(err, rollbackErr))
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: commit: %w", op, err)
	}

	return nil
}

func QuerierFromContext(ctx context.Context, fallback DBTX) DBTX {
	tx, ok := txFromContext(ctx)
	if ok {
		return tx
	}

	return fallback
}

func txFromContext(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txContextKey{}).(pgx.Tx)
	return tx, ok
}
