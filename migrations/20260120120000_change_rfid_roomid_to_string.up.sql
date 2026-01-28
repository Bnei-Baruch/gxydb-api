-- Изменяем тип room_id в таблице sessions
ALTER TABLE sessions ALTER COLUMN room_id TYPE varchar(255);

-- Изменяем тип gateway_feed в таблице sessions
ALTER TABLE sessions ALTER COLUMN gateway_feed TYPE varchar(255);

-- Изменяем тип room_id в таблице room_statistics
ALTER TABLE room_statistics ALTER COLUMN room_id TYPE varchar(255);

-- Изменяем тип room_id в таблице composites_rooms
ALTER TABLE composites_rooms ALTER COLUMN room_id TYPE varchar(255);

-- Обновляем существующие данные
UPDATE sessions SET room_id = room_id::varchar(255);
UPDATE sessions SET gateway_feed = gateway_feed::varchar(255);
UPDATE room_statistics SET room_id = room_id::varchar(255);
UPDATE composites_rooms SET room_id = room_id::varchar(255); 