DROP INDEX IF EXISTS sessions_room_id_idx;
DROP INDEX IF EXISTS sessions_user_id_idx;

DROP TABLE IF EXISTS composites_rooms;
DROP TABLE IF EXISTS composites;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS rooms;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS gateways;

DROP FUNCTION IF EXISTS now_utc();
