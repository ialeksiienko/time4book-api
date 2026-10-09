# Time4Book API

Go API for companies, resources, and reservations.

## Code layout

- `internal/app/core/domain/model`: business state and rules.
- `internal/app/core/usecases`: commands that change state.
- `internal/app/core/queries`: read operations and their result types.
- `internal/app/core/ports/in`: interfaces and request/response types used by HTTP.
- `internal/app/core/ports/out`: interfaces used for storage, tokens, and transactions.
- `internal/app/adapters/in/http`: HTTP routes and handlers.
- `internal/app/adapters/out/postgres`: database code and migrations.
- `internal/app/adapters/out/jwt`: token code.
- `internal/app/core/core.go`: connects commands and queries to the input ports.
- `internal/app/app.go`: connects the core to the adapters.

The HTTP adapter calls input ports. Commands and queries use output ports. The domain does not import adapters.

## Local checks

Run `go test ./...` and `go vet ./...`. Database tests need a PostgreSQL database with the migrations applied. Set `TIME4BOOK_TEST_DATABASE_URL` to run them; without it, they are skipped.

The migration files are in `internal/app/adapters/out/postgres/sql`. For local work, copy `.env.example` to `.env` and fill in your values. Values already set in the process environment take priority over `.env`.

The Docker image has no `.env` file. Set `DB_USER`, `DB_PASS`, `DB_HOST`, `DB_PORT`, `DB_NAME`, `JWT_SECRET`, `JWT_ACCESS_TOKEN_DURATION`, and `JWT_REFRESH_TOKEN_DURATION` when starting the container. The image sets `ENV=prod`. The app checks these values at startup and stops if one is missing or invalid. Docker can pass them with `--env` or `--env-file`; keep secret values out of the Dockerfile and image.
