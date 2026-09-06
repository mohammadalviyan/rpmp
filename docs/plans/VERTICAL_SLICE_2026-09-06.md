# Backend vertical slice 1

**Status:** Approved for backend Slice 1. Human approval received 2026-09-06.

## Goal

Slice 1 proves the RPMP backend shape with local authentication and one protected Dashboard read path. It issues a short-lived access JWT in an httpOnly cookie, exposes the current user, supports logout, and computes FR-2.2 summary values from a fixture-backed `Source`. It does not connect to Orchestrator or claim that fixture metrics are certified KRIs.

## In scope

- Go backend scaffold under `backend/`, but only after this plan receives human approval.
- Local Employee ID and password login.
- Viewer and Admin identities. Authorization remains in usecases.
- Access JWT in an httpOnly cookie. No refresh token.
- Login and logout audit records.
- Protected `me` and Dashboard summary endpoints.
- Stub `Source` adapter with deterministic fixtures.
- PostgreSQL migrations and sqlc queries for slice-owned data.
- Typed domain errors mapped to the shared HTTP error shape.
- Unit, handler, middleware, repository, and integration tests for this slice.

## Out of scope

- Any Next.js or other frontend scaffold.
- A real Orchestrator API, database, read replica, or export client.
- Robot start, stop, schedule, retry, or other control actions.
- Password reset, account administration, SSO/OIDC, refresh tokens, or business-unit authorization.
- Dashboard trends, errors, attention items, Use Cases screens, and health scores.
- Certified KRI definitions, SLA calculations, production reconciliation, or alerting.
- Hooks, a manifest writer, and automatic selection of later backlog cards.

## Module path

The current `origin` is `https://github.com/mohammadalviyan/rpmp.git`, so the future Go module is:

```text
github.com/mohammadalviyan/rpmp/backend
```

If the repository remote changes, derive the new import root from `git remote get-url origin`, remove its transport suffix, append `/backend`, then update `backend/go.mod` and every module import in the same change. Do not invent a replacement organization.

## HTTP contract

All JSON responses use `Content-Type: application/json`. Protected routes return `401` when authentication is absent or invalid. Usecases return `403` when an authenticated role lacks permission.

Errors use one stable shape:

```json
{
  "code": "invalid_credentials",
  "message": "Employee ID or password is incorrect."
}
```

The public message must not contain password, JWT, SQL, stack, or adapter details.

### `POST /api/v1/auth/login`

Request:

```json
{
  "employee_id": "12345678",
  "password": "user supplied secret"
}
```

Success is `200`, sets `rpmp_access`, and returns:

```json
{
  "user": {
    "id": "018f5f71-9cb9-7a61-97e7-d8f33f12b821",
    "employee_id": "12345678",
    "display_name": "Example Viewer",
    "role": "viewer"
  }
}
```

Known errors are `400 invalid_request`, `401 invalid_credentials`, `403 account_inactive`, and `500 internal_error`. Login success and failure are audited. Passwords and JWTs are never logged.

### `POST /api/v1/auth/logout`

Success is `204` with no body. The response expires `rpmp_access` using the same cookie attributes used at login. Logout is audited when a principal can be resolved. Deleting the browser cookie ends that browser session; a copied token remains valid until its short expiry because Slice 1 has no server-side revocation store.

### `GET /api/v1/auth/me`

Success is `200`:

```json
{
  "user": {
    "id": "018f5f71-9cb9-7a61-97e7-d8f33f12b821",
    "employee_id": "12345678",
    "display_name": "Example Viewer",
    "role": "viewer"
  }
}
```

The usecase reloads the user so a disabled account or changed role takes effect without waiting for JWT expiry.

### `GET /api/v1/dashboard/summary`

Optional `from` and `to` query parameters are RFC 3339 timestamps and define a half-open interval `[from, to)`. They must appear together. When omitted, the server uses the trailing 30 days ending at the request time. `from` must be earlier than `to`.

Success is `200`:

```json
{
  "period": {
    "from": "2026-08-07T00:00:00Z",
    "to": "2026-09-06T00:00:00Z",
    "timezone": "UTC"
  },
  "freshness": {
    "status": "fresh",
    "last_successful_refresh_at": "2026-09-06T00:00:00Z"
  },
  "kpis": {
    "total_use_cases": 12,
    "active_use_cases": 9,
    "execution_volume": 100,
    "success_rate": 94.0,
    "failed_executions": 6
  }
}
```

`success_rate` is a JSON number from `0` through `100`, rounded to two decimal places, or `null` when `execution_volume` is zero. Known errors are `400 invalid_period`, `401 unauthenticated`, `403 forbidden`, `503 source_unavailable`, and `500 internal_error`.

## Stub metric definitions

These definitions exist only to make fixtures and tests deterministic. They are not certified KRIs and must be reviewed during Phase 0 before a real source adapter is built.

- **Success, stub:** an execution whose already-normalized fixture status is `success`.
- **Failure, stub:** an execution whose already-normalized fixture status is `failure` or `exception`.
- **Active, stub:** a registered use case whose current RPMP fixture status is `active`.
- **Execution Volume, stub:** all fixture executions whose start time falls in `[from, to)`.
- **Success Rate, stub:** successful executions divided by Execution Volume, multiplied by 100. It is `null` when volume is zero.

The adapter owns fixture parsing and emits RPMP domain values. The usecase computes the summary. Vendor fields and enums do not enter the contract.

## JWT cookie

- Name: `rpmp_access`.
- Token: signed access JWT only. Claims are `iss`, `sub`, `role`, `iat`, `exp`, and `jti`; do not put secrets or profile data in claims.
- TTL: configurable, with a 30-minute default.
- Cookie: `HttpOnly`, host-only, `Path=/`, `SameSite=Lax`, and `Secure` outside explicit local HTTP development.
- Refresh token: no.
- Signing key: supplied by approved server-side secret management. Reject weak or missing production keys.
- State-changing auth routes: validate `Origin` against the configured RPMP origin. Reassess CSRF controls before adding authenticated mutation endpoints.

## Database tables

Only RPMP-owned identity and audit data is persisted in Slice 1. Dashboard fixtures remain files in the stub adapter.

### `users`

- `id uuid primary key`
- `employee_id text not null unique`
- `display_name text not null`
- `password_hash text not null`
- `role text not null` with `viewer` or `admin`
- `active boolean not null default true`
- `created_at timestamptz not null`
- `updated_at timestamptz not null`

### `audit_events`

- `id uuid primary key`
- `actor_user_id uuid null references users(id)`
- `action text not null` with `login` or `logout` for this slice
- `outcome text not null` with `success` or `failure`
- `request_id text not null`
- `occurred_at timestamptz not null`

Failed login can have a null actor. The audit record must not contain the submitted password, a password hash, or a JWT. Every query lists columns explicitly. Writes use explicit `RETURNING` only when a row must be returned. sqlc output is generated, never hand-edited.

There is no `sessions` or refresh-token table. Logout revocation is bounded by the short access-token TTL.

## Planned folder tree

The tree below is a future implementation target, not part of this documentation change.

```text
backend/
  go.mod
  cmd/
    api/
      main.go
  db/
    migrations/
    queries/
  internal/
    adapter/
      source.go
      stub/
        fixtures/
    auth/
    domain/
    handler/
    middleware/
    repo/
    usecase/
  sqlc.yaml
```

`cmd/api` only loads configuration and wires dependencies. Request flow is authentication middleware to handler to usecase to repository and/or adapter. Login is a usecase. Middleware proves identity only. Authorization stays in usecases.

## Test plan

1. **Domain and auth unit tests:** JWT issue/verify, expiry, signature rejection, required claims, role values, typed errors, and password hash comparison.
2. **Usecase unit tests:** valid and invalid login, inactive user, Viewer/Admin reads, disabled user on `me`, audit calls, logout audit, period validation, fixture totals, failure/exception grouping, rounding, and zero-volume `null`.
3. **Middleware and handler tests:** cookie attributes, missing/expired/tampered cookie, safe error JSON, status mapping, malformed requests, Origin rejection, logout expiry, and no secret leakage.
4. **Repository tests:** migrations on PostgreSQL, explicit sqlc queries, user lookup, role/active state, and audit inserts.
5. **Adapter tests:** deterministic fixture loading, normalized RPMP statuses, half-open period boundaries, freshness, and a typed unavailable error.
6. **API integration tests:** login to `me` to summary to logout, unauthorized access after cookie deletion, Viewer and Admin summary access, and source failure as `503`.
7. **Commands:** run formatting, `go test ./...`, `go test -race ./...`, and `go vet ./...`. Add the exact supported commands to `AGENTS.md` when the scaffold exists.

## Implementation order after approval

Each numbered item is a separate future session and must cite its backlog ID. An agent stops after its assigned card and does not select the next card.

1. **BE-01:** scaffold the minimum Go module, migrations, identity repository, auth package, auth usecases, middleware, handlers, routes, and tests.
2. Human reviews BE-01 and marks its dependants ready when its contract is accepted.
3. **BE-02:** add the `Source` port, stub fixtures, Dashboard summary usecase, handler, route, and tests.
4. Human reviews BE-02 and decides whether frontend work or Phase 0 discovery starts next.
5. **FE-01** and then **FE-02** may run only when their `blocked_by` cards are complete and a human selects each card.

## Open questions for a human

1. **Cookie deployment topology:** same origin or a separate API host? Default to a same-origin reverse proxy and a host-only cookie.
2. **Access TTL:** what does Security approve? Default to 30 minutes with no refresh token for Slice 1.
3. **Business timezone:** which timezone controls default periods? Default to UTC until Phase 0 names the business timezone.
4. **Audit retention and sink:** how long and where should auth audit events be retained? Default to 90 days in PostgreSQL, with export deferred until Security decides.
5. **Production user provisioning:** how are initial users created? Default to an explicit administrator-run seed/import command backed by approved secret input, with no built-in production credentials.
