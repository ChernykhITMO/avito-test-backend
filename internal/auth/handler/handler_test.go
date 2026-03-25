package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	authmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/auth/model"
	authservice "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/auth/service"
	"github.com/google/uuid"
)

type authRepoTestStub struct{}

func (authRepoTestStub) Create(ctx context.Context, params authmodel.CreateUserParams) error {
	return nil
}

func (authRepoTestStub) GetByEmail(ctx context.Context, email string) (*authmodel.UserWithPassword, error) {
	return nil, nil
}

type authTokenStub struct{}

func (authTokenStub) Issue(userID uuid.UUID, role authmodel.Role) (string, error) {
	return "token", nil
}

type authPasswordStub struct{}

func (authPasswordStub) Hash(password string) (string, error) { return "hash", nil }
func (authPasswordStub) Compare(hash, password string) error  { return nil }

type authTxStub struct{}

func (authTxStub) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func TestDummyLoginHandler(t *testing.T) {
	service := authservice.New(authRepoTestStub{}, authTokenStub{}, authPasswordStub{}, authTxStub{})
	handler := New(service)

	req := httptest.NewRequest(http.MethodPost, "/dummyLogin", strings.NewReader(`{"role":"admin"}`))
	rec := httptest.NewRecorder()

	handler.DummyLogin().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body["token"] != "token" {
		t.Fatalf("unexpected token: %s", body["token"])
	}
}

func TestRegisterHandlerReturnsResponseDTO(t *testing.T) {
	service := authservice.New(authRepoTestStub{}, authTokenStub{}, authPasswordStub{}, authTxStub{})
	handler := New(service)

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"email":"USER@example.com","password":"secret","role":"user"}`))
	rec := httptest.NewRecorder()

	handler.Register().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}

	var body struct {
		User struct {
			ID        string `json:"id"`
			Email     string `json:"email"`
			Role      string `json:"role"`
			CreatedAt string `json:"createdAt"`
		} `json:"user"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body.User.ID == "" {
		t.Fatal("expected user id in response")
	}
	if body.User.Email != "user@example.com" {
		t.Fatalf("unexpected email: %s", body.User.Email)
	}
	if body.User.Role != "user" {
		t.Fatalf("unexpected role: %s", body.User.Role)
	}
	if body.User.CreatedAt == "" {
		t.Fatal("expected createdAt in response")
	}
}

func TestDummyLoginHandlerRejectsEmptyRole(t *testing.T) {
	service := authservice.New(authRepoTestStub{}, authTokenStub{}, authPasswordStub{}, authTxStub{})
	handler := New(service)

	req := httptest.NewRequest(http.MethodPost, "/dummyLogin", strings.NewReader(`{"role":"   "}`))
	rec := httptest.NewRecorder()

	handler.DummyLogin().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestRegisterHandlerRejectsEmptyRequiredFields(t *testing.T) {
	service := authservice.New(authRepoTestStub{}, authTokenStub{}, authPasswordStub{}, authTxStub{})
	handler := New(service)

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"email":" ","password":"secret","role":"user"}`))
	rec := httptest.NewRecorder()

	handler.Register().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestLoginHandlerRejectsEmptyRequiredFields(t *testing.T) {
	service := authservice.New(authRepoTestStub{}, authTokenStub{}, authPasswordStub{}, authTxStub{})
	handler := New(service)

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"email":"user@example.com","password":"   "}`))
	rec := httptest.NewRecorder()

	handler.Login().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}
