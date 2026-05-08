-- +goose Up
CREATE TABLE roles(
    id VARCHAR(15) PRIMARY KEY NOT NULL,
    friendly_name VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ
);

INSERT INTO roles (id, friendly_name) VALUES 
('owner', 'Owner'),
('admin', 'Admin'),
('employee', 'Employee'),
('developer', 'Developer')
ON CONFLICT (id) DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS roles CASCADE;
