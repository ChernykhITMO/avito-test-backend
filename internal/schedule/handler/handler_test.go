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
	schedulemodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/schedule/model"
	scheduleservice "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/schedule/service"
	"github.com/google/uuid"
)

type scheduleRoomRepoStub struct{}

func (scheduleRoomRepoStub) Exists(ctx context.Context, roomID uuid.UUID) (bool, error) {
	return true, nil
}

type scheduleRepoHandlerStub struct{}

func (scheduleRepoHandlerStub) Create(ctx context.Context, schedule schedulemodel.Schedule) error {
	return nil
}
func (scheduleRepoHandlerStub) GetByRoomID(ctx context.Context, roomID uuid.UUID) (*schedulemodel.Schedule, error) {
	return nil, nil
}
func (scheduleRepoHandlerStub) UpdateGeneratedUntil(ctx context.Context, roomID uuid.UUID, generatedUntil time.Time) error {
	return nil
}

type scheduleGeneratorStub struct{}

func (scheduleGeneratorStub) EnsureRange(ctx context.Context, roomID uuid.UUID, toDate time.Time) error {
	return nil
}

type scheduleClock struct{}

func (scheduleClock) Now() time.Time { return time.Date(2026, 3, 21, 10, 0, 0, 0, time.UTC) }

type scheduleTxStub struct{}

func (scheduleTxStub) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func TestCreateScheduleHandler(t *testing.T) {
	service := scheduleservice.New(scheduleRoomRepoStub{}, scheduleRepoHandlerStub{}, scheduleGeneratorStub{}, scheduleTxStub{}, scheduleClock{}, 30)
	handler := New(service)
	jwtManager := jwtplatform.New("secret", time.Hour)
	token, err := jwtManager.Issue(authmodel.DummyAdminID, authmodel.RoleAdmin)
	if err != nil {
		t.Fatalf("Issue token: %v", err)
	}

	roomID := uuid.New()
	req := httptest.NewRequest(http.MethodPost, "/rooms/"+roomID.String()+"/schedule/create", strings.NewReader(`{"daysOfWeek":[1,2],"startTime":"09:00","endTime":"10:00"}`))
	req.SetPathValue("roomId", roomID.String())
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	middleware.Authenticate(jwtManager)(handler.Create()).ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}
}

func TestCreateScheduleHandlerRejectsMissingRequiredFields(t *testing.T) {
	service := scheduleservice.New(scheduleRoomRepoStub{}, scheduleRepoHandlerStub{}, scheduleGeneratorStub{}, scheduleTxStub{}, scheduleClock{}, 30)
	handler := New(service)
	jwtManager := jwtplatform.New("secret", time.Hour)
	token, err := jwtManager.Issue(authmodel.DummyAdminID, authmodel.RoleAdmin)
	if err != nil {
		t.Fatalf("Issue token: %v", err)
	}

	roomID := uuid.New()
	req := httptest.NewRequest(http.MethodPost, "/rooms/"+roomID.String()+"/schedule/create", strings.NewReader(`{"daysOfWeek":[],"startTime":"09:00","endTime":"10:00"}`))
	req.SetPathValue("roomId", roomID.String())
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	middleware.Authenticate(jwtManager)(handler.Create()).ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}
