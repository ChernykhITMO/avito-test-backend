package service

type Service struct {
	roomRepository     RoomRepository
	scheduleRepository ScheduleRepository
	slotGenerator      SlotGenerator
	transactor         Transactor
	clock              Clock
	windowDays         int
}

func New(
	roomRepository RoomRepository,
	scheduleRepository ScheduleRepository,
	slotGenerator SlotGenerator,
	transactor Transactor,
	clock Clock,
	windowDays int,
) *Service {
	return &Service{
		roomRepository:     roomRepository,
		scheduleRepository: scheduleRepository,
		slotGenerator:      slotGenerator,
		transactor:         transactor,
		clock:              clock,
		windowDays:         windowDays,
	}
}
