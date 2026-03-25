package service

import (
	"context"
	"fmt"

	authmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/auth/model"
	"github.com/google/uuid"
)

func (s *Service) DummyLogin(ctx context.Context, role string) (string, error) {
	const op = "internal.auth.service.Service.DummyLogin"

	_ = ctx

	parsedRole, ok := authmodel.ParseRole(role)
	if !ok {
		return "", ErrInvalidRole
	}

	var userID uuid.UUID
	switch parsedRole {
	case authmodel.RoleAdmin:
		userID = authmodel.DummyAdminID
	case authmodel.RoleUser:
		userID = authmodel.DummyUserID
	default:
		return "", ErrInvalidRole
	}

	token, err := s.tokenManager.Issue(userID, parsedRole)
	if err != nil {
		return "", fmt.Errorf("%s: issue token: %w", op, err)
	}

	return token, nil
}
