CREATE TABLE users (
    id uuid PRIMARY KEY,
    employee_id text NOT NULL UNIQUE,
    display_name text NOT NULL,
    password_hash text NOT NULL,
    role text NOT NULL CHECK (role IN ('viewer', 'admin')),
    active boolean NOT NULL DEFAULT true,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);

CREATE TABLE audit_events (
    id uuid PRIMARY KEY,
    actor_user_id uuid NULL REFERENCES users(id),
    action text NOT NULL CHECK (action IN ('login', 'logout')),
    outcome text NOT NULL CHECK (outcome IN ('success', 'failure')),
    request_id text NOT NULL,
    occurred_at timestamptz NOT NULL
);

CREATE INDEX audit_events_actor_user_id_idx ON audit_events (actor_user_id);
CREATE INDEX audit_events_occurred_at_idx ON audit_events (occurred_at);
