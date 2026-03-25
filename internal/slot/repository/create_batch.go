package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/postgres"
	slotmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/slot/model"
)

func (r *Repository) CreateBatch(ctx context.Context, slots []slotmodel.Slot) error {
	const op = "internal.slot.repository.Repository.CreateBatch"

	if len(slots) == 0 {
		return nil
	}

	querier := postgres.QuerierFromContext(ctx, r.db)

	var builder strings.Builder
	args := make([]any, 0, len(slots)*7)

	builder.WriteString(`
		INSERT INTO slots (id, room_id, schedule_id, slot_date, start_at, end_at, created_at)
		VALUES
	`)

	for index, slot := range slots {
		base := index * 7
		builder.WriteString(fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)", base+1, base+2, base+3, base+4, base+5, base+6, base+7))
		if index < len(slots)-1 {
			builder.WriteString(",")
		}

		args = append(args, slot.ID, slot.RoomID, slot.ScheduleID, slot.SlotDate, slot.StartAt, slot.EndAt, slot.CreatedAt)
	}

	builder.WriteString(` ON CONFLICT (room_id, start_at) DO NOTHING`)

	_, err := querier.Exec(ctx, builder.String(), args...)
	if err != nil {
		return fmt.Errorf("%s: exec insert slots batch: %w", op, err)
	}

	return nil
}
