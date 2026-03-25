package handler

type createRoomRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Capacity    *int    `json:"capacity"`
}

type roomResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Capacity    *int    `json:"capacity,omitempty"`
	CreatedAt   string  `json:"createdAt,omitempty"`
}

type createRoomResponse struct {
	Room roomResponse `json:"room"`
}

type listRoomsResponse struct {
	Rooms []roomResponse `json:"rooms"`
}
