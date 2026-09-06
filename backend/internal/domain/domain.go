package domain

import (
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

type Execution struct {
	ID        string
	UseCaseID string
	StartedAt time.Time
	Status    ExecutionStatus
}

type FreshnessStatus string

const FreshnessFresh FreshnessStatus = "fresh"

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

type DashboardSummary struct {
	Period    Period
	Freshness Freshness
	KPIs      KPIs
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
