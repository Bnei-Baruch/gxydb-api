-- Fix room_server_assignments.room_id type to VARCHAR to match sessions.room_id
-- Drop FK constraint and recreate table with correct type

DROP TABLE IF EXISTS room_server_assignments CASCADE;

CREATE TABLE IF NOT EXISTS room_server_assignments
(
    room_id       VARCHAR(255) PRIMARY KEY NOT NULL, -- Janus room ID (gateway_uid)
    gateway_name  VARCHAR(50)              NOT NULL,
    region        VARCHAR(50),
    assigned_at   TIMESTAMP                NOT NULL DEFAULT NOW(),
    last_used_at  TIMESTAMP                NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS room_server_assignments_last_used_at_idx
    ON room_server_assignments USING BTREE (last_used_at);
