package model

import "errors"

var (
	ErrBookingNotFound   = errors.New("booking not found")
	ErrSlotAlreadyBooked = errors.New("slot already booked")
	ErrForbiddenCancel   = errors.New("cannot cancel another user's booking")
	ErrSlotInPast        = errors.New("slot is in the past")
)
