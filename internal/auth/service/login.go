package service

import (
	"context"
	"fmt"
	"net/mail"
	"strings"
)

func (s *Service) Login(ctx context.Context, input LoginInput) (string, error) {
	const op = "internal.auth.service.Service.Login"

	email := strings.TrimSpace(strings.ToLower(input.Email))
	if _, err := mail.ParseAddress(email); err != nil {
		return "", ErrInvalidCredentials
	}

	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("%s: get user by email: %w", op, err)
	}

	if user == nil {
		return "", ErrInvalidCredentials
	}

	if err := s.passwords.Compare(user.PasswordHash, input.Password); err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := s.tokenManager.Issue(user.ID, user.Role)
	if err != nil {
		return "", fmt.Errorf("%s: issue token: %w", op, err)
	}

	return token, nil
}
