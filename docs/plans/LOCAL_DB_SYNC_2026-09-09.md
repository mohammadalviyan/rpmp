# Local Orchestrator DB sync (not Swagger)

**Date:** 2026-09-09  
**Status:** Ready to implement card-by-card

## Decision

Dashboard HTTP reads **RPMP Postgres only**. A job copies from a **read-only** local/replica Orchestrator-shaped database into RPMP tables. Not real-time. UI shows last refresh. Admin may trigger the same job (202 + lock). Orchestrator Swagger is out of v1.

See `docs/adrs/ADR_BACKEND.md` Decision 3 (updated) and Decision 8.

## Order

1. BE-05 RPMP sync schema (no source DSN required)
2. BE-06 `cmd/sync` (needs table inventory in the card or a `docs/plans` appendix once known)
3. BE-07 dashboard GET reads RPMP repo; stub file stays for tests
4. BE-08 sync status + Admin POST sync
5. FE-12 Overview `api-provider` + freshness copy
6. FE-13 Admin sync button

FE-12 can start after BE-07 if freshness on summary is enough; FE-13 needs BE-08.

## Non-goals

- Query Orchestrator from `GET /dashboard/*`
- Clone vendor table names into RPMP
- Real-time / websocket
- Email/report live
