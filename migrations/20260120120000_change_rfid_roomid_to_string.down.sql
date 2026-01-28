-- Откат изменений типов обратно на bigint
-- WARNING: This will convert gateway_uid back to internal rooms.id (data loss!)

-- Convert room_id from gateway_uid back to internal rooms.id
UPDATE sessions s
SET room_id = (SELECT r.id::varchar FROM rooms r WHERE r.gateway_uid = s.room_id)
WHERE room_id IS NOT NULL;

UPDATE room_statistics rs
SET room_id = (SELECT r.id::varchar FROM rooms r WHERE r.gateway_uid = rs.room_id)
WHERE room_id IS NOT NULL;

UPDATE composites_rooms cr
SET room_id = (SELECT r.id::varchar FROM rooms r WHERE r.gateway_uid = cr.room_id)
WHERE room_id IS NOT NULL;

-- Откат room_id в таблице composites_rooms
ALTER TABLE composites_rooms ALTER COLUMN room_id TYPE bigint USING room_id::bigint;

-- Откат room_id в таблице room_statistics
ALTER TABLE room_statistics ALTER COLUMN room_id TYPE bigint USING room_id::bigint;

-- Откат gateway_feed в таблице sessions
ALTER TABLE sessions ALTER COLUMN gateway_feed TYPE bigint USING gateway_feed::bigint;

-- Откат room_id в таблице sessions
ALTER TABLE sessions ALTER COLUMN room_id TYPE bigint USING room_id::bigint;

-- Restore FOREIGN KEY constraint
ALTER TABLE sessions ADD CONSTRAINT sessions_room_id_fkey FOREIGN KEY (room_id) REFERENCES rooms(id);

-- Restore index
CREATE INDEX IF NOT EXISTS sessions_room_id_idx ON sessions USING btree (room_id, created_at);
