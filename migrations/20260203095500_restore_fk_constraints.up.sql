-- Restore FK constraints after changing room_id to VARCHAR
-- Now FK points to rooms.gateway_uid instead of rooms.id

-- Add FK constraint from sessions.room_id to rooms.gateway_uid
ALTER TABLE sessions 
ADD CONSTRAINT sessions_room_id_fkey 
FOREIGN KEY (room_id) REFERENCES rooms(gateway_uid);

-- Add FK constraint from room_statistics.room_id to rooms.gateway_uid
ALTER TABLE room_statistics 
ADD CONSTRAINT room_statistics_room_id_fkey 
FOREIGN KEY (room_id) REFERENCES rooms(gateway_uid);

-- Add FK constraint from composites_rooms.room_id to rooms.gateway_uid
ALTER TABLE composites_rooms 
ADD CONSTRAINT composites_rooms_room_id_fkey 
FOREIGN KEY (room_id) REFERENCES rooms(gateway_uid);
