package worker

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

const refillBatchSize = 128

type scheduleRepository interface {
	ListRoomIDsNeedingGeneration(ctx context.Context, generatedUntilBefore time.Time, limit int) ([]uuid.UUID, error)
}

type slotGenerator interface {
	EnsureRange(ctx context.Context, roomID uuid.UUID, toDate time.Time) error
}

type clock interface {
	Now() time.Time
}

type SlotRefillWorker struct {
	logger               *slog.Logger
	scheduleRepository   scheduleRepository
	slotGenerator        slotGenerator
	clock                clock
	generationWindowDays int
	interval             time.Duration
}

func NewSlotRefillWorker(
	logger *slog.Logger,
	scheduleRepository scheduleRepository,
	slotGenerator slotGenerator,
	clock clock,
	generationWindowDays int,
	interval time.Duration,
) *SlotRefillWorker {
	return &SlotRefillWorker{
		logger:               logger,
		scheduleRepository:   scheduleRepository,
		slotGenerator:        slotGenerator,
		clock:                clock,
		generationWindowDays: generationWindowDays,
		interval:             interval,
	}
}

func (w *SlotRefillWorker) Run(ctx context.Context) {
	const op = "internal.worker.SlotRefillWorker.Run"

	if w == nil || w.interval <= 0 {
		if w != nil {
			w.logger.Info("slot refill worker skipped", slog.String("op", op), slog.Duration("interval", w.interval))
		}
		return
	}

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	w.runOnce(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.runOnce(ctx)
		}
	}
}

func (w *SlotRefillWorker) runOnce(ctx context.Context) {
	const op = "internal.worker.SlotRefillWorker.runOnce"

	targetDate := normalizeDate(w.clock.Now()).AddDate(0, 0, w.generationWindowDays)

	for {
		roomIDs, err := w.scheduleRepository.ListRoomIDsNeedingGeneration(ctx, targetDate, refillBatchSize)
		if err != nil {
			w.logger.Error("slot refill query failed", slog.String("op", op), slog.Any("error", fmt.Errorf("%s: list rooms needing generation: %w", op, err)))
			return
		}
		if len(roomIDs) == 0 {
			return
		}

		for _, roomID := range roomIDs {
			if err := w.slotGenerator.EnsureRange(ctx, roomID, targetDate); err != nil {
				w.logger.Error(
					"slot refill failed",
					slog.String("op", op),
					slog.String("room_id", roomID.String()),
					slog.Any("error", fmt.Errorf("%s: ensure range: %w", op, err)),
				)
			}
		}

		if len(roomIDs) < refillBatchSize {
			return
		}
	}
}

func normalizeDate(value time.Time) time.Time {
	utc := value.UTC()
	return time.Date(utc.Year(), utc.Month(), utc.Day(), 0, 0, 0, 0, time.UTC)
}
