-- Откат миграции string room_id
-- ВНИМАНИЕ: Это работает только если все gateway_uid это числа!

ALTER TABLE rooms ALTER COLUMN gateway_uid TYPE integer USING gateway_uid::integer;
