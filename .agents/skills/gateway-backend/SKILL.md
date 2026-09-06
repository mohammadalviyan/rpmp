---
name: gateway-backend
description: Route RPMP Go API, SQL, source-adapter, and backend architecture work. Use for backend/** or backend API and data-design documentation.
version: 1.0.0
---

# Backend gateway

## When to use

Use for planned `backend/**`, Go API behavior, PostgreSQL/sqlc design, Orchestrator source adapters, and `docs/adrs/ADR_BACKEND.md`.

Do not use for UI-only work.

## Read first

1. `AGENTS.md` hard invariants.
2. `docs/adrs/ADR_BACKEND.md`.
3. Only the PRD FR/persona/NFR IDs in scope.

No `go-best-practices` skill exists yet. Load no substitute. Add one later only when project conventions are real.

## Backend invariants

- Request flow is middleware (authn) → handler → usecase → repo and/or adapter. Packages: `handler`, `middleware`, `usecase`, `domain`, `repo`, `adapter`, `auth`. There is no application-layer `service/` package.
- Middleware is authn only (who, request ID, timeout). It calls `internal/auth`. Login and SSO callbacks are usecases, not middleware bodies.
- Authz (Viewer vs Admin) lives in usecase, not only in middleware.
- Vendor types and status enums stay in `internal/adapter`. The usecase may map already-normalized values into RPMP language. It does not parse vendor enums.
- No `SELECT *`. List columns explicitly. Use explicit `RETURNING` when rows must be returned.
- Prefer the proposed sqlc + pgx path unless an approved enterprise standard changes the ADR. Do not hand-edit generated sqlc files.
- Define typed Go domain errors and map them to HTTP codes. Raw error strings are not the API contract.
- Business metrics, health rules, and KRI formulas belong in Go usecases, not React.
- Auth is internal first and SSO-ready. Secrets stay server-side.
- RPMP is an information layer, not an Orchestrator control plane.

## Required follow-up

After source implementation:

1. Run Go formatting, tests, and vet commands defined by the repo.
2. Spawn `go-reviewer` in fresh context.
3. Fix every FAIL issue and review again.

Hooks are deferred. Reviewer writes `.data/feedback-loop.json` only when the file exists.
