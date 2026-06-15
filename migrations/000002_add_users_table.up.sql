CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    roles TEXT[] NOT NULL DEFAULT '{user}',
    oauth_provider VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Note: In PostgreSQL, casting VARCHAR to UUID might require a USING clause if there's existing data that is not a valid UUID string.
-- Assuming no data or existing data are valid UUID strings. If not, this might fail or require data cleanup first.
ALTER TABLE theaters ALTER COLUMN admin_id TYPE UUID USING admin_id::UUID;
ALTER TABLE theaters ADD CONSTRAINT fk_theaters_admin_id FOREIGN KEY (admin_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE bookings ALTER COLUMN user_id TYPE UUID USING user_id::UUID;
ALTER TABLE bookings ADD CONSTRAINT fk_bookings_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
