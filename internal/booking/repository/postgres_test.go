package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	bookingmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/booking/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	pgxmock "github.com/pashagolub/pgxmock/v4"
)

func TestBookingRepositoryQueries(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("NewPool: %v", err)
	}
	defer mock.Close()

	repo := New(mock)
	bookingID := uuid.New()
	slotID := uuid.New()
	userID := uuid.New()
	now := time.Now().UTC()
	link := "https://conference.local/test"

	mock.ExpectExec("INSERT INTO bookings").
		WithArgs(
			bookingID,
			slotID,
			userID,
			bookingmodel.StatusActive,
			&link,
			now,
			pgxmock.AnyArg(),
		).
		WillReturnError(&pgconn.PgError{Code: "23505"})

	err = repo.Create(context.Background(), bookingmodel.Booking{
		ID:             bookingID,
		SlotID:         slotID,
		UserID:         userID,
		Status:         bookingmodel.StatusActive,
		ConferenceLink: &link,
		CreatedAt:      now,
	})
	if !errors.Is(err, bookingmodel.ErrSlotAlreadyBooked) {
		t.Fatalf("expected slot already booked, got %v", err)
	}

	rows := pgxmock.NewRows([]string{"id", "slot_id", "user_id", "status", "conference_link", "created_at", "cancelled_at"}).
		AddRow(bookingID, slotID, userID, bookingmodel.StatusActive, link, now, nil)
	mock.ExpectQuery("SELECT id, slot_id, user_id, status, conference_link, created_at, cancelled_at FROM bookings").
		WithArgs(bookingID).
		WillReturnRows(rows)

	booking, err := repo.GetByID(context.Background(), bookingID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if booking == nil || booking.ID != bookingID {
		t.Fatalf("unexpected booking: %#v", booking)
	}

	mock.ExpectExec("UPDATE bookings SET status = 'cancelled'").
		WithArgs(bookingID, now).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	if err := repo.Cancel(context.Background(), bookingID, now); err != nil {
		t.Fatalf("Cancel: %v", err)
	}

	myRows := pgxmock.NewRows([]string{"id", "slot_id", "user_id", "status", "conference_link", "created_at", "cancelled_at"}).
		AddRow(bookingID, slotID, userID, bookingmodel.StatusActive, link, now, nil)
	mock.ExpectQuery("SELECT b.id, b.slot_id, b.user_id, b.status, b.conference_link, b.created_at, b.cancelled_at").
		WithArgs(userID, now).
		WillReturnRows(myRows)

	if _, err := repo.ListMyFuture(context.Background(), userID, now); err != nil {
		t.Fatalf("ListMyFuture: %v", err)
	}

	countRows := pgxmock.NewRows([]string{"count"}).AddRow(1)
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM bookings").WillReturnRows(countRows)

	allRows := pgxmock.NewRows([]string{"id", "slot_id", "user_id", "status", "conference_link", "created_at", "cancelled_at"}).
		AddRow(bookingID, slotID, userID, bookingmodel.StatusActive, link, now, nil)
	mock.ExpectQuery("SELECT id, slot_id, user_id, status, conference_link, created_at, cancelled_at FROM bookings").
		WithArgs(20, 0).
		WillReturnRows(allRows)

	if _, _, err := repo.ListAll(context.Background(), 1, 20); err != nil {
		t.Fatalf("ListAll: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}
