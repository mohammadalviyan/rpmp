package domain

import (
	"context"
	"errors"
	"time"
)

type Role string

const (
	RoleViewer Role = "viewer"
	RoleAdmin  Role = "admin"
)

func (r Role) CanRead() bool {
	return r == RoleViewer || r == RoleAdmin
}

type User struct {
	ID           string
	EmployeeID   string
	DisplayName  string
	PasswordHash string
	Role         Role
	Active       bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Principal struct {
	UserID string
	Role   Role
}

type AuditEvent struct {
	ID          string
	ActorUserID *string
	Action      string
	Outcome     string
	RequestID   string
	OccurredAt  time.Time
}

type UseCaseStatus string

const (
	UseCaseActive   UseCaseStatus = "active"
	UseCaseInactive UseCaseStatus = "inactive"
)

type ExecutionStatus string

const (
	ExecutionSuccess   ExecutionStatus = "success"
	ExecutionFailure   ExecutionStatus = "failure"
	ExecutionException ExecutionStatus = "exception"
)

func (s ExecutionStatus) IsFailure() bool {
	return s == ExecutionFailure || s == ExecutionException
}

type UseCase struct {
	ID     string
	Status UseCaseStatus
}

type StoredUseCase struct {
	ID        string
	SourceKey string
	Name      string
	Status    UseCaseStatus
	UpdatedAt time.Time
}

type Execution struct {
	ID        string
	UseCaseID string
	StartedAt time.Time
	Status    ExecutionStatus
}

type ExecutionOutcome string

const (
	ExecutionOutcomeSuccess ExecutionOutcome = "success"
	ExecutionOutcomeFailure ExecutionOutcome = "failure"
)

type StoredExecution struct {
	ID         string
	UseCaseID  string
	OccurredAt time.Time
	Outcome    ExecutionOutcome
	SourceRef  string
	CreatedAt  time.Time
}

type ExecutionError struct {
	ID          string
	ExecutionID *string
	UseCaseID   string
	OccurredAt  time.Time
	Code        string
	Label       string
	CreatedAt   time.Time
}

type SyncRunStatus string

const (
	SyncRunRunning SyncRunStatus = "running"
	SyncRunSuccess SyncRunStatus = "success"
	SyncRunFailure SyncRunStatus = "failure"
)

type SyncRun struct {
	ID          string
	StartedAt   time.Time
	FinishedAt  *time.Time
	Status      SyncRunStatus
	RowsRead    int32
	RowsWritten int32
	ErrorCode   *string
}

type SyncRunStart struct {
	ID        string
	StartedAt time.Time
}

type SyncRunFinish struct {
	ID          string
	FinishedAt  time.Time
	Status      SyncRunStatus
	RowsRead    int32
	RowsWritten int32
	ErrorCode   *string
}

type AggregateSnapshot struct {
	SourceSnapshotKey string
	ImportedAt        time.Time
	Rows              []ProcessAggregate
}

type ProcessAggregate struct {
	SourceProcessKey       string
	UseCaseSourceKey       string
	UseCaseName            string
	ProcessName            string
	PackageName            string
	EnvironmentName        *string
	ExecutingCount         int32
	PendingCount           int32
	SuspendedCount         int32
	ResumedCount           int32
	SuccessfulCount        int32
	ErrorCount             int32
	StoppedCount           int32
	AverageDurationSeconds *float64
	AveragePendingSeconds  *float64
	SourceTotalRows        int32
	SourceEntityKey        string
}

func (r ProcessAggregate) FailureCount() int32 {
	return r.ErrorCount + r.StoppedCount
}

type SnapshotPublishResult struct {
	RowsWritten int32
	Idempotent  bool
}

type SyncSession interface {
	StartSyncRun(context.Context, SyncRunStart) (SyncRun, error)
	PublishSnapshot(context.Context, string, AggregateSnapshot) (SnapshotPublishResult, error)
	FinishSyncRun(context.Context, SyncRunFinish) (SyncRun, error)
	Release(context.Context) error
}

type SyncSessionFactory interface {
	TryAcquireSyncSession(context.Context) (SyncSession, bool, error)
}

type FreshnessStatus string

const (
	FreshnessFresh      FreshnessStatus = "fresh"
	FreshnessSyncFailed FreshnessStatus = "sync_failed"
)

type Freshness struct {
	Status                  FreshnessStatus
	LastSuccessfulRefreshAt time.Time
}

type Period struct {
	From     time.Time
	To       time.Time
	Timezone string
}

type SourceSnapshot struct {
	UseCases   []UseCase
	Executions []Execution
	Freshness  Freshness
}

type KPIs struct {
	TotalUseCases    int
	ActiveUseCases   int
	ExecutionVolume  int
	SuccessRate      *float64
	FailedExecutions int
}

type DashboardSummaryAggregate struct {
	TotalUseCases    int
	ActiveUseCases   int
	ExecutionVolume  int
	SuccessfulCount  int
	FailedExecutions int
	Freshness        Freshness
}

type ExecutionTrendAggregate struct {
	Bucket  time.Time
	Success int
	Failure int
}

type DashboardSummary struct {
	Period    Period
	Freshness Freshness
	KPIs      KPIs
}

type ExecutionTrendPoint struct {
	Bucket  time.Time
	Label   string
	Success int
	Failure int
}

type ExecutionTrend struct {
	Period Period
	Points []ExecutionTrendPoint
}

type ErrorGroup struct {
	Code  string
	Label string
	Count int
}

type DashboardErrors struct {
	Period Period
	Groups []ErrorGroup
}

type ErrorKind string

const (
	KindInvalidRequest     ErrorKind = "invalid_request"
	KindInvalidCredentials ErrorKind = "invalid_credentials"
	KindAccountInactive    ErrorKind = "account_inactive"
	KindUnauthenticated    ErrorKind = "unauthenticated"
	KindForbidden          ErrorKind = "forbidden"
	KindInvalidPeriod      ErrorKind = "invalid_period"
	KindSourceUnavailable  ErrorKind = "source_unavailable"
	KindInternal           ErrorKind = "internal_error"
	KindNotFound           ErrorKind = "not_found"
	KindSyncInProgress     ErrorKind = "sync_in_progress"
)

type Error struct {
	Kind  ErrorKind
	Cause error
}

func (e *Error) Error() string {
	return string(e.Kind)
}

func (e *Error) Unwrap() error {
	return e.Cause
}

func NewError(kind ErrorKind, cause error) error {
	return &Error{Kind: kind, Cause: cause}
}

func ErrorKindOf(err error) ErrorKind {
	var domainErr *Error
	if errors.As(err, &domainErr) {
		return domainErr.Kind
	}
	return KindInternal
}
