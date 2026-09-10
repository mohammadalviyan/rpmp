CREATE TABLE aggregate_snapshots (
    id uuid PRIMARY KEY,
    sync_run_id uuid NOT NULL REFERENCES sync_runs(id),
    source_snapshot_key text NOT NULL UNIQUE,
    imported_at timestamptz NOT NULL,
    source_row_count integer NOT NULL CHECK (source_row_count >= 0),
    UNIQUE (id, sync_run_id)
);

CREATE INDEX aggregate_snapshots_imported_at_idx ON aggregate_snapshots (imported_at DESC);
CREATE INDEX aggregate_snapshots_sync_run_id_idx ON aggregate_snapshots (sync_run_id);

CREATE TABLE process_aggregate_rows (
    id uuid PRIMARY KEY,
    aggregate_snapshot_id uuid NOT NULL,
    sync_run_id uuid NOT NULL REFERENCES sync_runs(id),
    use_case_id uuid NOT NULL REFERENCES use_cases(id),
    source_process_key text NOT NULL,
    process_name text NOT NULL,
    package_name text NOT NULL,
    environment_name text NULL,
    executing_count integer NOT NULL CHECK (executing_count >= 0),
    pending_count integer NOT NULL CHECK (pending_count >= 0),
    suspended_count integer NOT NULL CHECK (suspended_count >= 0),
    resumed_count integer NOT NULL CHECK (resumed_count >= 0),
    successful_count integer NOT NULL CHECK (successful_count >= 0),
    error_count integer NOT NULL CHECK (error_count >= 0),
    stopped_count integer NOT NULL CHECK (stopped_count >= 0),
    average_duration_seconds double precision NULL CHECK (average_duration_seconds >= 0),
    average_pending_seconds double precision NULL CHECK (average_pending_seconds >= 0),
    source_total_rows integer NOT NULL CHECK (source_total_rows >= 0),
    source_entity_key text NOT NULL,
    imported_at timestamptz NOT NULL,
    FOREIGN KEY (aggregate_snapshot_id, sync_run_id)
        REFERENCES aggregate_snapshots(id, sync_run_id) ON DELETE CASCADE,
    UNIQUE (aggregate_snapshot_id, source_process_key),
    UNIQUE (aggregate_snapshot_id, source_entity_key)
);

CREATE INDEX process_aggregate_rows_use_case_id_idx
    ON process_aggregate_rows (use_case_id);
CREATE INDEX process_aggregate_rows_sync_run_id_idx
    ON process_aggregate_rows (sync_run_id);
