-- IMPORTANT: Drop ALL FOREIGN KEY constraints first before changing types
-- All room_id columns will now store Janus room ID (gateway_uid), not internal rooms.id

-- Drop FK constraints
ALTER TABLE sessions DROP CONSTRAINT IF EXISTS sessions_room_id_fkey;
ALTER TABLE room_statistics DROP CONSTRAINT IF EXISTS room_statistics_room_id_fkey;
ALTER TABLE composites_rooms DROP CONSTRAINT IF EXISTS composites_rooms_room_id_fkey;

-- Drop indexes (will recreate after type change)
DROP INDEX IF EXISTS sessions_room_id_idx;

-- Изменяем тип room_id в таблице sessions
ALTER TABLE sessions ALTER COLUMN room_id TYPE varchar(255);

-- Изменяем тип gateway_feed в таблице sessions
ALTER TABLE sessions ALTER COLUMN gateway_feed TYPE varchar(255);

-- Изменяем тип room_id в таблице room_statistics
ALTER TABLE room_statistics ALTER COLUMN room_id TYPE varchar(255);

-- Изменяем тип room_id в таблице composites_rooms
ALTER TABLE composites_rooms ALTER COLUMN room_id TYPE varchar(255);

-- Обновляем существующие данные (конвертируем внутренний rooms.id в Janus gateway_uid)
-- NOTE: This conversion assumes room_id currently contains rooms.id (internal ID)
-- We need to convert it to gateway_uid (Janus room ID)
UPDATE sessions s
SET room_id = (SELECT r.gateway_uid::varchar FROM rooms r WHERE r.id = s.room_id::bigint)
WHERE room_id IS NOT NULL AND room_id ~ '^\d+$';

UPDATE sessions SET gateway_feed = gateway_feed::varchar(255) WHERE gateway_feed IS NOT NULL;

UPDATE room_statistics rs
SET room_id = (SELECT r.gateway_uid::varchar FROM rooms r WHERE r.id = rs.room_id::bigint)
WHERE room_id IS NOT NULL AND room_id ~ '^\d+$';

UPDATE composites_rooms cr
SET room_id = (SELECT r.gateway_uid::varchar FROM rooms r WHERE r.id = cr.room_id::bigint)
WHERE room_id IS NOT NULL AND room_id ~ '^\d+$';

-- Recreate index on room_id (without FK constraint)
CREATE INDEX IF NOT EXISTS sessions_room_id_idx ON sessions USING btree (room_id, created_at); 