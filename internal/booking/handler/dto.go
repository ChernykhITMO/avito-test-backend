package handler

type createBookingRequest struct {
	SlotID               string `json:"slotId"`
	CreateConferenceLink bool   `json:"createConferenceLink"`
}

type bookingResponse struct {
	ID             string  `json:"id"`
	SlotID         string  `json:"slotId"`
	UserID         string  `json:"userId"`
	Status         string  `json:"status"`
	ConferenceLink *string `json:"conferenceLink,omitempty"`
	CreatedAt      string  `json:"createdAt,omitempty"`
}

type paginationResponse struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}

type createBookingResponse struct {
	Booking bookingResponse `json:"booking"`
}

type cancelBookingResponse struct {
	Booking bookingResponse `json:"booking"`
}

type listBookingsResponse struct {
	Bookings []bookingResponse `json:"bookings"`
}

type listAllBookingsResponse struct {
	Bookings   []bookingResponse  `json:"bookings"`
	Pagination paginationResponse `json:"pagination"`
}
