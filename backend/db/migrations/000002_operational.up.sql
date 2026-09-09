CREATE TABLE use_cases (
    id uuid PRIMARY KEY,
    source_key text NOT NULL UNIQUE,
    name text NOT NULL,
    status text NOT NULL CHECK (status IN ('active', 'inactive')),
    updated_at timestamptz NOT NULL
);

CREATE TABLE executions (
    id uuid PRIMARY KEY,
    use_case_id uuid NOT NULL REFERENCES use_cases(id),
    occurred_at timestamptz NOT NULL,
    outcome text NOT NULL CHECK (outcome IN ('success', 'failure')),
    source_ref text NOT NULL,
    created_at timestamptz NOT NULL
);

CREATE INDEX executions_occurred_at_idx ON executions (occurred_at);
CREATE INDEX executions_use_case_id_occurred_at_idx ON executions (use_case_id, occurred_at);

CREATE TABLE execution_errors (
    id uuid PRIMARY KEY,
    execution_id uuid NULL REFERENCES executions(id),
    use_case_id uuid NOT NULL REFERENCES use_cases(id),
    occurred_at timestamptz NOT NULL,
    code text NOT NULL,
    label text NOT NULL,
    created_at timestamptz NOT NULL
);

CREATE TABLE sync_runs (
    id uuid PRIMARY KEY,
    started_at timestamptz NOT NULL,
    finished_at timestamptz NULL,
    status text NOT NULL CHECK (status IN ('running', 'success', 'failure')),
    rows_read integer NOT NULL,
    rows_written integer NOT NULL,
    error_code text NULL
);

CREATE INDEX sync_runs_started_at_idx ON sync_runs (started_at DESC);
