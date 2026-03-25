package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	authmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/auth/model"
	jwtplatform "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/jwt"
	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/middleware"
	roommodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/room/model"
	roomservice "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/room/service"
	"github.com/google/uuid"
)

type roomRepoTestStub struct {
	rooms []roommodel.Room
}

func (s *roomRepoTestStub) Create(ctx context.Context, room roommodel.Room) error {
	s.rooms = append(s.rooms, room)
	return nil
}

func (s *roomRepoTestStub) List(ctx context.Context) ([]roommodel.Room, error) {
	return s.rooms, nil
}

func (s *roomRepoTestStub) Exists(ctx context.Context, roomID uuid.UUID) (bool, error) {
	return true, nil
}

func TestCreateRoomHandler(t *testing.T) {
	repo := &roomRepoTestStub{}
	service := roomservice.New(repo)
	handler := New(service)
	jwtManager := jwtplatform.New("secret", time.Hour)
	token, err := jwtManager.Issue(authmodel.DummyAdminID, authmodel.RoleAdmin)
	if err != nil {
		t.Fatalf("Issue token: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/rooms/create", strings.NewReader(`{"name":"Blue","capacity":8}`))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	middleware.Authenticate(jwtManager)(handler.Create()).ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}
	if len(repo.rooms) != 1 {
		t.Fatalf("expected 1 room to be created, got %d", len(repo.rooms))
	}
}

func TestCreateRoomHandlerRejectsInvalidCapacity(t *testing.T) {
	repo := &roomRepoTestStub{}
	service := roomservice.New(repo)
	handler := New(service)
	jwtManager := jwtplatform.New("secret", time.Hour)
	token, err := jwtManager.Issue(authmodel.DummyAdminID, authmodel.RoleAdmin)
	if err != nil {
		t.Fatalf("Issue token: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/rooms/create", strings.NewReader(`{"name":"Blue","capacity":0}`))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	middleware.Authenticate(jwtManager)(handler.Create()).ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
	if len(repo.rooms) != 0 {
		t.Fatalf("expected no rooms to be created, got %d", len(repo.rooms))
	}
}

func TestCreateRoomHandlerRejectsEmptyName(t *testing.T) {
	repo := &roomRepoTestStub{}
	service := roomservice.New(repo)
	handler := New(service)
	jwtManager := jwtplatform.New("secret", time.Hour)
	token, err := jwtManager.Issue(authmodel.DummyAdminID, authmodel.RoleAdmin)
	if err != nil {
		t.Fatalf("Issue token: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/rooms/create", strings.NewReader(`{"name":"   ","capacity":8}`))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	middleware.Authenticate(jwtManager)(handler.Create()).ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
	if len(repo.rooms) != 0 {
		t.Fatalf("expected no rooms to be created, got %d", len(repo.rooms))
	}
}
