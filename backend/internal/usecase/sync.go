package usecase

import (
	"context"
	"errors"
	"fmt"
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

func (r *SyncRunner) Run(ctx context.Context) (result SyncResult, err error) {
	session, acquired, err := r.sessions.TryAcquireSyncSession(ctx)
	if err != nil {
		return result, fmt.Errorf("acquire sync session: %w", err)
	}
	if !acquired {
		return result, domain.NewError(domain.KindSyncInProgress, errors.New("another sync holds the advisory lock"))
	}
	defer func() {
		releaseCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if releaseErr := session.Release(releaseCtx); releaseErr != nil {
			err = errors.Join(err, releaseErr)
		}
	}()

	runID := r.newID()
	startedAt := r.now().UTC()
	run, err := session.StartSyncRun(ctx, domain.SyncRunStart{ID: runID, StartedAt: startedAt})
	if err != nil {
		return result, fmt.Errorf("start sync run: %w", err)
	}
	result.Run = run

	snapshot, err := r.source.Read(ctx)
	if err != nil {
		finishErr := r.finishFailure(session, runID, 0, 0, syncErrorSourceRead)
		return result, errors.Join(fmt.Errorf("read aggregate source: %w", err), finishErr)
	}
	rowsRead := int32(len(snapshot.Rows))
	published, err := session.PublishSnapshot(ctx, runID, snapshot)
	if err != nil {
		finishErr := r.finishFailure(session, runID, rowsRead, 0, syncErrorSnapshotWrite)
		return result, errors.Join(fmt.Errorf("publish aggregate snapshot: %w", err), finishErr)
	}
	finishCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	finished, err := session.FinishSyncRun(finishCtx, domain.SyncRunFinish{
		ID: runID, FinishedAt: r.now().UTC(), Status: domain.SyncRunSuccess,
		RowsRead: rowsRead, RowsWritten: published.RowsWritten,
	})
	if err != nil {
		finishErr := r.finishFailure(
			session, runID, rowsRead, published.RowsWritten, syncErrorFinalize,
		)
		return result, errors.Join(fmt.Errorf("finish successful sync run: %w", err), finishErr)
	}
	result.Run = finished
	result.Idempotent = published.Idempotent
	return result, nil
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
