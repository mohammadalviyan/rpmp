# Use Cases from local RPMP API (not mock)

**Date:** 2026-09-11  
**Status:** Ready to implement card-by-card

## Goal

`/use-cases` and `/use-cases/[id]` read the local Go API when `RPMP_DATA_MODE=api`. Report, email, and settings stay on mock.

## Data that exists

CSV sync stores **aggregate snapshots**, not execution rows. One RPMP use case is one `use_cases` row keyed by `FullyQualifiedName`. Counts come from the **latest successful** `process_aggregate_rows` for that use case.

Do not invent owner, Attended/Unattended, Running/Warning/Stopped, issues tickets, or a weekly chart. Those are mock-only.

Failure on this slice is `error_count + stopped_count`, same as BE-06. `success_rate` is null when volume is zero. Volume is `successful_count + error_count + stopped_count`.

## Order

1. **BE-09** freeze `GET /api/v1/use-cases` and `GET /api/v1/use-cases/{id}`
2. **FE-14** list and detail consume that JSON in api mode

## Non-goals

- `GET .../executions`, weekly trend series, health, SLA, PATCH metadata
- Orchestrator HTTP, source DSN from the API process
- Mapping vendor `JobState` or mock `UseCase` fields onto the JSON
- Closing BE-06 live Postgres (separate leftover)
