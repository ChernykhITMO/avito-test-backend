-- +goose Up
CREATE UNIQUE INDEX slots_room_start_at_uidx
    ON slots (room_id, start_at);

CREATE INDEX slots_room_date_start_at_idx
    ON slots (room_id, slot_date, start_at);

CREATE UNIQUE INDEX bookings_active_slot_uidx
    ON bookings (slot_id)
    WHERE status = 'active';

CREATE INDEX bookings_user_created_at_idx
    ON bookings (user_id, created_at DESC);

CREATE INDEX bookings_created_at_idx
    ON bookings (created_at DESC, id DESC);

-- +goose Down
DROP INDEX IF EXISTS bookings_created_at_idx;
DROP INDEX IF EXISTS bookings_user_created_at_idx;
DROP INDEX IF EXISTS bookings_active_slot_uidx;
DROP INDEX IF EXISTS slots_room_date_start_at_idx;
DROP INDEX IF EXISTS slots_room_start_at_uidx;
