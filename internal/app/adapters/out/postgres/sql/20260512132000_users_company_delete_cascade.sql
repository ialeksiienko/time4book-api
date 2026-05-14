-- +goose Up
-- Clean up orphaned local users left behind by earlier company deletions.
DELETE FROM reservations
WHERE user_id IN (
    SELECT id
    FROM users
    WHERE company_id IS NULL
      AND role_id <> 'developer'
);

DELETE FROM users
WHERE company_id IS NULL
  AND role_id <> 'developer';

ALTER TABLE users
    DROP CONSTRAINT IF EXISTS users_company_id_fkey;

ALTER TABLE users
    ADD CONSTRAINT users_company_id_fkey
        FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE users
    DROP CONSTRAINT IF EXISTS users_company_id_fkey;

ALTER TABLE users
    ADD CONSTRAINT users_company_id_fkey
        FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE SET NULL;
