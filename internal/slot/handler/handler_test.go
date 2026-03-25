package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	authmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/auth/model"
	jwtplatform "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/jwt"
	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/middleware"
	schedulemodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/schedule/model"
	slotmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/slot/model"
	slotservice "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/slot/service"
	"github.com/google/uuid"
)

type slotHandlerRoomRepo struct{}

func (slotHandlerRoomRepo) Exists(ctx context.Context, roomID uuid.UUID) (bool, error) {
	return true, nil
}

type slotHandlerScheduleRepo struct {
	schedule *schedulemodel.Schedule
}

func (s slotHandlerScheduleRepo) GetByRoomID(ctx context.Context, roomID uuid.UUID) (*schedulemodel.Schedule, error) {
	return s.schedule, nil
}
func (slotHandlerScheduleRepo) UpdateGeneratedUntil(ctx context.Context, roomID uuid.UUID, generatedUntil time.Time) error {
	return nil
}

type slotHandlerRepo struct {
	slots []slotmodel.Slot
}

func (s slotHandlerRepo) ListAvailableByRoomAndDate(ctx context.Context, roomID uuid.UUID, date time.Time) ([]slotmodel.Slot, error) {
	return s.slots, nil
}
func (slotHandlerRepo) CreateBatch(ctx context.Context, slots []slotmodel.Slot) error { return nil }
func (slotHandlerRepo) GetByID(ctx context.Context, slotID uuid.UUID) (*slotmodel.Slot, error) {
	return nil, nil
}

type slotHandlerClock struct{}

func (slotHandlerClock) Now() time.Time { return time.Date(2026, 3, 21, 10, 0, 0, 0, time.UTC) }

type slotHandlerTx struct{}

func (slotHandlerTx) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func TestListAvailableSlotsHandler(t *testing.T) {
	roomID := uuid.New()
	service := slotservice.New(
		slotHandlerRoomRepo{},
		slotHandlerScheduleRepo{schedule: &schedulemodel.Schedule{
			ID:             uuid.New(),
			RoomID:         roomID,
			DaysOfWeek:     []int{1},
			StartMinute:    540,
			EndMinute:      600,
			GeneratedUntil: time.Date(2026, 3, 25, 0, 0, 0, 0, time.UTC),
		}},
		slotHandlerRepo{slots: []slotmodel.Slot{{
			ID:      uuid.New(),
			RoomID:  roomID,
			StartAt: time.Date(2026, 3, 23, 9, 0, 0, 0, time.UTC),
			EndAt:   time.Date(2026, 3, 23, 9, 30, 0, 0, time.UTC),
		}}},
		slotHandlerTx{},
		slotHandlerClock{},
	)
	handler := New(service)
	jwtManager := jwtplatform.New("secret", time.Hour)
	token, err := jwtManager.Issue(authmodel.DummyUserID, authmodel.RoleUser)
	if err != nil {
		t.Fatalf("Issue token: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/rooms/"+roomID.String()+"/slots/list?date=2026-03-23", nil)
	req.SetPathValue("roomId", roomID.String())
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	middleware.Authenticate(jwtManager)(handler.ListAvailable()).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}
