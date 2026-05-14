# time4book-api

Project for university.

## Міграції БД (Goose)

Скопіюй значення з [`.env.example`](./.env.example) у `.env` і підстав реальний рядок підключення в `GOOSE_DBSTRING`.

З каталогу репозиторію:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
export PATH="$PATH:$(go env GOPATH)/bin"   # інакше zsh напише: command not found: goose
```

Після `source .env` (або `goose -env .env`) змінні `GOOSE_DRIVER`, `GOOSE_DBSTRING`, `GOOSE_MIGRATION_DIR` вже задані — тоді **лише**:

```bash
goose up
```

Не додавай у кінці рядка `postgres` і DSN окремими аргументами: синтаксис `goose -dir … postgres "…" up` у такому вигляді збиває парсер (помилка на кшталт `postgres: no such command`).

Якщо хочеш усе в одній команді **без** env — драйвер і DSN йдуть **спочатку**, `-dir` **після** DSN ([док Goose](https://github.com/pressly/goose)):

```bash
goose postgres "postgresql://user:pass@host:5432/dbname" -dir internal/app/adapters/out/postgres/sql up
```

Перевірка / відкат (після того ж `source .env`):

```bash
goose status
goose down
```

**Роль `developer`:** без `companyId` у запиті списки ресурсів, резервацій і користувачів повертають дані **усіх** компаній. Щоб обмежити однією фірмою, додай query `?companyId=<uuid>` або заголовок `X-Company-ID`.

## Seeded developer account

After running database migrations, a **developer** user is available for local use:

| Field    | Value                                      |
|----------|---------------------------------------------|
| Email    | `vladbevl@gmail.com`                       |
| Password | `123856vlad`                               |
| Role     | `developer`                                |
| Company  | `Sandbox (developer)` (for API scoping)  |

Applied by Goose migration [`internal/app/adapters/out/postgres/sql/20260509000000_seed_developer_vlad.sql`](internal/app/adapters/out/postgres/sql/20260509000000_seed_developer_vlad.sql).

If `vladbevl@gmail.com` already exists before this migration runs, inserts that target that email may be skipped; adjust manually if needed.

On app startup in non-production environments, developer bootstrap also synchronizes this account's credentials to the password above.

**Security:** rotate or remove this seed in shared or production databases.
