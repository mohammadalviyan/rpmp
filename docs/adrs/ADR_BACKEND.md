# Backend Architecture Decision Record

**Status**: Proposed
**Date**: 2026-09-06
**Authors**: Development Team
**Version**: 1.3

## Document Information

This ADR records proposed backend decisions for the RPA Performance & Monitoring Platform (RPMP). Stack in PRD §34 is not locked against enterprise standards. Implement after docs-first work, not before.

RPMP is an information layer over Orchestrator / RPA data. It is not a robot controller.

## 1. Executive Summary

### Technology stack overview

- **Language**: Go
- **API**: REST over HTTP (`GET /api/v1/...` as sketched in PRD §18)
- **Database**: PostgreSQL (or the approved relational engine)
- **SQL access**: sqlc generating type-safe queries on top of pgx (proposed)
- **Integration**: v1 copies a read-only source Postgres via `cmd/sync` into RPMP tables (`internal/adapter`). Orchestrator HTTP/Swagger is later. HTTP never queries the source DSN.
- **Auth**: Internal first, SSO-ready. Short-lived access JWT in an httpOnly cookie; no refresh token in Slice 1
- **Layout**: `backend/cmd` + `backend/internal` (folders not created yet)

### Key choices

1. Go HTTP API with a thin router, not a heavyweight framework requirement
2. Layering `handler → usecase → repo`, plus HTTP `middleware` (authn) and `auth/` (how identity works)
3. Adapter boundary for vendor Orchestrator types (anti-corruption layer)
4. sqlc + pgx instead of a general ORM
5. Explicit SQL columns, typed domain errors mapped to HTTP
6. RPMP stores metadata and aggregations; it does not become the system of record for robot execution
7. Style: **pragmatic layered architecture with an anti-corruption layer**. Not full hexagonal. Not Clean Architecture four rings.
8. v1 dashboard reads RPMP Postgres only. Cron/`cmd/sync` plus optional Admin `POST /api/v1/sync` (202). Not real-time.

---

## 2. Architectural Decisions

### Decision 1: HTTP API in Go

#### Context

PRD §34 suggests a Go backend. The UI needs Dashboard and Use Cases JSON, not a generic Orchestrator proxy.

#### Decision

Expose a versioned REST API under `/api/v1`. Handlers live in `internal/handler`. Usecases compute RPMP metrics and enforce authz. Repos talk to PostgreSQL. Do not expose Orchestrator URLs or vendor status codes to the browser.

#### Rationale

- Matches the suggested direction and typical internal Go services
- REST matches the PRD endpoint list and is easy to cache
- A dedicated API contract lets the UI stay on business language (PRD §17)

#### Alternatives considered

**Next.js Route Handlers as the only backend**

- Pros: one deployable
- Cons: mixes RSC with integration credentials, harder adapter tests, weaker split for SSO later

**gRPC to the browser**

- Pros: typed contracts
- Cons: extra client tooling for an internal dashboard; REST is enough for MVP

#### Consequences

- OpenAPI or a checked-in example collection can follow once endpoints stabilize
- CORS and cookie domains must be explicit when frontend and API split hosts

---

### Decision 2: Layering `cmd` / `internal`

#### Context

Agents and humans need a predictable place for HTTP, authn, one-FR actions, SQL, and vendor I/O. A single `service/` package tends to become a god object. Full hexagonal ports for every dependency are more ceremony than MVP needs.

#### Decision

Style: **pragmatic layered architecture with an anti-corruption layer** (`internal/adapter`). Not full hexagonal. Not Clean Architecture four rings. A `Source` port/interface may live in `adapter` or `domain` later. Do not require ports for MVP.

```
backend/
  cmd/api/            # main, listen, wire dependencies only
  internal/
    handler/          # HTTP decode/encode, status codes. No Orchestrator. No SQL.
    middleware/       # authn (who), request ID, timeout. Thin. Calls internal/auth.
    usecase/          # one action per PRD FR (GetDashboard, ListUseCases, Login, …). Authz lives here.
    domain/           # RPMP models + typed errors. No HTTP, no SQL.
    repo/             # PostgreSQL / sqlc. Explicit columns.
    adapter/          # Orchestrator API or read-only source DB. Vendor types stay here.
    auth/             # identity machinery: password, token/session, later OIDC. Not HTTP middleware.
```

Request flow: middleware (authn) → handler → usecase → repo and/or adapter.

Handlers do not call Orchestrator. Usecases do not embed SQL if sqlc owns queries. Middleware does not contain Login or SSO callback business beyond calling `auth` / usecase as wired from `cmd`.

Login is a **usecase**, not logic inside middleware. Middleware only authenticates the request using `internal/auth`.

#### Rationale

- Standard Go `internal/` visibility
- One usecase per PRD FR keeps names and tests aligned with `FR-x.y`
- Avoids a god `Service` that mixes dashboard, list, login, and admin
- Adapter isolation matches PRD §16–17 (vendor terms stay out of the UI contract)
- Reviewers can grep `SELECT` only under `repo/` and sqlc SQL files
- Authn (middleware) stays thin so SSO later does not bloat every request path

#### Alternatives considered

**God `internal/service/` package**

- Rejected: one bag of methods grows across FRs, mixes authz with metrics, and hides which file owns Login vs GetDashboard

**Auth only in HTTP middleware**

- Rejected: Login and SSO callbacks are product actions; Viewer vs Admin authz must not exist only as a route gate. Tests would need a full HTTP stack for password and token logic

**Full hexagonal / four-ring Clean Architecture**

- Rejected for MVP: extra port types and directory rings without a second driver (CLI, worker) that needs them. Adapter as ACL is the isolation that matters now

**Flat `package main` with everything in one module path**

- Rejected: no adapter boundary, leaks vendor names into handlers

#### Consequences

- First vertical slice after docs should add this tree for one dashboard summary endpoint plus Login, not the whole PRD API
- Do not create `backend/` until a human approves `docs/plans/VERTICAL_SLICE_2026-09-06.md` and selects a backend backlog card

---

### Decision 3: Orchestrator vs read-only DB (v1 local sync)

#### Context

PRD §16.2: prefer official Orchestrator API, then approved DB/read replica. Direct production DB only if formally approved. Dashboard must not bind to raw production tables from the frontend. Orchestrator HTTP/Swagger is slow to productize. Slice 1 already serves dashboard JSON from a file stub.

#### Decision

**v1 (local / presentation-to-live):** do **not** call Orchestrator Swagger from `cmd/api`. Copy operational rows from a **read-only** source Postgres (local replica, dump restore, or approved replica) into **RPMP-owned** tables via `cmd/sync` on a schedule. Dashboard usecases read only RPMP Postgres.

Keep **one** `Source` mapping in `internal/adapter` (vendor columns and `JobState` stay there). HTTP never opens the Orchestrator connection string.

Freshness is explicit in the UI (`last_successful_refresh_at`, stale/failed). This is batch data, not real-time. Admin may trigger the same sync job asynchronously (`POST /api/v1/sync` returns 202). Cron remains the default. A manual button does not claim live data.

Orchestrator HTTP API stays a later adapter implementation behind the same domain types, not a blocker for Overview-on-real-rows.

#### Rationale

- Reporting SQL fits history and aggregates better than a awkward Swagger surface
- Request path stays fast and free of vendor credentials
- File stub remains for unit tests when `DATABASE_URL` source is unset
- Matches PRD: adapter boundary, no frontend-to-Orchestrator, batch refresh acceptable for MVP

#### Alternatives considered

**Frontend or handler queries Orchestrator DB**

- Rejected: leaks schema, ties latency to source, violates PRD

**Clone vendor table names into RPMP (`Jobs`, `Releases`)**

- Rejected: UI and usecases would speak vendor; schema churn breaks the app

**Wait for Swagger before any live numbers**

- Rejected: blocks presentation; API can replace the adapter later

**Only cache, no RPMP database**

- Rejected: metadata, comparison, and freshness need a store RPMP controls

#### Consequences

- Need a table inventory (Phase 0 lite) before `cmd/sync` SQL is written. Explicit columns only.
- Production Orchestrator DB still needs formal approval. Local dump/replica is the default for development.
- UI must not label the dashboard "live." Show last sync time and optional Admin sync.
- If API and DB both exist later: API for live robot status if needed; DB/sync for history. Same domain types.

---

### Decision 8: Sync job and manual trigger

#### Context

Copied data lags the source. Operators and demos need a way to refresh without pretending the dashboard is a live tail of Orchestrator.

#### Decision

- `cmd/sync`: CLI used by cron (interval documented, default 5–15 minutes in local compose/docs). Reads source DSN (`RPMP_SOURCE_DATABASE_URL`, read-only). Writes RPMP tables. Inserts `sync_runs`.
- One sync at a time (DB advisory lock or row lock). Overlapping cron/manual is a no-op or 409.
- `POST /api/v1/sync`: Admin only. Starts the same job in-process or exec; returns **202** with run id. Do not block the HTTP handler on a long source query.
- `GET /api/v1/sync/status` or reuse dashboard `freshness`: last run, status, timestamps. Viewer may read status. Only Admin may POST.
- Dashboard `503 source_unavailable` when there is no successful sync and no usable rows, or last run failed and product chooses fail-closed. Prefer showing last good numbers plus `freshness.status = sync_failed` unless empty.

#### Alternatives considered

**Sync inside GET /dashboard/summary**

- Rejected: first page load would hammer the source

**Viewer can press Sync**

- Rejected: easy to stampede the source DB

#### Consequences

- Need `sync_runs` in BE-05. Worker code shared by CLI and POST.
- Frontend Overview shows freshness copy and an Admin-only sync control.

---

### Decision 4: sqlc + pgx (no ORM required)

#### Context

RPMP queries are reporting-shaped: filters, aggregates, explicit columns. An ORM is optional overhead.

#### Decision

**Proposed:** author SQL in `.sql` files, generate Go with **sqlc**, run with **pgx**. No `SELECT *`. Migrations as versioned SQL.

**Noted alternative:** GORM or similar if the team already standardizes on it. If so, still forbid `Select("*")` / implicit all-columns and keep the adapter/repo split.

#### Rationale

- sqlc keeps SQL reviewable and column lists honest
- pgx is the usual Postgres driver for this style
- Matches the invariant in `AGENTS.md`

#### Alternatives considered

**database/sql + handwritten scan**

- Pros: zero codegen
- Cons: more scan bugs; sqlc is the leaner typed option

**GORM as default**

- Pros: familiar to some teams
- Cons: easy `SELECT *`, weaker reporting SQL. Acceptable only if enterprise standard requires it

#### Consequences

- Do not hand-edit sqlc output
- Indexes and retention follow NFR and data-governance sections of the PRD once volume is known

---

### Decision 5: Typed errors mapped to HTTP

#### Context

The UI needs stable error shapes (not found, forbidden, upstream unavailable). Stringly-typed `err.Error()` in JSON drifts.

#### Decision

Define sentinel or structured types in `internal/domain` (for example `UseCaseNotFound`, `SourceUnavailable`). Handlers map them to HTTP status and a small JSON `{ "code", "message" }`. Do not return adapter/vendor error strings as the public `message` without sanitizing.

#### Rationale

- Reviewable contract
- Hides Orchestrator internals from operators using the dashboard

#### Alternatives considered

**`fmt.Errorf` only**

- Rejected as the public API error model

#### Consequences

- Tests should assert codes, not English copy alone
- Logs may keep richer vendor detail server-side

---

### Decision 6: Internal auth with an access JWT cookie

#### Context

PRD requires internal authentication initially and an SSO-ready boundary. Enterprise IdP details remain unknown, but Slice 1 needs one browser authentication mechanism.

#### Decision

Use a signed, short-lived access JWT in an **httpOnly** cookie. Slice 1 does not issue a refresh token and does not store browser tokens in PostgreSQL or frontend JavaScript.

Defaults for Slice 1:

- Cookie name `rpmp_access`
- Configurable 30-minute TTL
- Host-only, `Path=/`, `SameSite=Lax`
- `Secure` outside explicit local HTTP development
- JWT claims limited to `iss`, `sub`, `role`, `iat`, `exp`, and `jti`
- Approved server-side secret management for the signing key
- Allowed-Origin validation on state-changing auth routes

Logout expires the cookie using the same attributes. Without a revocation store, a copied token remains valid until its short expiry. That bounded tradeoff is accepted for Slice 1 and must be reviewed before longer lifetimes or higher-risk mutation APIs are added.

Shape that stays stable:

- Login usecase
- `internal/auth` for JWT issue/verify
- Cookie middleware for authentication on later requests

Handlers stay unchanged when OIDC replaces local passwords: `auth/` and the Login/SSO usecase change; HTTP decode/encode does not.

#### Rationale

- Matches the locked Slice 1 direction while keeping the browser unable to read the token
- A short access-only lifetime limits stolen-token exposure without adding refresh rotation to the first slice
- Splitting how (auth) from when (middleware) and what (usecase) is what makes SSO a swap, not a handler rewrite

#### Alternatives considered

**API keys in localStorage**

- Rejected: XSS and accidental frontend storage of secrets

#### Consequences

- `/api/v1` dashboard routes require auth
- Audit fields (PRD observability/audit) attach to the authenticated principal
- SSO remains later work behind `internal/auth`; this decision does not select an enterprise IdP
- Security still must approve the production TTL, signing-key lifecycle, cookie topology, CSRF standard, and audit retention

---

### Decision 7: Auth package vs HTTP middleware

#### Context

Identity has three different jobs: prove who the caller is on each request (authn), decide what they may do (authz), and implement Login / later SSO (how credentials become a session). Putting all three in middleware makes SSO a middleware rewrite. Hiding Viewer vs Admin only in middleware lets a handler skip the check.

#### Decision

- **`internal/middleware`**: when. Authn of the current request (cookies/tokens), request ID, timeout. Thin. Calls `internal/auth`.
- **`internal/auth`**: how. Password hashing, access-JWT issue/verify, later OIDC client. Not HTTP middleware.
- **`internal/usecase`**: Login (and later SSO callback) as FR-shaped actions. Authz (Viewer vs Admin) lives here, not only as a route wrapper.

`cmd/api` wires the three. Middleware does not implement Login.

#### Rationale

- SSO would bloat middleware if callback and token exchange lived there
- Password and token logic are unit-testable without serving HTTP
- Viewer vs Admin must be enforced in usecase so a forgotten middleware still fails closed at the action

#### Alternatives considered

**Login inside middleware**

- Rejected: Login is an FR action with audit, errors, and cookies to set. It belongs in a usecase the handler calls

**Authz only in middleware**

- Rejected: UI hiding is not authorization; usecase must check role for configuration and metadata writes

#### Consequences

- Public routes (login) skip authn middleware or use a variant that does not require a session
- Protected routes run authn middleware, then handler, then usecase authz

---

## 3. Out of scope (backend MVP)

- Starting or stopping robots
- Writing to Orchestrator as a control plane
- Multi-tenant
- Full KRI automation (later phase)

## 4. Open questions

1. Orchestrator product/version and the **local source table/column inventory** for `cmd/sync` (PRD §33). Blocks BE-06 SQL, not BE-05.
2. Whether sqlc is allowed by enterprise Go standards
3. Sync interval for local cron (default 5–15 minutes until ops pick one)
4. Production source: approved replica vs dump vs (later) HTTP API. Swagger is not v1.

## 5. References

- `docs/PRD.md` §16–18, §33–35
- `docs/plans/LOCAL_DB_SYNC_2026-09-09.md`
- `docs/adrs/ADR_FRONTEND.md`
- `docs/adrs/ADR_AI_ORCHESTRATION.md`
- `AGENTS.md`
