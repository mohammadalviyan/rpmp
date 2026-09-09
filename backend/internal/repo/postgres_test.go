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

	for _, name := range []string{"000001_auth.up.sql", "000002_operational.up.sql"} {
		migration, err := os.ReadFile(filepath.Join("..", "..", "db", "migrations", name))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := pool.Exec(ctx, string(migration)); err != nil {
			t.Fatalf("apply migration %s: %v", name, err)
		}
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

	useCaseID := uuid.NewString()
	storedUseCase, err := repository.UpsertUseCase(ctx, domain.StoredUseCase{
		ID:        useCaseID,
		SourceKey: "invoice-processing",
		Name:      "Invoice Processing",
		Status:    domain.UseCaseActive,
		UpdatedAt: now,
	})
	if err != nil {
		t.Fatal(err)
	}
	updatedAt := now.Add(time.Minute)
	storedUseCase, err = repository.UpsertUseCase(ctx, domain.StoredUseCase{
		ID:        uuid.NewString(),
		SourceKey: "invoice-processing",
		Name:      "Invoice Processing Updated",
		Status:    domain.UseCaseInactive,
		UpdatedAt: updatedAt,
	})
	if err != nil {
		t.Fatal(err)
	}
	if storedUseCase.ID != useCaseID || storedUseCase.Name != "Invoice Processing Updated" ||
		storedUseCase.Status != domain.UseCaseInactive || !storedUseCase.UpdatedAt.Equal(updatedAt) {
		t.Fatalf("upserted use case = %#v", storedUseCase)
	}
	foundUseCase, err := repository.FindUseCaseBySourceKey(ctx, "invoice-processing")
	if err != nil || foundUseCase.ID != useCaseID {
		t.Fatalf("find use case: useCase=%#v err=%v", foundUseCase, err)
	}

	successExecutionID := uuid.NewString()
	failureExecutionID := uuid.NewString()
	executions := []domain.StoredExecution{
		{
			ID: successExecutionID, UseCaseID: useCaseID, OccurredAt: now,
			Outcome: domain.ExecutionOutcomeSuccess, SourceRef: "source-1", CreatedAt: now,
		},
		{
			ID: failureExecutionID, UseCaseID: useCaseID, OccurredAt: now.Add(time.Second),
			Outcome: domain.ExecutionOutcomeFailure, SourceRef: "source-2", CreatedAt: now,
		},
	}
	inserted, err := repository.InsertExecutions(ctx, executions)
	if err != nil || inserted != 2 {
		t.Fatalf("insert executions: count=%d err=%v", inserted, err)
	}
	foundSuccess, err := repository.FindExecutionByID(ctx, successExecutionID)
	if err != nil || foundSuccess.Outcome != domain.ExecutionOutcomeSuccess {
		t.Fatalf("find success execution: execution=%#v err=%v", foundSuccess, err)
	}
	foundFailure, err := repository.FindExecutionByID(ctx, failureExecutionID)
	if err != nil || foundFailure.Outcome != domain.ExecutionOutcomeFailure {
		t.Fatalf("find failure execution: execution=%#v err=%v", foundFailure, err)
	}

	linkedExecutionID := failureExecutionID
	linkedErrorID := uuid.NewString()
	unlinkedErrorID := uuid.NewString()
	executionErrors := []domain.ExecutionError{
		{
			ID: linkedErrorID, ExecutionID: &linkedExecutionID, UseCaseID: useCaseID,
			OccurredAt: now, Code: "faulted", Label: "Faulted", CreatedAt: now,
		},
		{
			ID: unlinkedErrorID, ExecutionID: nil, UseCaseID: useCaseID,
			OccurredAt: now, Code: "stopped", Label: "Stopped", CreatedAt: now,
		},
	}
	inserted, err = repository.InsertExecutionErrors(ctx, executionErrors)
	if err != nil || inserted != 2 {
		t.Fatalf("insert execution errors: count=%d err=%v", inserted, err)
	}
	foundLinkedError, err := repository.FindExecutionErrorByID(ctx, linkedErrorID)
	if err != nil || foundLinkedError.ExecutionID == nil || *foundLinkedError.ExecutionID != failureExecutionID ||
		foundLinkedError.Code != "faulted" || foundLinkedError.Label != "Faulted" {
		t.Fatalf("find linked execution error: executionError=%#v err=%v", foundLinkedError, err)
	}
	foundUnlinkedError, err := repository.FindExecutionErrorByID(ctx, unlinkedErrorID)
	if err != nil || foundUnlinkedError.ExecutionID != nil || foundUnlinkedError.Code != "stopped" {
		t.Fatalf("find unlinked execution error: executionError=%#v err=%v", foundUnlinkedError, err)
	}

	syncRunID := uuid.NewString()
	started, err := repository.StartSyncRun(ctx, domain.SyncRunStart{ID: syncRunID, StartedAt: now})
	if err != nil || started.Status != domain.SyncRunRunning || started.FinishedAt != nil ||
		started.RowsRead != 0 || started.RowsWritten != 0 || started.ErrorCode != nil {
		t.Fatalf("start sync run: run=%#v err=%v", started, err)
	}
	finishedAt := now.Add(2 * time.Minute)
	finished, err := repository.FinishSyncRun(ctx, domain.SyncRunFinish{
		ID: syncRunID, FinishedAt: finishedAt, Status: domain.SyncRunSuccess,
		RowsRead: 4, RowsWritten: 4,
	})
	if err != nil || finished.FinishedAt == nil || !finished.FinishedAt.Equal(finishedAt) ||
		finished.Status != domain.SyncRunSuccess || finished.RowsRead != 4 || finished.RowsWritten != 4 {
		t.Fatalf("finish sync run: run=%#v err=%v", finished, err)
	}
	foundSyncRun, err := repository.FindSyncRunByID(ctx, syncRunID)
	if err != nil || foundSyncRun.Status != domain.SyncRunSuccess {
		t.Fatalf("find sync run: run=%#v err=%v", foundSyncRun, err)
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

	for table, expected := range map[string]int{
		"use_cases":        5,
		"executions":       6,
		"execution_errors": 7,
		"sync_runs":        7,
	} {
		if err := pool.QueryRow(ctx, `
			SELECT count(column_name)
			FROM information_schema.columns
			WHERE table_schema = $1 AND table_name = $2
		`, schema, table).Scan(&columnCount); err != nil {
			t.Fatal(err)
		}
		if columnCount != expected {
			t.Fatalf("%s column count = %d, want %d", table, columnCount, expected)
		}
	}
}
