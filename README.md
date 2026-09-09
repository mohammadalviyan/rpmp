# RPMP

RPA Performance & Monitoring Platform. Internal information layer over Orchestrator / RPA operational data. MVP is Dashboard plus Use Cases.

RPMP is not a replacement for Orchestrator. It does not start or stop robots. Product language is RPMP (use case, execution, health, failure), not vendor job names.

Requirements live in [`docs/PRD.md`](docs/PRD.md). Cite `FR-x.y` (login is `FR-1.1` / `FR-1.2`, dashboard KPIs start at `FR-2.2`). Do not paste the PRD into tickets or agent prompts.

## Stack

- Next.js App Router in `frontend/` (usually `http://localhost:3000`)
- Go HTTP API in `backend/` (`RPMP_ADDR` defaults to `:8080`)
- PostgreSQL for RPMP-owned identity and audit data

## Clone and branches

Never push `main` / `master`. Work on a branch. Prefixes: `feature/`, `fix/`, `refactor/`, `docs/`, `chore/`.

## Run locally

PostgreSQL must exist first. Apply `backend/db/migrations/000001_auth.up.sql` (for example `psql "$RPMP_DATABASE_URL" -f backend/db/migrations/000001_auth.up.sql`). There is no migrate wrapper in the repo.

There is no seed CLI. Login needs a `users` row whose `password_hash` is bcrypt. Hash with `backend/internal/auth` (`Passwords.Hash`, `golang.org/x/crypto/bcrypt`). Frozen JSON uses employee ID `12345678`. That is an ID, not a default password. Pick a password, hash it, insert the row yourself.

API env (file `backend/.env` or the process environment). Names match `backend/cmd/api`:

```sh
RPMP_DATABASE_URL=postgres://USER:PASS@127.0.0.1:5432/rpmp
RPMP_JWT_SECRET=at-least-32-characters-long-secret
RPMP_ALLOWED_ORIGIN=http://localhost:3000
RPMP_LOCAL_HTTP=true
# RPMP_ADDR=:8080
```

From `backend/`:

```sh
go run ./cmd/api
```

`GET /` is `404`. Real routes are `/api/v1/...` (login, me, logout, dashboard summary, execution-trend, errors).

Frontend: copy `frontend/.env.example`. Next reads `RPMP_API_ORIGIN`.

```sh
RPMP_API_ORIGIN=http://127.0.0.1:8080
```

From `frontend/`:

```sh
npm install
npm run dev
```

`RPMP_ALLOWED_ORIGIN` must match the browser origin (`http://localhost:3000`), not the API origin.

### Ports

A design-sample dashboard may also bind `:8080`. RPMP HTML is Next on `:3000`. The Go API is `:8080` by default. Do not run the sample and the API on the same port. The sample is layout reference only, not RPMP.

## Commands

From `backend/`:

```sh
go test ./...
go vet ./...
```

Repository integration tests skip unless `DATABASE_URL` is set (that name is test-only, not `RPMP_DATABASE_URL`).

From `frontend/`:

```sh
npm run typecheck
npm run lint
npm test
npm run dev
```

## Workflow

- One backlog ID per session. Frontend cards live in [`docs/backlog/FE/`](docs/backlog/FE/). Backend cards live in [`docs/backlog/BE/`](docs/backlog/BE/).
- Local DB sync (not Swagger): [`docs/plans/LOCAL_DB_SYNC_2026-09-09.md`](docs/plans/LOCAL_DB_SYNC_2026-09-09.md). Next backend card is **BE-05**.
- On an FE card, do not edit `backend/`.
- Execute the card's Execution prompt (agent) or treat in scope / out of scope / frozen JSON / acceptance as the spec (by hand).
- Mark that card complete and open a PR.
- Do not auto-pick the next card.

## FE start here

Read [`docs/onboarding-frontend.md`](docs/onboarding-frontend.md) before changing UI.

Design tokens: [`docs/design/bri-tokens.css`](docs/design/bri-tokens.css).

New FE cards (lead only) start from [`docs/backlog/FE/_TEMPLATE.md`](docs/backlog/FE/_TEMPLATE.md) after the API JSON is frozen.

## Agents

[`AGENTS.md`](AGENTS.md) is for Cursor / Claude. It is not a substitute for this README.
