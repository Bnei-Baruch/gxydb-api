-- Fix room_server_assignments to use Janus room ID (gateway_uid) instead of internal rooms.id
-- This aligns with sessions.room_id which now stores Janus room ID

-- 1. Drop existing table (since it's new and likely empty)
DROP TABLE IF EXISTS room_server_assignments;

-- 2. Recreate with VARCHAR room_id (stores Janus room ID, not internal rooms.id)
CREATE TABLE IF NOT EXISTS room_server_assignments
(
    room_id       VARCHAR(255) PRIMARY KEY            NOT NULL, -- Janus room ID (gateway_uid), e.g. "2171"
    gateway_name  VARCHAR(50)                         NOT NULL,
    region        VARCHAR(50),
    assigned_at   TIMESTAMP                           NOT NULL DEFAULT NOW(),
    last_used_at  TIMESTAMP                           NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS room_server_assignments_last_used_at_idx
    ON room_server_assignments USING BTREE (last_used_at);

-- Note: No FOREIGN KEY constraint because we're storing gateway_uid (string) not rooms.id (int64)
-- This is intentional for Janus string ID support
