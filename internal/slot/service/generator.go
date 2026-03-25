package service

import (
	"time"

	schedulemodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/schedule/model"
	slotmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/slot/model"
	"github.com/google/uuid"
)

func generateSlots(schedule schedulemodel.Schedule, fromDate, toDate time.Time) []slotmodel.Slot {
	slots := make([]slotmodel.Slot, 0)

	for current := normalizeSlotDate(fromDate); !current.After(normalizeSlotDate(toDate)); current = current.AddDate(0, 0, 1) {
		if !schedulemodel.IncludesDate(schedule.DaysOfWeek, current) {
			continue
		}

		for minute := schedule.StartMinute; minute < schedule.EndMinute; minute += 30 {
			startAt := current.Add(time.Duration(minute) * time.Minute)
			endAt := startAt.Add(30 * time.Minute)

			slots = append(slots, slotmodel.Slot{
				ID:         uuid.New(),
				RoomID:     schedule.RoomID,
				ScheduleID: schedule.ID,
				SlotDate:   current,
				StartAt:    startAt,
				EndAt:      endAt,
				CreatedAt:  time.Now().UTC(),
			})
		}
	}

	return slots
}

func normalizeSlotDate(value time.Time) time.Time {
	return time.Date(value.UTC().Year(), value.UTC().Month(), value.UTC().Day(), 0, 0, 0, 0, time.UTC)
}
