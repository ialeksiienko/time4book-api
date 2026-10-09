-- +goose Up
CREATE EXTENSION IF NOT EXISTS btree_gist;

ALTER TABLE reservations
    ADD CONSTRAINT reservations_valid_interval CHECK (end_date > start_date);

ALTER TABLE reservations
    ADD CONSTRAINT reservations_active_no_overlap
    EXCLUDE USING gist (
        resource_id WITH =,
        tstzrange(start_date, end_date, '[)') WITH &&
    ) WHERE (status = 'active');

-- +goose Down
ALTER TABLE reservations DROP CONSTRAINT reservations_active_no_overlap;
ALTER TABLE reservations DROP CONSTRAINT reservations_valid_interval;
