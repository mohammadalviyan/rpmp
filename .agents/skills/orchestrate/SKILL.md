---
name: orchestrate
description: Run the short RPMP workflow for work that changes files. Auto-invoke for features, fixes, refactors, code scaffolding, and documentation edits. Do not invoke for read-only questions, explanations, search-only tasks, or pure conversation.
version: 1.0.0
---

# RPMP orchestration v1

See `docs/adrs/ADR_AI_ORCHESTRATION.md`.

## Triage

| Class | Criteria | Path |
| --- | --- | --- |
| TRIVIAL | Typo or one-line config/doc fix | 1, 8, 16 |
| SMALL | Well-understood, one domain, limited files | 1, 2, 3, 4, 8, 13, 16 |
| MEDIUM | Multi-file or some product/design uncertainty | Same as SMALL, with `grill-me` gate |
| LARGE | New subsystem or cross-domain architecture | Plan with user, then run both gateways |

For MEDIUM, ask whether to run `grill-me` before Phase 3. Run it when an in-scope FR or NFR is TBD. LARGE work requires explicit planning; do not silently choose architecture.

## Short phase protocol

1. **Setup**: read `AGENTS.md`; identify requested files and branch state. Read only the in-scope PRD FR/persona/NFR IDs.
2. **Triage**: classify size and domains.
3. **Intent**: restate outcome, constraints, acceptance criteria, and unresolved assumptions.
4. **Gateway**: load `gateway-backend`, `gateway-frontend`, or both. Backend first for an API-to-UI vertical slice.
8. **Implement**: delegate source work to the matching thin developer when available. Do not pass the full PRD or session history.
13. **Test**: run the smallest relevant checks, then the domain suite. Fix failures in scope.
16. **Done**: review the diff, obtain the matching reviewer verdict for source changes, and report paths, tests, assumptions, and blockers.

## Docs-first guard

Docs-first setup is complete, but the repo remains documentation-only until a human approves `docs/plans/VERTICAL_SLICE_2026-09-06.md`. Before that approval, Phase 8 may write only `docs/**`, ADRs, `.agents/skills/**`, and tool veneer/config files. Approval permits backend Slice 1 work only. `frontend/` remains deferred to its backlog cards.

## Backlog execution

- Every code task must cite one backlog ID under `docs/backlog/`.
- One future agent session executes one backlog ID.
- Confirm every `blocked_by` ID is complete and require a human to select the card.
- Stop after the selected card. Never auto-pick or begin the next card.
- Do not create `backend/` before Slice 1 plan approval. Do not create `frontend/` until a human selects an unblocked frontend card.

## Domain routing

- Backend: Go API, SQL, source adapters, backend ADR or API design. Load `gateway-backend`; use `backend-developer`, then `go-reviewer`.
- Frontend: Next routes, RSC/client components, UI, frontend ADR. Load `gateway-frontend`; use `frontend-developer`, then `frontend-reviewer`.
- Cross-domain: settle the API contract first. Run backend, review, then frontend, review.

There is no shared-package lane in RPMP v1.

## State

`.data/SCHEMA.md` defines optional state. `manifest-writer` and enforcement hooks are deferred. Do not make state-file availability a completion requirement in v1.

## Hard rules

- Never invent Orchestrator vendor fields for the frontend contract.
- Never paste the whole PRD into an agent prompt.
- Never spawn developer and reviewer in parallel.
- Never hand-edit generated files.
- Never mark tests PASS unless they ran.
