package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/config"
	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/postgres"
	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/server"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	logger           *slog.Logger
	server           *http.Server
	db               *pgxpool.Pool
	slotRefillWorker workerRunner
	shutdownTimeout  int
}

type workerRunner interface {
	Run(ctx context.Context)
}

func New(ctx context.Context, cfg config.Config, logger *slog.Logger) (*App, error) {
	const op = "internal.app.New"

	if err := postgres.Migrate(ctx, cfg.Database.URL, "migrations"); err != nil {
		return nil, fmt.Errorf("%s: migrate: %w", op, err)
	}

	db, err := postgres.Open(ctx, postgres.Config{
		URL:             cfg.Database.URL,
		MaxConns:        cfg.Database.MaxConns,
		MinConns:        cfg.Database.MinConns,
		MaxConnIdleTime: time.Duration(cfg.Database.MaxConnIdleTimeSeconds) * time.Second,
		MaxConnLifetime: time.Duration(cfg.Database.MaxConnLifeTimeSeconds) * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: open db: %w", op, err)
	}

	modules := buildModules(cfg, logger, db)

	return &App{
		logger:           logger,
		server:           server.New(cfg.HTTP, modules.handler),
		db:               db,
		slotRefillWorker: modules.slotRefillWorker,
		shutdownTimeout:  cfg.HTTP.ShutdownTimeoutSeconds,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	const op = "internal.app.App.Run"

	errCh := make(chan error, 1)

	go a.slotRefillWorker.Run(ctx)

	go func() {
		a.logger.Info("http server started", slog.String("addr", a.server.Addr))

		err := a.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("%s: listen and serve: %w", op, err)
			return
		}

		close(errCh)
	}()

	select {
	case <-ctx.Done():
		a.logger.Info("shutdown signal received")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			time.Duration(a.shutdownTimeout)*time.Second,
		)
		defer cancel()

		if err := a.server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("%s: shutdown server: %w", op, err)
		}

		a.db.Close()

		return nil
	case err := <-errCh:
		a.db.Close()
		return err
	}
}
