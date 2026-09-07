---
id: FE-XX # <!-- lead: next unused FE-nn -->
title: "" # <!-- lead: short UI outcome -->
status: blocked # blocked | ready | complete. Stay blocked until blocked_by is complete and a human assigns this card.
domain: frontend
blocked_by:
  - BE-XX # <!-- lead: frozen API card IDs and any prior FE IDs -->
fr:
  - FR-x.y # <!-- lead: only the in-scope PRD IDs. Do not dump the PRD. -->
---

# FE-XX: _title_

_Lead: one paragraph. What the UI changes. Name the API path. Point at tokens or layout PNG if this card uses them. State that the design sample is not RPMP._

## In scope

- _Lead: routes, components, and data the FE agent may touch._

## Out of scope

- Backend files and new API fields.
- Orchestrator vendor names and invented KPIs.
- Email, Generate Report, and other sample-only products unless this card explicitly ships them.
- _Lead: anything else this card must not do._

## Frozen contract

_Lead: paste JSON from curl against the live local API after the matching BE card is complete. Do not invent keys._

Path: `_Lead: METHOD /api/v1/..._`

```json
{}
```

## Acceptance

- [ ] UI matches In scope and Frozen contract.
- [ ] No invented KPI or vendor fields.
- [ ] JWT stays httpOnly. Frontend JavaScript never reads it.
- [ ] `npm run typecheck`, `npm run lint`, and `npm test` pass from `frontend/`.
- [ ] No `backend/` edits.

## Execution prompt

```text
Implement backlog card FE-XX only. A human selected this card.

Read AGENTS.md, docs/adrs/ADR_FRONTEND.md, and docs/backlog/FE/FE-XX.md. Load gateway-frontend.

Implement this ID only. Stay inside In scope. Honor Out of scope and Frozen contract. Do not edit backend/. Do not invent API fields, KPIs, or Orchestrator vendor terms. Do not enable Email or Generate Report unless this card says to.

Run npm run typecheck, npm run lint, and npm test from frontend/. Stop after FE-XX. Do not pick the next card.
```
