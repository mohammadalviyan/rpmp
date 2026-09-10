package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mohammadalviyan/rpmp/backend/internal/domain"
)

type fakeAggregateSource struct {
	snapshot domain.AggregateSnapshot
	err      error
}

func (s fakeAggregateSource) Read(context.Context) (domain.AggregateSnapshot, error) {
	return s.snapshot, s.err
}

type fakeSyncSessions struct {
	session  *fakeSyncSession
	acquired bool
	err      error
}

func (s fakeSyncSessions) TryAcquireSyncSession(context.Context) (domain.SyncSession, bool, error) {
	return s.session, s.acquired, s.err
}

type fakeSyncSession struct {
	started    *domain.SyncRunStart
	published  *domain.AggregateSnapshot
	finished   *domain.SyncRunFinish
	released   bool
	publish    domain.SnapshotPublishResult
	publishErr error
	finishErr  error
}

func (s *fakeSyncSession) StartSyncRun(_ context.Context, start domain.SyncRunStart) (domain.SyncRun, error) {
	s.started = &start
	return domain.SyncRun{ID: start.ID, StartedAt: start.StartedAt, Status: domain.SyncRunRunning}, nil
}

func (s *fakeSyncSession) PublishSnapshot(
	_ context.Context,
	_ string,
	snapshot domain.AggregateSnapshot,
) (domain.SnapshotPublishResult, error) {
	s.published = &snapshot
	return s.publish, s.publishErr
}

func (s *fakeSyncSession) FinishSyncRun(_ context.Context, finish domain.SyncRunFinish) (domain.SyncRun, error) {
	s.finished = &finish
	if s.finishErr != nil {
		return domain.SyncRun{}, s.finishErr
	}
	return domain.SyncRun{
		ID: finish.ID, FinishedAt: &finish.FinishedAt, Status: finish.Status,
		RowsRead: finish.RowsRead, RowsWritten: finish.RowsWritten, ErrorCode: finish.ErrorCode,
	}, nil
}

func (s *fakeSyncSession) Release(context.Context) error {
	s.released = true
	return nil
}

func TestSyncRunnerPublishesAndFinishesSuccess(t *testing.T) {
	session := &fakeSyncSession{publish: domain.SnapshotPublishResult{RowsWritten: 2}}
	snapshot := domain.AggregateSnapshot{Rows: []domain.ProcessAggregate{{}, {}}}
	runner := NewSyncRunner(
		fakeAggregateSource{snapshot: snapshot},
		fakeSyncSessions{session: session, acquired: true},
	)
	now := time.Date(2026, 9, 9, 12, 0, 0, 0, time.UTC)
	runner.now = func() time.Time { return now }
	runner.newID = func() string { return "f33ad9a6-fc61-4ce9-b528-588a40a72f31" }

	result, err := runner.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if result.Run.Status != domain.SyncRunSuccess ||
		result.Run.RowsRead != 2 ||
		result.Run.RowsWritten != 2 ||
		session.finished == nil ||
		session.finished.ErrorCode != nil ||
		!session.released {
		t.Fatalf("unexpected result=%#v session=%#v", result, session)
	}
}

func TestSyncRunnerRecordsSanitizedSourceFailure(t *testing.T) {
	session := &fakeSyncSession{}
	runner := NewSyncRunner(
		fakeAggregateSource{err: errors.New("row 2 contains secret raw details")},
		fakeSyncSessions{session: session, acquired: true},
	)

	if _, err := runner.Run(context.Background()); err == nil {
		t.Fatal("source failure accepted")
	}
	if session.finished == nil ||
		session.finished.Status != domain.SyncRunFailure ||
		session.finished.ErrorCode == nil ||
		*session.finished.ErrorCode != syncErrorSourceRead ||
		!session.released {
		t.Fatalf("failure was not recorded safely: %#v", session)
	}
}

func TestSyncRunnerRollsFailedPublishIntoStableCode(t *testing.T) {
	session := &fakeSyncSession{publishErr: errors.New("database detail")}
	runner := NewSyncRunner(
		fakeAggregateSource{snapshot: domain.AggregateSnapshot{
			Rows: []domain.ProcessAggregate{{}, {}},
		}},
		fakeSyncSessions{session: session, acquired: true},
	)

	if _, err := runner.Run(context.Background()); err == nil {
		t.Fatal("publish failure accepted")
	}
	if session.finished == nil ||
		session.finished.Status != domain.SyncRunFailure ||
		session.finished.RowsRead != 2 ||
		session.finished.RowsWritten != 0 ||
		session.finished.ErrorCode == nil ||
		*session.finished.ErrorCode != syncErrorSnapshotWrite {
		t.Fatalf("publish failure was not recorded safely: %#v", session)
	}
}

func TestSyncRunnerNeverClaimsSuccessOnUnusableSession(t *testing.T) {
	unusable := errors.New("sync session connection is unusable")
	session := &fakeSyncSession{
		publishErr: errors.New("commit failed after request cancellation"),
		finishErr:  unusable,
	}
	runner := NewSyncRunner(
		fakeAggregateSource{snapshot: domain.AggregateSnapshot{
			Rows: []domain.ProcessAggregate{{}},
		}},
		fakeSyncSessions{session: session, acquired: true},
	)

	result, err := runner.Run(context.Background())
	if err == nil {
		t.Fatal("unusable session reported success")
	}
	if !errors.Is(err, unusable) {
		t.Fatalf("error does not report the unusable session: %v", err)
	}
	if result.Run.Status == domain.SyncRunSuccess || result.Idempotent {
		t.Fatalf("result claims progress: %#v", result)
	}
}

func TestSyncRunnerReportsFailedFinishAfterSuccessfulPublish(t *testing.T) {
	unusable := errors.New("sync session connection is unusable")
	session := &fakeSyncSession{
		publish:   domain.SnapshotPublishResult{RowsWritten: 1},
		finishErr: unusable,
	}
	runner := NewSyncRunner(
		fakeAggregateSource{snapshot: domain.AggregateSnapshot{
			Rows: []domain.ProcessAggregate{{}},
		}},
		fakeSyncSessions{session: session, acquired: true},
	)

	result, err := runner.Run(context.Background())
	if err == nil || !errors.Is(err, unusable) {
		t.Fatalf("finish failure was not surfaced: %v", err)
	}
	if result.Run.Status == domain.SyncRunSuccess {
		t.Fatalf("result claims success: %#v", result)
	}
	if !session.released {
		t.Fatal("session was not released")
	}
}

func TestSyncRunnerRejectsOverlap(t *testing.T) {
	runner := NewSyncRunner(
		fakeAggregateSource{},
		fakeSyncSessions{acquired: false},
	)

	_, err := runner.Run(context.Background())
	if domain.ErrorKindOf(err) != domain.KindSyncInProgress {
		t.Fatalf("error kind = %q, want %q", domain.ErrorKindOf(err), domain.KindSyncInProgress)
	}
}

func TestSyncRunnerReportsIdempotentSnapshot(t *testing.T) {
	session := &fakeSyncSession{publish: domain.SnapshotPublishResult{Idempotent: true}}
	runner := NewSyncRunner(
		fakeAggregateSource{snapshot: domain.AggregateSnapshot{Rows: []domain.ProcessAggregate{{}}}},
		fakeSyncSessions{session: session, acquired: true},
	)

	result, err := runner.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !result.Idempotent || result.Run.RowsWritten != 0 {
		t.Fatalf("result = %#v", result)
	}
}
