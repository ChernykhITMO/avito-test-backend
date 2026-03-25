package model

import (
	"testing"
	"time"
)

func TestScheduleHelpers(t *testing.T) {
	days, err := NormalizeDaysOfWeek([]int{5, 1, 3, 3})
	if err != nil {
		t.Fatalf("NormalizeDaysOfWeek: %v", err)
	}

	if len(days) != 3 {
		t.Fatalf("unexpected days: %#v", days)
	}
	if days[0] != 1 || days[1] != 3 || days[2] != 5 {
		t.Fatalf("unexpected normalized days: %#v", days)
	}

	minutes, err := ParseClock("09:30")
	if err != nil {
		t.Fatalf("ParseClock: %v", err)
	}
	if minutes != 570 {
		t.Fatalf("unexpected minutes: %d", minutes)
	}

	if !IncludesDate(days, time.Date(2026, 3, 23, 0, 0, 0, 0, time.UTC)) {
		t.Fatal("expected monday to be included")
	}
}
