-- +goose Up
CREATE TABLE users(
    id UUID PRIMARY KEY NOT NULL,
    company_id UUID REFERENCES companies(id) ON DELETE SET NULL,
    firstname VARCHAR(20) NOT NULL,
    lastname VARCHAR(35) NOT NULL,
    email TEXT NOT NULL UNIQUE,
    role_id VARCHAR(50) NOT NULL REFERENCES roles(id) ON DELETE RESTRICT,
    status VARCHAR(10) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ
);

-- +goose Down
DROP TABLE IF EXISTS users CASCADE;
