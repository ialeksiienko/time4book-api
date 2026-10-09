-- +goose Up
CREATE TABLE user_credentials(
    user_id UUID PRIMARY KEY NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS user_credentials CASCADE;
