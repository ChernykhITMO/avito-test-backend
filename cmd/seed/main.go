package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	adminUserID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	regularUser = uuid.MustParse("00000000-0000-0000-0000-000000000002")
)

type roomSeed struct {
	id          uuid.UUID
	scheduleID  uuid.UUID
	name        string
	description string
	capacity    int
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		fmt.Fprintln(os.Stderr, "DATABASE_URL is required")
		os.Exit(1)
	}

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "connect db: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "ping db: %v\n", err)
		os.Exit(1)
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "begin tx: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if err := seedUsers(ctx, tx); err != nil {
		fmt.Fprintf(os.Stderr, "seed users: %v\n", err)
		os.Exit(1)
	}

	rooms := []roomSeed{
		{
			id:          uuid.MustParse("10000000-0000-0000-0000-000000000001"),
			scheduleID:  uuid.MustParse("20000000-0000-0000-0000-000000000001"),
			name:        "Blue Room",
			description: "Main floor meeting room",
			capacity:    8,
		},
		{
			id:          uuid.MustParse("10000000-0000-0000-0000-000000000002"),
			scheduleID:  uuid.MustParse("20000000-0000-0000-0000-000000000002"),
			name:        "Green Room",
			description: "Second floor meeting room",
			capacity:    6,
		},
		{
			id:          uuid.MustParse("10000000-0000-0000-0000-000000000003"),
			scheduleID:  uuid.MustParse("20000000-0000-0000-0000-000000000003"),
			name:        "Board Room",
			description: "Large room for workshops",
			capacity:    12,
		},
	}

	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	generatedUntil := today.AddDate(0, 0, 14)

	for _, room := range rooms {
		if err := seedRoomAndSchedule(ctx, tx, room, generatedUntil); err != nil {
			fmt.Fprintf(os.Stderr, "seed room/schedule: %v\n", err)
			os.Exit(1)
		}
	}

	for _, room := range rooms {
		if err := seedSlots(ctx, tx, room, today, today.AddDate(0, 0, 6), now); err != nil {
			fmt.Fprintf(os.Stderr, "seed slots: %v\n", err)
			os.Exit(1)
		}
	}

	if err := seedBookings(ctx, tx, rooms[0].id, regularUser, today, now); err != nil {
		fmt.Fprintf(os.Stderr, "seed bookings: %v\n", err)
		os.Exit(1)
	}

	if err := tx.Commit(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "commit tx: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("seed completed")
}

func seedUsers(ctx context.Context, tx pgx.Tx) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO users (id, email, role)
		VALUES
			($1, 'admin@example.com', 'admin'),
			($2, 'user@example.com', 'user')
		ON CONFLICT (id) DO NOTHING
	`, adminUserID, regularUser)
	return err
}

func seedRoomAndSchedule(ctx context.Context, tx pgx.Tx, room roomSeed, generatedUntil time.Time) error {
	if _, err := tx.Exec(ctx, `
		INSERT INTO rooms (id, name, description, capacity)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO NOTHING
	`, room.id, room.name, room.description, room.capacity); err != nil {
		return err
	}

	_, err := tx.Exec(ctx, `
		INSERT INTO schedules (id, room_id, days_of_week, start_minute, end_minute, generated_until)
		VALUES ($1, $2, ARRAY[1,2,3,4,5,6,7]::smallint[], 540, 1080, $3)
		ON CONFLICT (room_id) DO UPDATE
		SET generated_until = GREATEST(schedules.generated_until, EXCLUDED.generated_until)
	`, room.scheduleID, room.id, generatedUntil)
	return err
}

func seedSlots(ctx context.Context, tx pgx.Tx, room roomSeed, fromDate, toDate, now time.Time) error {
	for date := fromDate; !date.After(toDate); date = date.AddDate(0, 0, 1) {
		for minute := 540; minute < 1080; minute += 30 {
			startAt := date.Add(time.Duration(minute) * time.Minute)
			endAt := startAt.Add(30 * time.Minute)

			_, err := tx.Exec(ctx, `
				INSERT INTO slots (id, room_id, schedule_id, slot_date, start_at, end_at, created_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
				ON CONFLICT (room_id, start_at) DO NOTHING
			`, uuid.New(), room.id, room.scheduleID, date, startAt, endAt, now)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func seedBookings(ctx context.Context, tx pgx.Tx, roomID, userID uuid.UUID, today, now time.Time) error {
	rows, err := tx.Query(ctx, `
		SELECT id
		FROM slots
		WHERE room_id = $1
		  AND slot_date = $2
		  AND start_at > $3
		ORDER BY start_at ASC
		LIMIT 2
	`, roomID, today.AddDate(0, 0, 1), now)
	if err != nil {
		return err
	}
	defer rows.Close()

	slotIDs := make([]uuid.UUID, 0, 2)
	for rows.Next() {
		var slotID uuid.UUID
		if err := rows.Scan(&slotID); err != nil {
			return err
		}
		slotIDs = append(slotIDs, slotID)
	}
	if rows.Err() != nil {
		return rows.Err()
	}

	for _, slotID := range slotIDs {
		_, err := tx.Exec(ctx, `
			INSERT INTO bookings (id, slot_id, user_id, status, created_at)
			VALUES ($1, $2, $3, 'active', $4)
			ON CONFLICT (slot_id) WHERE status = 'active' DO NOTHING
		`, uuid.New(), slotID, userID, now)
		if err != nil {
			return err
		}
	}

	return nil
}
