package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/app"
	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/config"
	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/logger"
)

func run() int {
	const op = "cmd.api.run"

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		slog.New(slog.NewJSONHandler(os.Stdout, nil)).Error("failed to load config", slog.Any("error", err), slog.String("op", op))
		return 1
	}

	log := logger.New(cfg.LogLevel)

	application, err := app.New(ctx, cfg, log)
	if err != nil {
		log.Error("failed to initialize application", slog.Any("error", err), slog.String("op", op))
		return 1
	}

	if err := application.Run(ctx); err != nil {
		log.Error("application stopped with error", slog.Any("error", err), slog.String("op", op))
		return 1
	}

	return 0
}
