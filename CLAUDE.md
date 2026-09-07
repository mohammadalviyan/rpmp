# CLAUDE.md

Same map as `AGENTS.md`. Canonical copy lives in both files so Claude Code does not need an extra hop.

RPA Performance & Monitoring Platform (RPMP). Internal information layer over Orchestrator / RPA operational data. MVP: Dashboard + Use Cases. Not a replacement for Orchestrator.

Docs-first setup is complete. A human approved `docs/plans/VERTICAL_SLICE_2026-09-06.md` on 2026-09-06. Backend exists through BE-04 and frontend through FE-04. Further work is one human-selected backlog card at a time. Suggested stack (PRD §34, not locked beyond implemented slices): Next.js/React frontend, Go HTTP API, PostgreSQL, internal auth designed SSO-ready. No `packages/shared` in v1.

## Tech stack (proposed)

- **Frontend**: Next.js App Router, React, Tailwind CSS, shadcn/ui (internal admin UI)
- **Backend**: Go HTTP API (`cmd` + `internal`)
- **Database**: PostgreSQL (sqlc + pgx; see `docs/adrs/ADR_BACKEND.md`)
- **Integration**: Orchestrator API or approved read-only DB, behind a backend adapter
- **Auth**: Internal first, SSO-ready. Access JWT in an httpOnly cookie; no refresh token in Slice 1. Secrets stay on the server.

## Project structure

The backend and frontend scaffolds exist.

**Backend** (`backend/`): request flow `middleware (authn) → handler → usecase → repo` and/or `adapter`.

- `cmd/api/`: main, listen, wire dependencies only
- `internal/handler/`: HTTP decode/encode, status codes. No Orchestrator. No SQL.
- `internal/middleware/`: authn (who), request ID, timeout. Thin. Calls `internal/auth`.
- `internal/usecase/`: one action per PRD FR (GetDashboard, ListUseCases, Login). Authz lives here.
- `internal/domain/`: RPMP models and typed errors. No HTTP, no SQL.
- `internal/repo/`: PostgreSQL / sqlc. Explicit columns.
- `internal/adapter/`: Orchestrator API or read-only source DB. Vendor types stay here.
- `internal/auth/`: identity machinery (password, token/session, later OIDC). Not HTTP middleware.

Login is a usecase. Middleware only authenticates the request using `internal/auth`.

**Frontend** (`frontend/`): Next.js App Router:

- `app/`: routes (Dashboard, Use Cases). RSC for dashboard reads; client components where interactivity needs them
- `components/`: UI (`ui/` for shadcn primitives, feature components elsewhere)
- `lib/`: typed clients for the RPMP API only (not Orchestrator)

**Docs and agent layout (exist now):**

- `docs/PRD.md`: product requirements
- `docs/adrs/`: architecture decisions
- `docs/plans/`: implementation plans (`<NAME>_<DATE>.md`)
- `.agents/skills/`: skill source of truth
- `.claude/` / `.cursor/`: tool veneer (symlinks + thin agents / Cursor rules)
- `.data/`: per-session orchestration state (gitignored except `SCHEMA.md`)

## Documentation

- Product requirements: `docs/PRD.md` (cite `FR-x.y`). Draft dump preserved in `docs/PRD.draft-v1.md`.
- `docs/adrs/ADR_BACKEND.md`: Go API, PostgreSQL, adapters, errors, auth
- `docs/adrs/ADR_FRONTEND.md`: Next.js, RSC/client split, UI copy rules
- `docs/adrs/ADR_AI_ORCHESTRATION.md`: agent workflow (scaled-down v1)
- Skills in `.agents/skills/` for process. Do not paste the full PRD into prompts.

## AI orchestration (v1)

Shared knowledge lives in `.agents/skills/`. Claude Code and Cursor symlink into that tree.

- **`orchestrate`**: short protocol (phases 1, 2, 3 intent, 4 gateway, 8 implement, 13 test, 16 done). Auto-invoke for file-changing work only. Not for read-only questions.
- **Gateways**: `gateway-backend` and `gateway-frontend`. No shared package lane.
- **Agents (when source folders exist)**: `backend-developer`, `frontend-developer`, `go-reviewer`, `frontend-reviewer`. Thin, one-shot, no further spawn.
- **Approval gate:** Slice 1 is approved. `backend/` and `frontend/` exist through BE-04 / FE-04. Execute only the human-selected backlog card.
- **Backlog sessions:** code work cites exactly one backlog ID. One agent session executes one ID, stops, and never auto-picks the next card.
- **MEDIUM+**: offer `grill-me` before discovery. Load the matching gateway before implementation.
- **Hooks**: deferred until folders and agents are in use. `manifest-writer` deferred. See ADR_AI_ORCHESTRATION.

## Hard invariants

1. **No `SELECT *`.** Explicit column lists. Explicit `RETURNING` when a write must return rows.
2. **Typed domain errors in Go**, mapped to HTTP. No raw strings as the API error contract.
3. **Orchestrator vendor terms stay in backend adapters.** The UI and the API it consumes use RPMP language (use case, execution, health, failure). See PRD §17.
4. **No business logic duplicated in React** that belongs in Go (normalization, KRI formulas, status mapping).
5. **Auth is internal first, SSO-ready.** No secrets, tokens, or service-account credentials in the frontend.
6. **Do not hand-edit generated files** (sqlc output, Next generated types, lockfile noise).
7. **Agent layout:** product docs in `docs/`; skills in `.agents/skills/`. Point to PRD sections; do not dump the full PRD into prompts.
8. **One backlog card per session.** Do not expand `backend/` on a frontend card or `frontend/` on a backend card. Do not invent or auto-pick the next card.
9. **RPMP is an information layer**, not a robot controller. No start/stop jobs, no Orchestrator replacement features in MVP.
10. **Snake_case DB columns** in SQL and repo structs that map to tables. API JSON for the UI may use the contract in `docs/PRD.md` once it is locked.

## Git workflow

- Never push directly to `master` / `main`. Use a branch.
- Prefixes: `feature/`, `fix/`, `refactor/`, `docs/`, `chore/`.
- Non-trivial plans go in `docs/plans/<NAME>_<DATE>.md`.

## Dev commands

Run from `backend/`:

```sh
go fmt ./...
go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0 generate
go test ./...
go test -race ./...
go vet ./...
```

Repository integration tests run when `DATABASE_URL` points to a disposable PostgreSQL database. They skip when it is unset.

Run from `frontend/`:

```sh
npm run typecheck
npm run lint
npm test
```
