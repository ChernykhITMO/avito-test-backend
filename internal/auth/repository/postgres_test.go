package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	authmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/auth/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	pgxmock "github.com/pashagolub/pgxmock/v4"
)

func TestAuthRepositoryQueries(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("NewPool: %v", err)
	}
	defer mock.Close()

	repo := New(mock)
	user := authmodel.User{
		ID:        uuid.New(),
		Email:     "user@example.com",
		Role:      authmodel.RoleUser,
		CreatedAt: time.Now().UTC(),
	}
	params := authmodel.CreateUserParams{
		User:         user,
		PasswordHash: "hash",
	}

	mock.ExpectExec("INSERT INTO users").
		WithArgs(user.ID, user.Email, "hash", user.Role, user.CreatedAt).
		WillReturnError(&pgconn.PgError{Code: "23505"})

	err = repo.Create(context.Background(), params)
	if !errors.Is(err, authmodel.ErrEmailAlreadyExists) {
		t.Fatalf("expected email exists, got %v", err)
	}

	rows := pgxmock.NewRows([]string{"id", "email", "role", "created_at", "password_hash"}).
		AddRow(user.ID, user.Email, user.Role, user.CreatedAt, "hash")
	mock.ExpectQuery("SELECT id, email, role, created_at, COALESCE\\(password_hash, ''\\)").
		WithArgs(user.Email).
		WillReturnRows(rows)

	got, err := repo.GetByEmail(context.Background(), user.Email)
	if err != nil {
		t.Fatalf("GetByEmail: %v", err)
	}
	if got == nil || got.Email != user.Email {
		t.Fatalf("unexpected user: %#v", got)
	}
}
