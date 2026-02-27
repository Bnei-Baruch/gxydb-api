-- Drop FK constraints

ALTER TABLE sessions DROP CONSTRAINT IF EXISTS sessions_room_id_fkey;
ALTER TABLE room_statistics DROP CONSTRAINT IF EXISTS room_statistics_room_id_fkey;
ALTER TABLE composites_rooms DROP CONSTRAINT IF EXISTS composites_rooms_room_id_fkey;
