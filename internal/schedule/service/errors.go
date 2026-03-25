package service

import "errors"

var (
	ErrRoomNotFound          = errors.New("room not found")
	ErrInvalidSchedule       = errors.New("invalid schedule")
	ErrScheduleAlreadyExists = errors.New("schedule already exists")
)
