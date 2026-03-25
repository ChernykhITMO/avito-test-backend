package service

type Service struct {
	roomRepository     RoomRepository
	scheduleRepository ScheduleRepository
	slotRepository     SlotRepository
	transactor         Transactor
	clock              Clock
}

func New(
	roomRepository RoomRepository,
	scheduleRepository ScheduleRepository,
	slotRepository SlotRepository,
	transactor Transactor,
	clock Clock,
) *Service {
	return &Service{
		roomRepository:     roomRepository,
		scheduleRepository: scheduleRepository,
		slotRepository:     slotRepository,
		transactor:         transactor,
		clock:              clock,
	}
}
