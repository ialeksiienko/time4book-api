-- +goose Up
ALTER TABLE reservations
ADD COLUMN IF NOT EXISTS company_id UUID;

UPDATE reservations r
SET company_id = res.company_id
FROM resources res
WHERE r.resource_id = res.id
  AND r.company_id IS NULL;

ALTER TABLE reservations
ALTER COLUMN company_id SET NOT NULL;

CREATE INDEX IF NOT EXISTS idx_reservations_company_id ON reservations(company_id);

-- +goose Down
DROP INDEX IF EXISTS idx_reservations_company_id;
ALTER TABLE reservations DROP COLUMN IF EXISTS company_id;
