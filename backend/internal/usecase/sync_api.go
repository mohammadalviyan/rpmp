package usecase

import (
	"context"
	"time"

	"github.com/mohammadalviyan/rpmp/backend/internal/domain"
)

type SyncStatusRepository interface {
	LatestSyncRun(context.Context) (domain.SyncRun, error)
}

type SyncStatus struct {
	RunID       *string
	Status      domain.SyncRunStatus
	StartedAt   *time.Time
	FinishedAt  *time.Time
	RowsWritten int32
}

type GetSyncStatus struct {
	runs SyncStatusRepository
}

func NewGetSyncStatus(runs SyncStatusRepository) *GetSyncStatus {
	return &GetSyncStatus{runs: runs}
}

func (u *GetSyncStatus) Execute(ctx context.Context, principal domain.Principal) (SyncStatus, error) {
	if !principal.Role.CanRead() {
		return SyncStatus{}, domain.NewError(domain.KindForbidden, nil)
	}
	run, err := u.runs.LatestSyncRun(ctx)
	if err != nil {
		if domain.ErrorKindOf(err) == domain.KindNotFound {
			return SyncStatus{Status: domain.SyncRunNever}, nil
		}
		return SyncStatus{}, err
	}
	runID := run.ID
	startedAt := run.StartedAt.UTC()
	return SyncStatus{
		RunID:       &runID,
		Status:      run.Status,
		StartedAt:   &startedAt,
		FinishedAt:  run.FinishedAt,
		RowsWritten: run.RowsWritten,
	}, nil
}

type BackgroundSyncRunner interface {
	Start(context.Context, context.Context) (domain.SyncRun, error)
}

type StartSync struct {
	runner    BackgroundSyncRunner
	workerCtx context.Context
}

func NewStartSync(runner BackgroundSyncRunner, workerCtx context.Context) *StartSync {
	return &StartSync{runner: runner, workerCtx: workerCtx}
}

func (u *StartSync) Execute(ctx context.Context, principal domain.Principal) (domain.SyncRun, error) {
	if principal.Role != domain.RoleAdmin {
		return domain.SyncRun{}, domain.NewError(domain.KindForbidden, nil)
	}
	return u.runner.Start(ctx, u.workerCtx)
}
