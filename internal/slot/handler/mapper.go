package handler

import (
	"time"

	slotmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/slot/model"
)

func toSlotResponse(slot slotmodel.Slot) slotResponse {
	return slotResponse{
		ID:     slot.ID.String(),
		RoomID: slot.RoomID.String(),
		Start:  slot.StartAt.UTC().Format(time.RFC3339),
		End:    slot.EndAt.UTC().Format(time.RFC3339),
	}
}
