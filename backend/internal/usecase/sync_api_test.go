package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mohammadalviyan/rpmp/backend/internal/domain"
)

type syncStatusRepositoryStub struct {
	run domain.SyncRun
	err error
}

func (s syncStatusRepositoryStub) LatestSyncRun(context.Context) (domain.SyncRun, error) {
	return s.run, s.err
}

func TestGetSyncStatusReturnsNeverWhenNoRunExists(t *testing.T) {
	get := NewGetSyncStatus(syncStatusRepositoryStub{
		err: domain.NewError(domain.KindNotFound, errors.New("no rows")),
	})

	status, err := get.Execute(context.Background(), domain.Principal{Role: domain.RoleViewer})
	if err != nil {
		t.Fatal(err)
	}
	if status.Status != domain.SyncRunNever || status.RunID != nil ||
		status.StartedAt != nil || status.FinishedAt != nil || status.RowsWritten != 0 {
		t.Fatalf("status = %#v", status)
	}
}

func TestGetSyncStatusReturnsLatestRunForReader(t *testing.T) {
	finishedAt := time.Date(2026, 9, 9, 0, 1, 0, 0, time.UTC)
	get := NewGetSyncStatus(syncStatusRepositoryStub{run: domain.SyncRun{
		ID:          "018f5f71-9cb9-7a61-97e7-d8f33f12b821",
		Status:      domain.SyncRunSuccess,
		StartedAt:   finishedAt.Add(-time.Minute),
		FinishedAt:  &finishedAt,
		RowsWritten: 100,
	}})

	for _, role := range []domain.Role{domain.RoleViewer, domain.RoleAdmin} {
		status, err := get.Execute(context.Background(), domain.Principal{Role: role})
		if err != nil {
			t.Fatal(err)
		}
		if status.RunID == nil || *status.RunID != "018f5f71-9cb9-7a61-97e7-d8f33f12b821" ||
			status.Status != domain.SyncRunSuccess || status.RowsWritten != 100 {
			t.Fatalf("role %s status = %#v", role, status)
		}
	}
}

type backgroundSyncRunnerStub struct {
	run    domain.SyncRun
	err    error
	called bool
}

func (s *backgroundSyncRunnerStub) Start(context.Context, context.Context) (domain.SyncRun, error) {
	s.called = true
	return s.run, s.err
}

func TestStartSyncRequiresAdmin(t *testing.T) {
	runner := &backgroundSyncRunnerStub{}
	start := NewStartSync(runner, context.Background())

	if _, err := start.Execute(
		context.Background(),
		domain.Principal{Role: domain.RoleViewer},
	); domain.ErrorKindOf(err) != domain.KindForbidden {
		t.Fatalf("viewer error = %v", err)
	}
	if runner.called {
		t.Fatal("viewer started synchronization")
	}

	runner.run = domain.SyncRun{ID: "run-1", Status: domain.SyncRunRunning}
	run, err := start.Execute(context.Background(), domain.Principal{Role: domain.RoleAdmin})
	if err != nil || run.ID != "run-1" || !runner.called {
		t.Fatalf("admin run=%#v err=%v called=%t", run, err, runner.called)
	}
}
