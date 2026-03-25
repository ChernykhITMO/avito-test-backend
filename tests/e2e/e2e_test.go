package e2e

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	authhandler "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/auth/handler"
	authmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/auth/model"
	authservice "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/auth/service"
	bookinghandler "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/booking/handler"
	bookingmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/booking/model"
	bookingservice "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/booking/service"
	conferencemock "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/conference/mock"
	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/httpcommon"
	jwtplatform "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/jwt"
	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/middleware"
	roomhandler "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/room/handler"
	roommodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/room/model"
	roomservice "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/room/service"
	schedulehandler "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/schedule/handler"
	schedulemodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/schedule/model"
	scheduleservice "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/schedule/service"
	slothandler "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/slot/handler"
	slotmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/slot/model"
	slotservice "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/slot/service"
	"github.com/google/uuid"
)

type testApp struct {
	handler http.Handler
}

type noopTx struct{}

func (noopTx) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

type fixedClock struct {
	now time.Time
}

func (c fixedClock) Now() time.Time { return c.now }

type authRepoStub struct{}

func (authRepoStub) Create(ctx context.Context, params authmodel.CreateUserParams) error {
	return nil
}
func (authRepoStub) GetByEmail(ctx context.Context, email string) (*authmodel.UserWithPassword, error) {
	return nil, nil
}

type passwordStub struct{}

func (passwordStub) Hash(password string) (string, error) { return "hash", nil }
func (passwordStub) Compare(hash, password string) error  { return nil }

type store struct {
	mu        sync.Mutex
	rooms     map[uuid.UUID]roommodel.Room
	schedules map[uuid.UUID]schedulemodel.Schedule
	slots     map[uuid.UUID]slotmodel.Slot
	bookings  map[uuid.UUID]bookingmodel.Booking
}

func newStore() *store {
	return &store{
		rooms:     make(map[uuid.UUID]roommodel.Room),
		schedules: make(map[uuid.UUID]schedulemodel.Schedule),
		slots:     make(map[uuid.UUID]slotmodel.Slot),
		bookings:  make(map[uuid.UUID]bookingmodel.Booking),
	}
}

type roomRepo struct{ s *store }

func (r roomRepo) Create(ctx context.Context, room roommodel.Room) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.rooms[room.ID] = room
	return nil
}
func (r roomRepo) List(ctx context.Context) ([]roommodel.Room, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	result := make([]roommodel.Room, 0, len(r.s.rooms))
	for _, room := range r.s.rooms {
		result = append(result, room)
	}
	return result, nil
}
func (r roomRepo) Exists(ctx context.Context, roomID uuid.UUID) (bool, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	_, ok := r.s.rooms[roomID]
	return ok, nil
}

type scheduleRepo struct{ s *store }

func (r scheduleRepo) Create(ctx context.Context, schedule schedulemodel.Schedule) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	if _, ok := r.s.schedules[schedule.RoomID]; ok {
		return schedulemodel.ErrScheduleAlreadyExists
	}
	r.s.schedules[schedule.RoomID] = schedule
	return nil
}
func (r scheduleRepo) GetByRoomID(ctx context.Context, roomID uuid.UUID) (*schedulemodel.Schedule, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	schedule, ok := r.s.schedules[roomID]
	if !ok {
		return nil, nil
	}
	return &schedule, nil
}
func (r scheduleRepo) UpdateGeneratedUntil(ctx context.Context, roomID uuid.UUID, generatedUntil time.Time) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	schedule := r.s.schedules[roomID]
	schedule.GeneratedUntil = generatedUntil
	r.s.schedules[roomID] = schedule
	return nil
}

type slotRepo struct{ s *store }

func (r slotRepo) ListAvailableByRoomAndDate(ctx context.Context, roomID uuid.UUID, date time.Time) ([]slotmodel.Slot, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()

	active := make(map[uuid.UUID]struct{})
	for _, booking := range r.s.bookings {
		if booking.Status == bookingmodel.StatusActive {
			active[booking.SlotID] = struct{}{}
		}
	}

	result := make([]slotmodel.Slot, 0)
	for _, slot := range r.s.slots {
		if slot.RoomID == roomID && sameDate(slot.SlotDate, date) {
			if _, taken := active[slot.ID]; !taken {
				result = append(result, slot)
			}
		}
	}
	return result, nil
}
func (r slotRepo) CreateBatch(ctx context.Context, slots []slotmodel.Slot) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	for _, slot := range slots {
		r.s.slots[slot.ID] = slot
	}
	return nil
}
func (r slotRepo) GetByID(ctx context.Context, slotID uuid.UUID) (*slotmodel.Slot, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	slot, ok := r.s.slots[slotID]
	if !ok {
		return nil, nil
	}
	return &slot, nil
}

type bookingRepo struct{ s *store }

func (r bookingRepo) Create(ctx context.Context, booking bookingmodel.Booking) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	for _, existing := range r.s.bookings {
		if existing.SlotID == booking.SlotID && existing.Status == bookingmodel.StatusActive {
			return bookingmodel.ErrSlotAlreadyBooked
		}
	}
	r.s.bookings[booking.ID] = booking
	return nil
}
func (r bookingRepo) GetByID(ctx context.Context, bookingID uuid.UUID) (*bookingmodel.Booking, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	booking, ok := r.s.bookings[bookingID]
	if !ok {
		return nil, nil
	}
	return &booking, nil
}
func (r bookingRepo) GetByIDForUpdate(ctx context.Context, bookingID uuid.UUID) (*bookingmodel.Booking, error) {
	return r.GetByID(ctx, bookingID)
}
func (r bookingRepo) Cancel(ctx context.Context, bookingID uuid.UUID, cancelledAt time.Time) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	booking := r.s.bookings[bookingID]
	booking.Status = bookingmodel.StatusCancelled
	booking.CancelledAt = &cancelledAt
	r.s.bookings[bookingID] = booking
	return nil
}
func (r bookingRepo) ListMyFuture(ctx context.Context, userID uuid.UUID, now time.Time) ([]bookingmodel.Booking, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	result := make([]bookingmodel.Booking, 0)
	for _, booking := range r.s.bookings {
		slot := r.s.slots[booking.SlotID]
		if booking.UserID == userID && !slot.StartAt.Before(now) {
			result = append(result, booking)
		}
	}
	return result, nil
}
func (r bookingRepo) ListAll(ctx context.Context, page, pageSize int) ([]bookingmodel.Booking, int, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	result := make([]bookingmodel.Booking, 0, len(r.s.bookings))
	for _, booking := range r.s.bookings {
		result = append(result, booking)
	}
	return result, len(result), nil
}

func TestCreateRoomScheduleAndBookingFlow(t *testing.T) {
	app := newTestApp(t)

	adminToken := loginAs(t, app, "admin")
	userToken := loginAs(t, app, "user")

	roomID := createRoom(t, app, adminToken)
	createSchedule(t, app, adminToken, roomID)
	slotID := getFirstSlot(t, app, userToken, roomID)
	bookingID := createBooking(t, app, userToken, slotID)

	if bookingID == "" {
		t.Fatal("expected booking id")
	}
}

func TestCancelBookingFlow(t *testing.T) {
	app := newTestApp(t)

	adminToken := loginAs(t, app, "admin")
	userToken := loginAs(t, app, "user")

	roomID := createRoom(t, app, adminToken)
	createSchedule(t, app, adminToken, roomID)
	slotID := getFirstSlot(t, app, userToken, roomID)
	bookingID := createBooking(t, app, userToken, slotID)

	req := httptest.NewRequest(http.MethodPost, "/bookings/"+bookingID+"/cancel", nil)
	req.SetPathValue("bookingId", bookingID)
	req.Header.Set("Authorization", "Bearer "+userToken)
	rec := httptest.NewRecorder()
	app.handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func newTestApp(t *testing.T) *testApp {
	t.Helper()

	store := newStore()
	now := time.Date(2026, 3, 23, 8, 0, 0, 0, time.UTC)
	clock := fixedClock{now: now}
	tx := noopTx{}
	jwtManager := jwtplatform.New("secret", time.Hour)

	roomRepo := roomRepo{s: store}
	scheduleRepo := scheduleRepo{s: store}
	slotRepo := slotRepo{s: store}
	bookingRepo := bookingRepo{s: store}

	authSvc := authservice.New(authRepoStub{}, jwtManager, passwordStub{}, tx)
	roomSvc := roomservice.New(roomRepo)
	slotSvc := slotservice.New(roomRepo, scheduleRepo, slotRepo, tx, clock)
	scheduleSvc := scheduleservice.New(roomRepo, scheduleRepo, slotSvc, tx, clock, 30)
	bookingSvc := bookingservice.New(bookingRepo, slotRepo, conferencemock.New("https://conference.local/booking/", "ok"), clock, tx)

	authHandler := authhandler.New(authSvc)
	roomHandler := roomhandler.New(roomSvc)
	scheduleHandler := schedulehandler.New(scheduleSvc)
	slotHandler := slothandler.New(slotSvc)
	bookingHandler := bookinghandler.New(bookingSvc)

	mux := http.NewServeMux()
	authenticated := middleware.Authenticate(jwtManager)
	mux.Handle("POST /dummyLogin", authHandler.DummyLogin())
	mux.Handle("POST /rooms/create", authenticated(roomHandler.Create()))
	mux.Handle("POST /rooms/{roomId}/schedule/create", authenticated(scheduleHandler.Create()))
	mux.Handle("GET /rooms/{roomId}/slots/list", authenticated(slotHandler.ListAvailable()))
	mux.Handle("POST /bookings/create", authenticated(bookingHandler.Create()))
	mux.Handle("POST /bookings/{bookingId}/cancel", authenticated(bookingHandler.Cancel()))
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		httpcommon.WriteError(w, http.StatusNotFound, httpcommon.CodeNotFound, "route not found")
	}))

	return &testApp{handler: mux}
}

func loginAs(t *testing.T, app *testApp, role string) string {
	t.Helper()

	rec := app.do(t, http.MethodPost, "/dummyLogin", "", strings.NewReader(`{"role":"`+role+`"}`))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var body struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode token: %v", err)
	}
	return body.Token
}

func createRoom(t *testing.T, app *testApp, token string) string {
	t.Helper()

	rec := app.do(t, http.MethodPost, "/rooms/create", token, strings.NewReader(`{"name":"Blue","capacity":8}`))
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}

	var body struct {
		Room struct {
			ID string `json:"id"`
		} `json:"room"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode room: %v", err)
	}
	return body.Room.ID
}

func createSchedule(t *testing.T, app *testApp, token, roomID string) {
	t.Helper()

	req := httptest.NewRequest(http.MethodPost, "/rooms/"+roomID+"/schedule/create", strings.NewReader(`{"daysOfWeek":[1],"startTime":"09:00","endTime":"10:00"}`))
	req.SetPathValue("roomId", roomID)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	app.handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
}

func getFirstSlot(t *testing.T, app *testApp, token, roomID string) string {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, "/rooms/"+roomID+"/slots/list?date=2026-03-23", nil)
	req.SetPathValue("roomId", roomID)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	app.handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var body struct {
		Slots []struct {
			ID string `json:"id"`
		} `json:"slots"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode slots: %v", err)
	}
	if len(body.Slots) == 0 {
		t.Fatal("expected at least one slot")
	}
	return body.Slots[0].ID
}

func createBooking(t *testing.T, app *testApp, token, slotID string) string {
	t.Helper()

	rec := app.do(t, http.MethodPost, "/bookings/create", token, strings.NewReader(`{"slotId":"`+slotID+`"}`))
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}

	var body struct {
		Booking struct {
			ID string `json:"id"`
		} `json:"booking"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode booking: %v", err)
	}
	return body.Booking.ID
}

func (a *testApp) do(t *testing.T, method, path, token string, body io.Reader) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(method, path, body)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	rec := httptest.NewRecorder()
	a.handler.ServeHTTP(rec, req)
	return rec
}

func sameDate(a, b time.Time) bool {
	ua := a.UTC()
	ub := b.UTC()
	return ua.Year() == ub.Year() && ua.YearDay() == ub.YearDay()
}
