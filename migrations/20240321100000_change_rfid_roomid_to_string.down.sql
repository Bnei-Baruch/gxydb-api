-- Revert room_id type in composites_rooms table
ALTER TABLE composites_rooms ALTER COLUMN room_id TYPE bigint USING room_id::bigint;

-- Revert room_id type in room_statistics table
ALTER TABLE room_statistics ALTER COLUMN room_id TYPE bigint USING room_id::bigint;

-- Revert gateway_feed type in sessions table
ALTER TABLE sessions ALTER COLUMN gateway_feed TYPE bigint USING gateway_feed::bigint;

-- Revert room_id type in sessions table
ALTER TABLE sessions ALTER COLUMN room_id TYPE bigint USING room_id::bigint;
