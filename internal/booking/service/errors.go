package service

import "errors"

var (
	ErrSlotNotFound      = errors.New("slot not found")
	ErrSlotInPast        = errors.New("slot is in the past")
	ErrSlotAlreadyBooked = errors.New("slot already booked")
	ErrBookingNotFound   = errors.New("booking not found")
	ErrForbiddenCancel   = errors.New("cannot cancel another user's booking")
)
