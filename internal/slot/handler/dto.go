package handler

type slotResponse struct {
	ID     string `json:"id"`
	RoomID string `json:"roomId"`
	Start  string `json:"start"`
	End    string `json:"end"`
}

type listSlotsResponse struct {
	Slots []slotResponse `json:"slots"`
}
