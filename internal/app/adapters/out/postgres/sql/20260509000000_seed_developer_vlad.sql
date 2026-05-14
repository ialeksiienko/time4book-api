-- +goose Up
-- Developer login: vladbevl@gmail.com (password in README).

INSERT INTO users (id, company_id, firstname, lastname, email, role_id, status, created_at)
SELECT
    'a1111111-1111-4111-a111-111111111111'::uuid,
    NULL,
    'Vlad',
    'Bevl',
    'vladbevl@gmail.com',
    'developer',
    'active',
    now()
WHERE NOT EXISTS (
    SELECT 1 FROM users WHERE email = 'vladbevl@gmail.com'
);

INSERT INTO user_credentials (user_id, email, password_hash)
SELECT
    u.id,
    'vladbevl@gmail.com',
    '$2a$14$1HK.7j1.aOvltUh0sZioMuAkWVftKdy.4Egn83pLolCXsm/bC4ghu'
FROM users u
WHERE u.email = 'vladbevl@gmail.com'
  AND NOT EXISTS (
      SELECT 1 FROM user_credentials c WHERE c.user_id = u.id
  );

INSERT INTO companies (
    id, owner_id, name, nip, address, industry, status, created_at
)
SELECT
    'b2222222-2222-4222-a222-222222222222'::uuid,
    u.id,
    'Sandbox (developer)',
    NULL,
    NULL,
    NULL,
    'active',
    now()
FROM users u
WHERE u.email = 'vladbevl@gmail.com'
  AND NOT EXISTS (
      SELECT 1 FROM companies WHERE id = 'b2222222-2222-4222-a222-222222222222'::uuid
  );

UPDATE users
SET company_id = 'b2222222-2222-4222-a222-222222222222'::uuid
WHERE email = 'vladbevl@gmail.com'
  AND company_id IS NULL;

-- +goose Down
UPDATE users SET company_id = NULL
WHERE email = 'vladbevl@gmail.com'
  AND company_id = 'b2222222-2222-4222-a222-222222222222'::uuid;

DELETE FROM companies WHERE id = 'b2222222-2222-4222-a222-222222222222'::uuid;
DELETE FROM user_credentials WHERE email = 'vladbevl@gmail.com';
DELETE FROM users WHERE email = 'vladbevl@gmail.com';
