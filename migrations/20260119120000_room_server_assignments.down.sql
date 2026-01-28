-- Revert room_server_assignments to use internal rooms.id

DROP TABLE IF EXISTS room_server_assignments;

CREATE TABLE IF NOT EXISTS room_server_assignments
(
    room_id       BIGINT PRIMARY KEY REFERENCES rooms NOT NULL,
    gateway_name  VARCHAR(50)                         NOT NULL,
    region        VARCHAR(50),
    assigned_at   TIMESTAMP                           NOT NULL DEFAULT NOW(),
    last_used_at  TIMESTAMP                           NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS room_server_assignments_last_used_at_idx
    ON room_server_assignments USING BTREE (last_used_at);
