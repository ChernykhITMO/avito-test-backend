package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func Migrate(ctx context.Context, dsn string, dir string) error {
	const op = "internal.platform.postgres.Migrate"

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("%s: sql open: %w", op, err)
	}
	defer func() {
		_ = db.Close()
	}()

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("%s: set dialect: %w", op, err)
	}

	if err := goose.UpContext(ctx, db, dir); err != nil {
		return fmt.Errorf("%s: goose up: %w", op, err)
	}

	return nil
}
