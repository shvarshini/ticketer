ALTER TABLE bookings DROP CONSTRAINT fk_bookings_user_id;
ALTER TABLE bookings ALTER COLUMN user_id TYPE VARCHAR(255);

ALTER TABLE theaters DROP CONSTRAINT fk_theaters_admin_id;
ALTER TABLE theaters ALTER COLUMN admin_id TYPE VARCHAR(255);

DROP TABLE IF EXISTS users;
