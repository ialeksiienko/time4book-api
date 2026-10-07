-- +goose Up
CREATE TABLE companies(
    id UUID PRIMARY KEY NOT NULL,
    owner_id UUID NOT NULL,
    name VARCHAR(128) NOT NULL,
    nip VARCHAR(10),
    address VARCHAR(128),
    industry VARCHAR(32),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ
);

-- +goose Down
DROP TABLE IF EXISTS companies CASCADE;
