-- name: UpsertUseCase :one
INSERT INTO use_cases (
    id,
    source_key,
    name,
    status,
    updated_at
) VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (source_key) DO UPDATE SET
    name = EXCLUDED.name,
    status = EXCLUDED.status,
    updated_at = EXCLUDED.updated_at
RETURNING id, source_key, name, status, updated_at;

-- name: InsertExecutions :copyfrom
INSERT INTO executions (
    id,
    use_case_id,
    occurred_at,
    outcome,
    source_ref,
    created_at
) VALUES ($1, $2, $3, $4, $5, $6);

-- name: InsertExecutionErrors :copyfrom
INSERT INTO execution_errors (
    id,
    execution_id,
    use_case_id,
    occurred_at,
    code,
    label,
    created_at
) VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: StartSyncRun :one
INSERT INTO sync_runs (
    id,
    started_at,
    finished_at,
    status,
    rows_read,
    rows_written,
    error_code
) VALUES ($1, $2, NULL, 'running', 0, 0, NULL)
RETURNING id, started_at, finished_at, status, rows_read, rows_written, error_code;

-- name: FinishSyncRun :one
UPDATE sync_runs
SET
    finished_at = $2,
    status = $3,
    rows_read = $4,
    rows_written = $5,
    error_code = $6
WHERE id = $1
RETURNING id, started_at, finished_at, status, rows_read, rows_written, error_code;

-- name: GetUseCaseBySourceKey :one
SELECT id, source_key, name, status, updated_at
FROM use_cases
WHERE source_key = $1;

-- name: GetExecutionByID :one
SELECT id, use_case_id, occurred_at, outcome, source_ref, created_at
FROM executions
WHERE id = $1;

-- name: GetExecutionErrorByID :one
SELECT id, execution_id, use_case_id, occurred_at, code, label, created_at
FROM execution_errors
WHERE id = $1;

-- name: GetSyncRunByID :one
SELECT id, started_at, finished_at, status, rows_read, rows_written, error_code
FROM sync_runs
WHERE id = $1;

-- name: InsertAggregateSnapshot :one
INSERT INTO aggregate_snapshots (
    id,
    sync_run_id,
    source_snapshot_key,
    imported_at,
    source_row_count
) VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (source_snapshot_key) DO NOTHING
RETURNING id, sync_run_id, source_snapshot_key, imported_at, source_row_count;

-- name: InsertProcessAggregateRows :copyfrom
INSERT INTO process_aggregate_rows (
    id,
    aggregate_snapshot_id,
    sync_run_id,
    use_case_id,
    source_process_key,
    process_name,
    package_name,
    environment_name,
    executing_count,
    pending_count,
    suspended_count,
    resumed_count,
    successful_count,
    error_count,
    stopped_count,
    average_duration_seconds,
    average_pending_seconds,
    source_total_rows,
    source_entity_key,
    imported_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
    $11, $12, $13, $14, $15, $16, $17, $18, $19, $20
);

-- name: GetAggregateSnapshotBySourceKey :one
SELECT id, sync_run_id, source_snapshot_key, imported_at, source_row_count
FROM aggregate_snapshots
WHERE source_snapshot_key = $1;

-- name: CountProcessAggregateRowsBySnapshot :one
SELECT count(id)
FROM process_aggregate_rows
WHERE aggregate_snapshot_id = $1;
