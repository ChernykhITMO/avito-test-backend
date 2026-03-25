package service

type Service struct {
	userRepository UserRepository
	tokenManager   TokenManager
	passwords      PasswordManager
	transactor     Transactor
}

func New(
	userRepository UserRepository,
	tokenManager TokenManager,
	passwords PasswordManager,
	transactor Transactor,
) *Service {
	return &Service{
		userRepository: userRepository,
		tokenManager:   tokenManager,
		passwords:      passwords,
		transactor:     transactor,
	}
}
