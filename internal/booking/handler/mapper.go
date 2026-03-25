package handler

import (
	"time"

	bookingmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/booking/model"
)

func toBookingResponse(booking bookingmodel.Booking) bookingResponse {
	return bookingResponse{
		ID:             booking.ID.String(),
		SlotID:         booking.SlotID.String(),
		UserID:         booking.UserID.String(),
		Status:         string(booking.Status),
		ConferenceLink: booking.ConferenceLink,
		CreatedAt:      booking.CreatedAt.UTC().Format(time.RFC3339),
	}
}
