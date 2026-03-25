package httpcommon

import "net/http"

const (
	CodeInvalidRequest    = "INVALID_REQUEST"
	CodeUnauthorized      = "UNAUTHORIZED"
	CodeNotFound          = "NOT_FOUND"
	CodeRoomNotFound      = "ROOM_NOT_FOUND"
	CodeSlotNotFound      = "SLOT_NOT_FOUND"
	CodeSlotAlreadyBooked = "SLOT_ALREADY_BOOKED"
	CodeBookingNotFound   = "BOOKING_NOT_FOUND"
	CodeForbidden         = "FORBIDDEN"
	CodeScheduleExists    = "SCHEDULE_EXISTS"
	CodeInternalError     = "INTERNAL_ERROR"
)

const (
	defaultInvalidRequestMessage = "invalid request"
	defaultInternalErrorMessage  = "internal server error"
	defaultUnauthorizedMessage   = "unauthorized"
	defaultForbiddenMessage      = "forbidden"
)

type ErrorPayload struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func WriteError(w http.ResponseWriter, status int, code, message string) {
	WriteJSON(w, status, ErrorPayload{
		Error: ErrorBody{
			Code:    code,
			Message: message,
		},
	})
}

func WriteInvalidRequest(w http.ResponseWriter) {
	WriteError(w, http.StatusBadRequest, CodeInvalidRequest, defaultInvalidRequestMessage)
}

func WriteInternalError(w http.ResponseWriter) {
	WriteError(w, http.StatusInternalServerError, CodeInternalError, defaultInternalErrorMessage)
}

func WriteUnauthorized(w http.ResponseWriter, message string) {
	if message == "" {
		message = defaultUnauthorizedMessage
	}

	WriteError(w, http.StatusUnauthorized, CodeUnauthorized, message)
}

func WriteForbidden(w http.ResponseWriter, message string) {
	if message == "" {
		message = defaultForbiddenMessage
	}

	WriteError(w, http.StatusForbidden, CodeForbidden, message)
}
