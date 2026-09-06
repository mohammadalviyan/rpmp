package repo

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mohammadalviyan/rpmp/backend/internal/domain"
)

func TestPostgresRepositoryAndMigration(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL is not set")
	}
	ctx := context.Background()
	admin, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer admin.Close()

	schema := "be01_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	if _, err := admin.Exec(ctx, `CREATE SCHEMA "`+schema+`"`); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if _, err := admin.Exec(ctx, `DROP SCHEMA "`+schema+`" CASCADE`); err != nil {
			t.Errorf("drop test schema: %v", err)
		}
	}()

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	config.ConnConfig.RuntimeParams["search_path"] = schema
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	migration, err := os.ReadFile(filepath.Join("..", "..", "db", "migrations", "000001_auth.up.sql"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, string(migration)); err != nil {
		t.Fatalf("apply migration: %v", err)
	}

	userID := uuid.NewString()
	now := time.Now().UTC().Truncate(time.Microsecond)
	_, err = pool.Exec(ctx, `
		INSERT INTO users (id, employee_id, display_name, password_hash, role, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, userID, "12345678", "Example Viewer", "bcrypt-hash", "viewer", true, now, now)
	if err != nil {
		t.Fatal(err)
	}

	repository := NewPostgres(pool)
	byEmployee, err := repository.FindByEmployeeID(ctx, "12345678")
	if err != nil {
		t.Fatal(err)
	}
	if byEmployee.ID != userID || byEmployee.Role != domain.RoleViewer || !byEmployee.Active {
		t.Fatalf("unexpected user: %#v", byEmployee)
	}
	byID, err := repository.FindByID(ctx, userID)
	if err != nil || byID.EmployeeID != "12345678" {
		t.Fatalf("find by ID: user=%#v err=%v", byID, err)
	}
	if _, err := repository.FindByEmployeeID(ctx, "missing"); domain.ErrorKindOf(err) != domain.KindNotFound {
		t.Fatalf("missing user error = %v", err)
	}

	eventID := uuid.NewString()
	if err := repository.Record(ctx, domain.AuditEvent{
		ID: eventID, ActorUserID: &userID, Action: "login", Outcome: "success",
		RequestID: "request-1", OccurredAt: now,
	}); err != nil {
		t.Fatal(err)
	}
	var count int
	if err := pool.QueryRow(ctx, `
		SELECT count(id)
		FROM audit_events
		WHERE id = $1 AND actor_user_id = $2 AND action = $3 AND outcome = $4 AND request_id = $5
	`, eventID, userID, "login", "success", "request-1").Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("audit count = %d", count)
	}

	var columnCount int
	if err := pool.QueryRow(ctx, `
		SELECT count(column_name)
		FROM information_schema.columns
		WHERE table_schema = $1 AND table_name IN ('users', 'audit_events')
	`, schema).Scan(&columnCount); err != nil {
		t.Fatal(err)
	}
	if columnCount != 14 {
		t.Fatal(fmt.Sprintf("migration column count = %d, want 14", columnCount))
	}
}
