-- +goose Up
CREATE TABLE reservations (
    id UUID PRIMARY KEY NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    resource_id UUID NOT NULL REFERENCES resources(id) ON DELETE RESTRICT,
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,
    description TEXT,
    status VARCHAR(25) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_reservations_resource_id ON reservations(resource_id);
CREATE INDEX idx_reservations_user_id ON reservations(user_id);
CREATE INDEX idx_reservations_status ON reservations(status);

-- +goose Down
DROP TABLE IF EXISTS reservations CASCADE;
