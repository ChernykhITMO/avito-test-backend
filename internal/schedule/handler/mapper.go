package handler

import schedulemodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/schedule/model"

func toScheduleResponse(schedule schedulemodel.Schedule) scheduleResponse {
	return scheduleResponse{
		ID:         schedule.ID.String(),
		RoomID:     schedule.RoomID.String(),
		DaysOfWeek: schedule.DaysOfWeek,
		StartTime:  schedulemodel.FormatClock(schedule.StartMinute),
		EndTime:    schedulemodel.FormatClock(schedule.EndMinute),
	}
}
