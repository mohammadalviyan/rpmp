package repo

import (
	"context"
	"errors"
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

	for _, name := range []string{
		"000001_auth.up.sql",
		"000002_operational.up.sql",
		"000003_aggregate_snapshots.up.sql",
	} {
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

	dashboardPeriod := domain.Period{
		From: now,
		To:   now.Add(2 * time.Second),
	}
	if _, err := repository.AggregateDashboardSummary(ctx, dashboardPeriod); domain.ErrorKindOf(err) != domain.KindSourceUnavailable {
		t.Fatalf("rows without a successful sync must be unavailable: %v", err)
	}
	if _, err := repository.LatestSyncRun(ctx); domain.ErrorKindOf(err) != domain.KindNotFound {
		t.Fatalf("latest sync without runs = %v", err)
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
	latestSyncRun, err := repository.LatestSyncRun(ctx)
	if err != nil || latestSyncRun.ID != syncRunID || latestSyncRun.Status != domain.SyncRunSuccess {
		t.Fatalf("latest sync run: run=%#v err=%v", latestSyncRun, err)
	}

	dashboardSummary, err := repository.AggregateDashboardSummary(ctx, dashboardPeriod)
	if err != nil {
		t.Fatal(err)
	}
	if dashboardSummary.TotalUseCases != 1 || dashboardSummary.ActiveUseCases != 0 ||
		dashboardSummary.ExecutionVolume != 2 || dashboardSummary.SuccessfulCount != 1 ||
		dashboardSummary.FailedExecutions != 1 ||
		dashboardSummary.Freshness.Status != domain.FreshnessFresh ||
		!dashboardSummary.Freshness.LastSuccessfulRefreshAt.Equal(finishedAt) {
		t.Fatalf("dashboard summary aggregate = %#v", dashboardSummary)
	}
	halfOpenSummary, err := repository.AggregateDashboardSummary(ctx, domain.Period{
		From: now,
		To:   now.Add(time.Second),
	})
	if err != nil {
		t.Fatal(err)
	}
	if halfOpenSummary.ExecutionVolume != 1 || halfOpenSummary.SuccessfulCount != 1 ||
		halfOpenSummary.FailedExecutions != 0 {
		t.Fatalf("dashboard [from,to) aggregate = %#v", halfOpenSummary)
	}
	trend, err := repository.AggregateDashboardExecutionTrend(ctx, dashboardPeriod)
	if err != nil {
		t.Fatal(err)
	}
	if len(trend) != 1 || trend[0].Success != 1 || trend[0].Failure != 1 ||
		trend[0].Bucket.Month() != now.Month() {
		t.Fatalf("dashboard trend aggregate = %#v", trend)
	}
	errorGroups, err := repository.AggregateDashboardErrors(ctx, dashboardPeriod)
	if err != nil {
		t.Fatal(err)
	}
	if len(errorGroups) != 2 ||
		errorGroups[0] != (domain.ErrorGroup{Code: "faulted", Label: "Faulted", Count: 1}) ||
		errorGroups[1] != (domain.ErrorGroup{Code: "stopped", Label: "Stopped", Count: 1}) {
		t.Fatalf("dashboard error aggregates = %#v", errorGroups)
	}

	laterFailureID := uuid.NewString()
	if _, err := repository.StartSyncRun(ctx, domain.SyncRunStart{
		ID: laterFailureID, StartedAt: finishedAt.Add(time.Minute),
	}); err != nil {
		t.Fatal(err)
	}
	failureCode := "source_read_failed"
	if _, err := repository.FinishSyncRun(ctx, domain.SyncRunFinish{
		ID: laterFailureID, FinishedAt: finishedAt.Add(2 * time.Minute), Status: domain.SyncRunFailure,
		ErrorCode: &failureCode,
	}); err != nil {
		t.Fatal(err)
	}
	lastGood, err := repository.AggregateDashboardSummary(ctx, dashboardPeriod)
	if err != nil {
		t.Fatal(err)
	}
	if lastGood.Freshness.Status != domain.FreshnessSyncFailed ||
		!lastGood.Freshness.LastSuccessfulRefreshAt.Equal(finishedAt) ||
		lastGood.ExecutionVolume != 2 {
		t.Fatalf("last-good dashboard aggregate = %#v", lastGood)
	}

	syncRepository := NewSyncPostgres(pool)
	session, acquired, err := syncRepository.TryAcquireSyncSession(ctx)
	if err != nil || !acquired {
		t.Fatalf("acquire sync session: acquired=%t err=%v", acquired, err)
	}
	overlap, overlapAcquired, err := syncRepository.TryAcquireSyncSession(ctx)
	if err != nil || overlapAcquired || overlap != nil {
		t.Fatalf("overlapping session: session=%#v acquired=%t err=%v", overlap, overlapAcquired, err)
	}
	aggregateRunID := uuid.NewString()
	if _, err := session.StartSyncRun(ctx, domain.SyncRunStart{ID: aggregateRunID, StartedAt: now}); err != nil {
		t.Fatal(err)
	}
	environment := "Production"
	duration := 12.5
	snapshot := domain.AggregateSnapshot{
		SourceSnapshotKey: "csv-sha256:integration",
		ImportedAt:        now,
		Rows: []domain.ProcessAggregate{
			{
				SourceProcessKey: "10", UseCaseSourceKey: "Area/Name", UseCaseName: "Area/Name",
				ProcessName: "Error Process", PackageName: "error.pkg",
				SuccessfulCount: 0, ErrorCount: 3, StoppedCount: 0,
				SourceTotalRows: 2, SourceEntityKey: "110",
			},
			{
				SourceProcessKey: "11", UseCaseSourceKey: "Area/Name", UseCaseName: "Area/Name",
				ProcessName: "Stopped Process", PackageName: "stopped.pkg",
				EnvironmentName: &environment, SuccessfulCount: 0, ErrorCount: 0, StoppedCount: 4,
				AverageDurationSeconds: &duration, SourceTotalRows: 2, SourceEntityKey: "111",
			},
		},
	}
	published, err := session.PublishSnapshot(ctx, aggregateRunID, snapshot)
	if err != nil || published.RowsWritten != 2 || published.Idempotent {
		t.Fatalf("publish snapshot: result=%#v err=%v", published, err)
	}
	if _, err := session.FinishSyncRun(ctx, domain.SyncRunFinish{
		ID: aggregateRunID, FinishedAt: now, Status: domain.SyncRunSuccess,
		RowsRead: 2, RowsWritten: 2,
	}); err != nil {
		t.Fatal(err)
	}
	if err := session.Release(ctx); err != nil {
		t.Fatal(err)
	}

	var aggregateRows, groupedUseCases, failures int
	if err := pool.QueryRow(ctx, `
		SELECT
		    count(par.id),
		    count(DISTINCT par.use_case_id),
		    coalesce(sum(par.error_count + par.stopped_count), 0)
		FROM process_aggregate_rows par
		JOIN aggregate_snapshots snapshots ON snapshots.id = par.aggregate_snapshot_id
		WHERE snapshots.source_snapshot_key = $1
	`, snapshot.SourceSnapshotKey).Scan(&aggregateRows, &groupedUseCases, &failures); err != nil {
		t.Fatal(err)
	}
	if aggregateRows != 2 || groupedUseCases != 1 || failures != 7 {
		t.Fatalf("aggregate rows=%d use_cases=%d failures=%d", aggregateRows, groupedUseCases, failures)
	}

	secondSession, acquired, err := syncRepository.TryAcquireSyncSession(ctx)
	if err != nil || !acquired {
		t.Fatalf("reacquire sync session: acquired=%t err=%v", acquired, err)
	}
	secondRunID := uuid.NewString()
	if _, err := secondSession.StartSyncRun(ctx, domain.SyncRunStart{ID: secondRunID, StartedAt: now}); err != nil {
		t.Fatal(err)
	}
	idempotent, err := secondSession.PublishSnapshot(ctx, secondRunID, snapshot)
	if err != nil || !idempotent.Idempotent || idempotent.RowsWritten != 0 {
		t.Fatalf("idempotent publish: result=%#v err=%v", idempotent, err)
	}
	if _, err := secondSession.FinishSyncRun(ctx, domain.SyncRunFinish{
		ID: secondRunID, FinishedAt: now, Status: domain.SyncRunSuccess,
		RowsRead: 2, RowsWritten: 0,
	}); err != nil {
		t.Fatal(err)
	}
	if err := secondSession.Release(ctx); err != nil {
		t.Fatal(err)
	}

	failedSession, acquired, err := syncRepository.TryAcquireSyncSession(ctx)
	if err != nil || !acquired {
		t.Fatalf("acquire session for rollback: acquired=%t err=%v", acquired, err)
	}
	failedRunID := uuid.NewString()
	if _, err := failedSession.StartSyncRun(ctx, domain.SyncRunStart{ID: failedRunID, StartedAt: now}); err != nil {
		t.Fatal(err)
	}
	conflicting := domain.AggregateSnapshot{
		SourceSnapshotKey: "csv-sha256:conflict",
		ImportedAt:        now,
		Rows: []domain.ProcessAggregate{
			{
				SourceProcessKey: "20", UseCaseSourceKey: "Area/Conflict", UseCaseName: "Area/Conflict",
				ProcessName: "First", PackageName: "first.pkg", SourceTotalRows: 2, SourceEntityKey: "120",
			},
			{
				SourceProcessKey: "20", UseCaseSourceKey: "Area/Conflict", UseCaseName: "Area/Conflict",
				ProcessName: "Duplicate", PackageName: "duplicate.pkg", SourceTotalRows: 2, SourceEntityKey: "121",
			},
		},
	}
	if _, err := failedSession.PublishSnapshot(ctx, failedRunID, conflicting); err == nil {
		t.Fatal("duplicate source process key accepted")
	}
	var conflictSnapshots int
	if err := pool.QueryRow(ctx, `
		SELECT count(id)
		FROM aggregate_snapshots
		WHERE source_snapshot_key = $1
	`, conflicting.SourceSnapshotKey).Scan(&conflictSnapshots); err != nil {
		t.Fatal(err)
	}
	if conflictSnapshots != 0 {
		t.Fatalf("failed publish left %d snapshots", conflictSnapshots)
	}
	failureCode = "snapshot_write_failed"
	recorded, err := failedSession.FinishSyncRun(ctx, domain.SyncRunFinish{
		ID: failedRunID, FinishedAt: now, Status: domain.SyncRunFailure,
		RowsRead: 2, RowsWritten: 0, ErrorCode: &failureCode,
	})
	if err != nil || recorded.Status != domain.SyncRunFailure {
		t.Fatalf("record failure after rollback: run=%#v err=%v", recorded, err)
	}

	canceledCtx, cancelRequest := context.WithCancel(ctx)
	tx, err := failedSession.(*syncSession).conn.Begin(canceledCtx)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Exec(canceledCtx, `SELECT 1`); err != nil {
		t.Fatal(err)
	}
	cancelRequest()
	if err := failedSession.(*syncSession).rollback(tx); err != nil {
		t.Fatalf("rollback after request cancellation: %v", err)
	}
	if failedSession.(*syncSession).poisoned {
		t.Fatal("confirmed rollback poisoned a healthy session")
	}
	if _, err := failedSession.FinishSyncRun(ctx, domain.SyncRunFinish{
		ID: failedRunID, FinishedAt: now, Status: domain.SyncRunFailure,
		RowsRead: 2, RowsWritten: 0, ErrorCode: &failureCode,
	}); err != nil {
		t.Fatalf("session unusable after cleanup: %v", err)
	}
	if err := failedSession.Release(ctx); err != nil {
		t.Fatal(err)
	}

	poisonedSession, acquired, err := syncRepository.TryAcquireSyncSession(ctx)
	if err != nil || !acquired {
		t.Fatalf("acquire session for discard: acquired=%t err=%v", acquired, err)
	}
	if err := poisonedSession.(*syncSession).discard(); err != nil {
		t.Fatal(err)
	}
	poisonedRunID := uuid.NewString()
	if _, err := poisonedSession.StartSyncRun(ctx, domain.SyncRunStart{
		ID: poisonedRunID, StartedAt: now,
	}); !errors.Is(err, errSyncSessionUnusable) {
		t.Fatalf("start on discarded session: %v", err)
	}
	if _, err := poisonedSession.PublishSnapshot(ctx, poisonedRunID, snapshot); !errors.Is(err, errSyncSessionUnusable) {
		t.Fatalf("publish on discarded session: %v", err)
	}
	if _, err := poisonedSession.FinishSyncRun(ctx, domain.SyncRunFinish{
		ID: failedRunID, FinishedAt: now, Status: domain.SyncRunSuccess, RowsRead: 2, RowsWritten: 2,
	}); !errors.Is(err, errSyncSessionUnusable) {
		t.Fatalf("finish on discarded session: %v", err)
	}
	if err := poisonedSession.Release(ctx); err != nil {
		t.Fatalf("release discarded session: %v", err)
	}
	var stillFailed string
	if err := pool.QueryRow(ctx, `
		SELECT status
		FROM sync_runs
		WHERE id = $1
	`, failedRunID).Scan(&stillFailed); err != nil {
		t.Fatal(err)
	}
	if stillFailed != string(domain.SyncRunFailure) {
		t.Fatalf("discarded session changed run status to %q", stillFailed)
	}
	afterDiscard, acquired, err := syncRepository.TryAcquireSyncSession(ctx)
	if err != nil || !acquired {
		t.Fatalf("advisory lock not free after discard: acquired=%t err=%v", acquired, err)
	}
	if err := afterDiscard.Release(ctx); err != nil {
		t.Fatal(err)
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
		"use_cases":              5,
		"executions":             6,
		"execution_errors":       7,
		"sync_runs":              7,
		"aggregate_snapshots":    5,
		"process_aggregate_rows": 20,
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
