package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/mohammadalviyan/rpmp/backend/internal/domain"
)

const (
	syncErrorSourceRead    = "source_read_failed"
	syncErrorSnapshotWrite = "snapshot_write_failed"
	syncErrorFinalize      = "sync_finalize_failed"
)

type AggregateSource interface {
	Read(context.Context) (domain.AggregateSnapshot, error)
}

type SyncResult struct {
	Run        domain.SyncRun
	Idempotent bool
}

type SyncRunner struct {
	source   AggregateSource
	sessions domain.SyncSessionFactory
	now      func() time.Time
	newID    func() string
}

func NewSyncRunner(source AggregateSource, sessions domain.SyncSessionFactory) *SyncRunner {
	return &SyncRunner{
		source: source, sessions: sessions, now: time.Now, newID: uuid.NewString,
	}
}

func (r *SyncRunner) Run(ctx context.Context) (SyncResult, error) {
	result, session, err := r.start(ctx)
	if err != nil {
		return result, err
	}
	return r.execute(ctx, session, result)
}

func (r *SyncRunner) Start(ctx, workerCtx context.Context) (domain.SyncRun, error) {
	result, session, err := r.start(ctx)
	if err != nil {
		return domain.SyncRun{}, err
	}
	if workerCtx == nil {
		workerCtx = context.Background()
	}
	go func() {
		if _, runErr := r.execute(workerCtx, session, result); runErr != nil {
			slog.Error("background sync failed", "run_id", result.Run.ID)
		}
	}()
	return result.Run, nil
}

func (r *SyncRunner) start(ctx context.Context) (SyncResult, domain.SyncSession, error) {
	var result SyncResult
	session, acquired, err := r.sessions.TryAcquireSyncSession(ctx)
	if err != nil {
		return result, nil, fmt.Errorf("acquire sync session: %w", err)
	}
	if !acquired {
		return result, nil, domain.NewError(domain.KindSyncInProgress, errors.New("another sync holds the advisory lock"))
	}

	runID := r.newID()
	startedAt := r.now().UTC()
	run, err := session.StartSyncRun(ctx, domain.SyncRunStart{ID: runID, StartedAt: startedAt})
	if err != nil {
		releaseCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return result, nil, errors.Join(
			fmt.Errorf("start sync run: %w", err),
			session.Release(releaseCtx),
		)
	}
	result.Run = run
	return result, session, nil
}

func (r *SyncRunner) execute(
	ctx context.Context,
	session domain.SyncSession,
	result SyncResult,
) (completed SyncResult, err error) {
	completed = result
	defer func() {
		releaseCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if releaseErr := session.Release(releaseCtx); releaseErr != nil {
			err = errors.Join(err, releaseErr)
		}
	}()

	snapshot, err := r.source.Read(ctx)
	if err != nil {
		finishErr := r.finishFailure(session, result.Run.ID, 0, 0, syncErrorSourceRead)
		return completed, errors.Join(fmt.Errorf("read aggregate source: %w", err), finishErr)
	}
	rowsRead := int32(len(snapshot.Rows))
	published, err := session.PublishSnapshot(ctx, result.Run.ID, snapshot)
	if err != nil {
		finishErr := r.finishFailure(session, result.Run.ID, rowsRead, 0, syncErrorSnapshotWrite)
		return completed, errors.Join(fmt.Errorf("publish aggregate snapshot: %w", err), finishErr)
	}
	finishCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	finished, err := session.FinishSyncRun(finishCtx, domain.SyncRunFinish{
		ID: result.Run.ID, FinishedAt: r.now().UTC(), Status: domain.SyncRunSuccess,
		RowsRead: rowsRead, RowsWritten: published.RowsWritten,
	})
	if err != nil {
		finishErr := r.finishFailure(
			session, result.Run.ID, rowsRead, published.RowsWritten, syncErrorFinalize,
		)
		return completed, errors.Join(fmt.Errorf("finish successful sync run: %w", err), finishErr)
	}
	completed.Run = finished
	completed.Idempotent = published.Idempotent
	return completed, nil
}

func (r *SyncRunner) finishFailure(
	session domain.SyncSession,
	runID string,
	rowsRead int32,
	rowsWritten int32,
	errorCode string,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := session.FinishSyncRun(ctx, domain.SyncRunFinish{
		ID: runID, FinishedAt: r.now().UTC(), Status: domain.SyncRunFailure,
		RowsRead: rowsRead, RowsWritten: rowsWritten, ErrorCode: &errorCode,
	})
	if err != nil {
		return fmt.Errorf("finish failed sync run: %w", err)
	}
	return nil
}
