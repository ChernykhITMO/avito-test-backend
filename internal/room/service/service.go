package service

type Service struct {
	roomRepository RoomRepository
}

func New(roomRepository RoomRepository) *Service {
	return &Service{roomRepository: roomRepository}
}
