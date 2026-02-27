-- Миграция для поддержки string room_id в Janus
-- Этап 1: Изменяем gateway_uid с INTEGER на VARCHAR

-- Изменяем тип gateway_uid в таблице rooms
ALTER TABLE rooms ALTER COLUMN gateway_uid TYPE varchar(255);

-- Обновляем существующие данные (int → string)
UPDATE rooms SET gateway_uid = gateway_uid::varchar(255);

-- Комментарий для будущего:
-- После этой миграции rooms.gateway_uid может хранить как числа ("2171")
-- так и string ID от Janus ("room-uuid-123")
