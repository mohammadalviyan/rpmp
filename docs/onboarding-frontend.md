# Frontend first day

Audience: one frontend developer. You may use Cursor Agent or write code by hand. The lead generates FE cards. You do not invent cards from the PRD.

## First-day path

1. Read this file, [`README.md`](../README.md), [`docs/adrs/ADR_FRONTEND.md`](adrs/ADR_FRONTEND.md), [`docs/design/bri-tokens.css`](design/bri-tokens.css), and **only** the assigned FE card. Do not read the full PRD. If you use an agent, optionally load the `gateway-frontend` skill.
2. Run the stack as in the README. Confirm login and the existing dashboard in the browser (`http://localhost:3000`).
3. Open exactly one `docs/backlog/FE/<ID>.md` that a human marked `ready` and assigned to you.
4. If using Cursor: new Agent chat, workspace this repo, paste that card's Execution prompt unchanged. If coding by hand: treat In scope, Out of scope, Frozen contract, and Acceptance as the spec.
5. Do not edit `backend/`. Do not invent KPI or Orchestrator vendor fields. Do not enable Email or Generate Report from the design sample.
6. Run frontend checks, open a PR, set `status: complete` on **that card only**. Stop.

## Cards

The lead writes FE-05+ from [`docs/backlog/FE/_TEMPLATE.md`](backlog/FE/_TEMPLATE.md) after the API JSON is frozen (curl or equivalent). Frontend does not add cards from the PRD.

If `blocked_by` IDs are not `complete`, stop. Do not implement a blocked card.

## Sample vs RPMP

A design sample dashboard on `http://localhost:8080` is layout reference only. It is not RPMP. RPMP HTML is Next on `:3000`. The Go API is `:8080` by default. `GET /` on the API is `404`. Real routes are `/api/v1/...`.
