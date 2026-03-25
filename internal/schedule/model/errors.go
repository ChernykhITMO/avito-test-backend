package model

import "errors"

var (
	ErrScheduleAlreadyExists = errors.New("schedule already exists")
	ErrInvalidSchedule       = errors.New("invalid schedule")
)
