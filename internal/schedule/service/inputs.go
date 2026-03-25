package service

import "github.com/google/uuid"

type CreateInput struct {
	RoomID     uuid.UUID
	DaysOfWeek []int
	StartTime  string
	EndTime    string
}
