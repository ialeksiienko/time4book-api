-- +goose Up
CREATE TABLE resources (
    id UUID PRIMARY KEY NOT NULL,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    name VARCHAR(128) NOT NULL,
    type VARCHAR(20) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    location VARCHAR(128) NOT NULL DEFAULT '',
    max_reservation_minutes INT,
    available_from VARCHAR(5),
    available_to VARCHAR(5),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    unavailable_from DATE,
    unavailable_to DATE,
    unavailable_reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS resources CASCADE;
