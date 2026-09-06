---
name: gateway-frontend
description: Route RPMP Next.js routes, RSC/client data, UI, and frontend architecture work. Use for frontend/** and frontend UI documentation.
version: 1.0.0
---

# Frontend gateway

## When to use

Use for planned `frontend/**`, Next.js App Router work, Dashboard or Use Cases UI, and `docs/adrs/ADR_FRONTEND.md`.

Do not use for Go-only work.

## Read first

1. `AGENTS.md` hard invariants.
2. `docs/adrs/ADR_FRONTEND.md`.
3. Only the PRD FR/persona/NFR IDs in scope.

## Intent table

| Intent | Guidance |
| --- | --- |
| Route/layout/auth gate | App Router skill, if one exists later |
| Dashboard read | Prefer RSC fetch from RPMP API |
| Filters, polling, interactive chart/table | Client Component; add TanStack Query only if cache/polling is justified |
| UI component/form/table | Use shadcn/Tailwind skill only when installed |
| API type | Derive from a locked OpenAPI/Go contract; never guess |

Load Next-related skills only when they exist in `.agents/skills/`. Do not import skills from another repo.

## Frontend invariants

- RSC for initial dashboard reads; Client Components only where interaction needs them.
- No Orchestrator vendor terms or fields in UI copy or UI-owned API types. Use RPMP language.
- Do not duplicate health, normalization, or KRI business logic in React.
- No secrets or service-account credentials in frontend code. Browser auth uses httpOnly cookies.
- No `packages/shared` in v1. Do not guess backend DTOs. Use the agreed API contract.
- Do not hand-edit generated Next or API client files.
- RPMP remains read-only for routine monitoring. Do not add robot control actions.

## Required follow-up

After source implementation:

1. Run frontend type, lint, and test commands defined by the repo.
2. Spawn `frontend-reviewer` in fresh context.
3. Fix every FAIL issue and review again.

Hooks are deferred. Reviewer writes `.data/feedback-loop.json` only when the file exists.
