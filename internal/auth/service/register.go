package service

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	authmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/auth/model"
	"github.com/google/uuid"
)

func (s *Service) Register(ctx context.Context, input RegisterInput) (*authmodel.User, error) {
	const op = "internal.auth.service.Service.Register"

	parsedRole, ok := authmodel.ParseRole(input.Role)
	if !ok {
		return nil, ErrInvalidRole
	}

	email := strings.TrimSpace(strings.ToLower(input.Email))
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, ErrInvalidCredentials
	}

	if strings.TrimSpace(input.Password) == "" {
		return nil, ErrInvalidCredentials
	}

	passwordHash, err := s.passwords.Hash(input.Password)
	if err != nil {
		return nil, fmt.Errorf("%s: hash password: %w", op, err)
	}

	user := authmodel.User{
		ID:    uuid.New(),
		Email: email,
		Role:  parsedRole,
	}
	createdAt := time.Now().UTC()
	user.CreatedAt = createdAt

	err = s.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := s.userRepository.Create(txCtx, authmodel.CreateUserParams{
			User:         user,
			PasswordHash: passwordHash,
		}); err != nil {
			return fmt.Errorf("%s: create user: %w", op, err)
		}

		return nil
	})
	if err != nil {
		if errors.Is(err, authmodel.ErrEmailAlreadyExists) {
			return nil, ErrEmailAlreadyExists
		}

		return nil, fmt.Errorf("%s: transaction: %w", op, err)
	}

	return &user, nil
}
