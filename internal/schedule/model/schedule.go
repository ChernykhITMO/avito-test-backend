package model

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Schedule struct {
	ID             uuid.UUID
	RoomID         uuid.UUID
	DaysOfWeek     []int
	StartMinute    int
	EndMinute      int
	GeneratedUntil time.Time
	CreatedAt      time.Time
}

func NormalizeDaysOfWeek(days []int) ([]int, error) {
	if len(days) == 0 {
		return nil, ErrInvalidSchedule
	}

	unique := make(map[int]struct{}, len(days))
	normalized := make([]int, 0, len(days))

	for _, day := range days {
		if day < int(WeekdayMonday) || day > int(WeekdaySunday) {
			return nil, ErrInvalidSchedule
		}

		if _, ok := unique[day]; ok {
			continue
		}
		unique[day] = struct{}{}
		normalized = append(normalized, day)
	}

	sort.Ints(normalized)
	return normalized, nil
}

func ParseClock(value string) (int, error) {
	parts := strings.Split(value, ":")
	if len(parts) != 2 {
		return 0, ErrInvalidSchedule
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, ErrInvalidSchedule
	}

	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, ErrInvalidSchedule
	}

	if hours < 0 || hours > 23 || minutes < 0 || minutes > 59 {
		return 0, ErrInvalidSchedule
	}

	return hours*60 + minutes, nil
}

func FormatClock(totalMinutes int) string {
	return fmt.Sprintf("%02d:%02d", totalMinutes/60, totalMinutes%60)
}

func IncludesDate(days []int, date time.Time) bool {
	weekday := toScheduleWeekday(date.Weekday())
	for _, day := range days {
		if day == int(weekday) {
			return true
		}
	}

	return false
}

func toScheduleWeekday(weekday time.Weekday) Weekday {
	if weekday == time.Sunday {
		return WeekdaySunday
	}

	return Weekday(weekday)
}
