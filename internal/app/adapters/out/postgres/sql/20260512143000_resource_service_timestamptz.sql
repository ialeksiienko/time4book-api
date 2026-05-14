-- +goose Up
ALTER TABLE resources
    ALTER COLUMN unavailable_from TYPE TIMESTAMPTZ
    USING (unavailable_from::timestamp AT TIME ZONE 'Europe/Warsaw');

ALTER TABLE resources
    ALTER COLUMN unavailable_to TYPE TIMESTAMPTZ
    USING (unavailable_to::timestamp AT TIME ZONE 'Europe/Warsaw');

-- +goose Down
ALTER TABLE resources
    ALTER COLUMN unavailable_from TYPE DATE
    USING (unavailable_from AT TIME ZONE 'Europe/Warsaw')::date;

ALTER TABLE resources
    ALTER COLUMN unavailable_to TYPE DATE
    USING (unavailable_to AT TIME ZONE 'Europe/Warsaw')::date;
