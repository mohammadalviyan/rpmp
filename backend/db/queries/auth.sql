-- name: GetUserByEmployeeID :one
SELECT id, employee_id, display_name, password_hash, role, active, created_at, updated_at
FROM users
WHERE employee_id = $1;

-- name: GetUserByID :one
SELECT id, employee_id, display_name, password_hash, role, active, created_at, updated_at
FROM users
WHERE id = $1;

-- name: CreateAuditEvent :exec
INSERT INTO audit_events (
    id,
    actor_user_id,
    action,
    outcome,
    request_id,
    occurred_at
) VALUES ($1, $2, $3, $4, $5, $6);
