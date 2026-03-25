package service

import (
	"context"
	"errors"
	"testing"

	authmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/auth/model"
	"github.com/google/uuid"
)

type authRepoStub struct {
	createFn     func(ctx context.Context, params authmodel.CreateUserParams) error
	getByEmailFn func(ctx context.Context, email string) (*authmodel.UserWithPassword, error)
}

func (s authRepoStub) Create(ctx context.Context, params authmodel.CreateUserParams) error {
	return s.createFn(ctx, params)
}

func (s authRepoStub) GetByEmail(ctx context.Context, email string) (*authmodel.UserWithPassword, error) {
	return s.getByEmailFn(ctx, email)
}

type tokenManagerStub struct {
	issueFn func(userID uuid.UUID, role authmodel.Role) (string, error)
}

func (s tokenManagerStub) Issue(userID uuid.UUID, role authmodel.Role) (string, error) {
	return s.issueFn(userID, role)
}

type passwordManagerStub struct {
	hashFn    func(password string) (string, error)
	compareFn func(hash, password string) error
}

func (s passwordManagerStub) Hash(password string) (string, error) {
	return s.hashFn(password)
}

func (s passwordManagerStub) Compare(hash, password string) error {
	return s.compareFn(hash, password)
}

type noopTransactor struct{}

func (noopTransactor) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func TestDummyLogin(t *testing.T) {
	service := New(
		authRepoStub{},
		tokenManagerStub{issueFn: func(userID uuid.UUID, role authmodel.Role) (string, error) {
			if userID != authmodel.DummyAdminID {
				t.Fatalf("unexpected user id: %s", userID)
			}
			if role != authmodel.RoleAdmin {
				t.Fatalf("unexpected role: %s", role)
			}

			return "token", nil
		}},
		passwordManagerStub{},
		noopTransactor{},
	)

	token, err := service.DummyLogin(context.Background(), "admin")
	if err != nil {
		t.Fatalf("DummyLogin returned error: %v", err)
	}
	if token != "token" {
		t.Fatalf("unexpected token: %s", token)
	}
}

func TestRegisterCreatesUser(t *testing.T) {
	var created bool

	service := New(
		authRepoStub{
			createFn: func(ctx context.Context, params authmodel.CreateUserParams) error {
				created = true
				if params.User.Email != "user@example.com" {
					t.Fatalf("unexpected email: %s", params.User.Email)
				}
				if params.PasswordHash != "hashed" {
					t.Fatalf("unexpected password hash: %s", params.PasswordHash)
				}

				return nil
			},
			getByEmailFn: func(ctx context.Context, email string) (*authmodel.UserWithPassword, error) {
				return nil, nil
			},
		},
		tokenManagerStub{},
		passwordManagerStub{
			hashFn: func(password string) (string, error) {
				if password != "secret" {
					t.Fatalf("unexpected password: %s", password)
				}

				return "hashed", nil
			},
			compareFn: func(hash, password string) error { return nil },
		},
		noopTransactor{},
	)

	user, err := service.Register(context.Background(), RegisterInput{
		Email:    "User@example.com",
		Password: "secret",
		Role:     "user",
	})
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}
	if !created {
		t.Fatal("expected user to be created")
	}
	if user.Role != authmodel.RoleUser {
		t.Fatalf("unexpected role: %s", user.Role)
	}
}

func TestLoginInvalidCredentials(t *testing.T) {
	service := New(
		authRepoStub{
			createFn: func(ctx context.Context, params authmodel.CreateUserParams) error { return nil },
			getByEmailFn: func(ctx context.Context, email string) (*authmodel.UserWithPassword, error) {
				return &authmodel.UserWithPassword{
					User: authmodel.User{
						ID:    uuid.New(),
						Email: email,
						Role:  authmodel.RoleUser,
					},
					PasswordHash: "hash",
				}, nil
			},
		},
		tokenManagerStub{issueFn: func(userID uuid.UUID, role authmodel.Role) (string, error) { return "token", nil }},
		passwordManagerStub{
			hashFn: func(password string) (string, error) { return "hash", nil },
			compareFn: func(hash, password string) error {
				return errors.New("mismatch")
			},
		},
		noopTransactor{},
	)

	_, err := service.Login(context.Background(), LoginInput{
		Email:    "user@example.com",
		Password: "bad",
	})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected invalid credentials, got %v", err)
	}
}
