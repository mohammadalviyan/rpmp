# Local Orchestrator DB sync (not Swagger)

**Date:** 2026-09-09  
**Status:** BE-06 interim CSV mode in progress

## Decision

Dashboard HTTP reads **RPMP Postgres only**. A job copies from a **read-only** local/replica Orchestrator-shaped database into RPMP tables. Not real-time. UI shows last refresh. Admin may trigger the same job (202 + lock). Orchestrator Swagger is out of v1.

See `docs/adrs/ADR_BACKEND.md` Decision 3 (updated) and Decision 8.

## BE-06 interim CSV mode

`RPMP_SOURCE_DATABASE_URL` is currently firewalled. A human approved CSV ingestion as an interim adapter. `cmd/sync` requires `RPMP_DATABASE_URL` and `RPMP_SOURCE_CSV_PATH` in this mode; it does not require the source database URL.

`docs/report/DataSourceOrchestrator.csv` is authoritative. It contains 213 process rows, including zero-activity rows, with this exact inventory:

`ProcessId, ProcessName, PackageName, EnvironmentName, FullyQualifiedName, CountExecuting, CountPending, CountSuspended, CountResumed, CountSuccessful, CountErrors, CountStopped, AverageDurationInSeconds, AveragePendingTimeInSeconds, TotalRows, EntityId`

The adapter groups RPMP use cases by the full `FullyQualifiedName`, marks each active, and retains every process row as aggregate snapshot data. It stores successful, error, and stopped counts separately; compatibility failure count is error plus stopped. It records `imported_at` only because the reporting period is unknown. It does not create synthetic execution rows.

`docs/report/FinalReportGoals.csv` is reference only and has known omissions. The sync does not ingest it.

Live read-only source Postgres remains the target adapter. Its acceptance check stays pending until firewall access is available.

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
