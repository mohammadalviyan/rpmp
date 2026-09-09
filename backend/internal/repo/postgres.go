package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mohammadalviyan/rpmp/backend/internal/domain"
	"github.com/mohammadalviyan/rpmp/backend/internal/repo/dbgen"
)

type Postgres struct {
	queries *dbgen.Queries
}

func NewPostgres(db dbgen.DBTX) *Postgres {
	return &Postgres{queries: dbgen.New(db)}
}

func (r *Postgres) FindByEmployeeID(ctx context.Context, employeeID string) (domain.User, error) {
	user, err := r.queries.GetUserByEmployeeID(ctx, employeeID)
	if err != nil {
		return domain.User{}, mapQueryError(err)
	}
	return mapUser(user)
}

func (r *Postgres) FindByID(ctx context.Context, id string) (domain.User, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return domain.User{}, domain.NewError(domain.KindNotFound, err)
	}
	user, err := r.queries.GetUserByID(ctx, pgtype.UUID{Bytes: parsed, Valid: true})
	if err != nil {
		return domain.User{}, mapQueryError(err)
	}
	return mapUser(user)
}

func (r *Postgres) Record(ctx context.Context, event domain.AuditEvent) error {
	id, err := uuid.Parse(event.ID)
	if err != nil {
		return fmt.Errorf("parse audit event ID: %w", err)
	}
	var actorID pgtype.UUID
	if event.ActorUserID != nil {
		parsed, parseErr := uuid.Parse(*event.ActorUserID)
		if parseErr != nil {
			return fmt.Errorf("parse audit actor ID: %w", parseErr)
		}
		actorID = pgtype.UUID{Bytes: parsed, Valid: true}
	}
	return r.queries.CreateAuditEvent(ctx, dbgen.CreateAuditEventParams{
		ID:          pgtype.UUID{Bytes: id, Valid: true},
		ActorUserID: actorID,
		Action:      event.Action,
		Outcome:     event.Outcome,
		RequestID:   event.RequestID,
		OccurredAt:  pgtype.Timestamptz{Time: event.OccurredAt, Valid: true},
	})
}

func (r *Postgres) UpsertUseCase(ctx context.Context, useCase domain.StoredUseCase) (domain.StoredUseCase, error) {
	id, err := parseUUID(useCase.ID, "use case")
	if err != nil {
		return domain.StoredUseCase{}, err
	}
	row, err := r.queries.UpsertUseCase(ctx, dbgen.UpsertUseCaseParams{
		ID:        id,
		SourceKey: useCase.SourceKey,
		Name:      useCase.Name,
		Status:    string(useCase.Status),
		UpdatedAt: pgtype.Timestamptz{Time: useCase.UpdatedAt, Valid: true},
	})
	if err != nil {
		return domain.StoredUseCase{}, mapQueryError(err)
	}
	return mapStoredUseCase(row)
}

func (r *Postgres) InsertExecutions(ctx context.Context, executions []domain.StoredExecution) (int64, error) {
	params := make([]dbgen.InsertExecutionsParams, len(executions))
	for i, execution := range executions {
		id, err := parseUUID(execution.ID, "execution")
		if err != nil {
			return 0, err
		}
		useCaseID, err := parseUUID(execution.UseCaseID, "execution use case")
		if err != nil {
			return 0, err
		}
		params[i] = dbgen.InsertExecutionsParams{
			ID:         id,
			UseCaseID:  useCaseID,
			OccurredAt: pgtype.Timestamptz{Time: execution.OccurredAt, Valid: true},
			Outcome:    string(execution.Outcome),
			SourceRef:  execution.SourceRef,
			CreatedAt:  pgtype.Timestamptz{Time: execution.CreatedAt, Valid: true},
		}
	}
	return r.queries.InsertExecutions(ctx, params)
}

func (r *Postgres) InsertExecutionErrors(ctx context.Context, executionErrors []domain.ExecutionError) (int64, error) {
	params := make([]dbgen.InsertExecutionErrorsParams, len(executionErrors))
	for i, executionError := range executionErrors {
		id, err := parseUUID(executionError.ID, "execution error")
		if err != nil {
			return 0, err
		}
		useCaseID, err := parseUUID(executionError.UseCaseID, "execution error use case")
		if err != nil {
			return 0, err
		}
		var executionID pgtype.UUID
		if executionError.ExecutionID != nil {
			executionID, err = parseUUID(*executionError.ExecutionID, "execution error execution")
			if err != nil {
				return 0, err
			}
		}
		params[i] = dbgen.InsertExecutionErrorsParams{
			ID:          id,
			ExecutionID: executionID,
			UseCaseID:   useCaseID,
			OccurredAt:  pgtype.Timestamptz{Time: executionError.OccurredAt, Valid: true},
			Code:        executionError.Code,
			Label:       executionError.Label,
			CreatedAt:   pgtype.Timestamptz{Time: executionError.CreatedAt, Valid: true},
		}
	}
	return r.queries.InsertExecutionErrors(ctx, params)
}

func (r *Postgres) StartSyncRun(ctx context.Context, start domain.SyncRunStart) (domain.SyncRun, error) {
	id, err := parseUUID(start.ID, "sync run")
	if err != nil {
		return domain.SyncRun{}, err
	}
	row, err := r.queries.StartSyncRun(ctx, dbgen.StartSyncRunParams{
		ID:        id,
		StartedAt: pgtype.Timestamptz{Time: start.StartedAt, Valid: true},
	})
	if err != nil {
		return domain.SyncRun{}, err
	}
	return mapSyncRun(row)
}

func (r *Postgres) FinishSyncRun(ctx context.Context, finish domain.SyncRunFinish) (domain.SyncRun, error) {
	id, err := parseUUID(finish.ID, "sync run")
	if err != nil {
		return domain.SyncRun{}, err
	}
	var errorCode pgtype.Text
	if finish.ErrorCode != nil {
		errorCode = pgtype.Text{String: *finish.ErrorCode, Valid: true}
	}
	row, err := r.queries.FinishSyncRun(ctx, dbgen.FinishSyncRunParams{
		ID:          id,
		FinishedAt:  pgtype.Timestamptz{Time: finish.FinishedAt, Valid: true},
		Status:      string(finish.Status),
		RowsRead:    finish.RowsRead,
		RowsWritten: finish.RowsWritten,
		ErrorCode:   errorCode,
	})
	if err != nil {
		return domain.SyncRun{}, mapQueryError(err)
	}
	return mapSyncRun(row)
}

func (r *Postgres) FindUseCaseBySourceKey(ctx context.Context, sourceKey string) (domain.StoredUseCase, error) {
	row, err := r.queries.GetUseCaseBySourceKey(ctx, sourceKey)
	if err != nil {
		return domain.StoredUseCase{}, mapQueryError(err)
	}
	return mapStoredUseCase(row)
}

func (r *Postgres) FindExecutionByID(ctx context.Context, id string) (domain.StoredExecution, error) {
	parsed, err := parseUUID(id, "execution")
	if err != nil {
		return domain.StoredExecution{}, err
	}
	row, err := r.queries.GetExecutionByID(ctx, parsed)
	if err != nil {
		return domain.StoredExecution{}, mapQueryError(err)
	}
	return mapStoredExecution(row)
}

func (r *Postgres) FindExecutionErrorByID(ctx context.Context, id string) (domain.ExecutionError, error) {
	parsed, err := parseUUID(id, "execution error")
	if err != nil {
		return domain.ExecutionError{}, err
	}
	row, err := r.queries.GetExecutionErrorByID(ctx, parsed)
	if err != nil {
		return domain.ExecutionError{}, mapQueryError(err)
	}
	return mapExecutionError(row)
}

func (r *Postgres) FindSyncRunByID(ctx context.Context, id string) (domain.SyncRun, error) {
	parsed, err := parseUUID(id, "sync run")
	if err != nil {
		return domain.SyncRun{}, err
	}
	row, err := r.queries.GetSyncRunByID(ctx, parsed)
	if err != nil {
		return domain.SyncRun{}, mapQueryError(err)
	}
	return mapSyncRun(row)
}

func mapQueryError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.NewError(domain.KindNotFound, err)
	}
	return err
}

func mapUser(user dbgen.User) (domain.User, error) {
	if !user.ID.Valid || !user.CreatedAt.Valid || !user.UpdatedAt.Valid {
		return domain.User{}, errors.New("user contains invalid database values")
	}
	return domain.User{
		ID:           uuid.UUID(user.ID.Bytes).String(),
		EmployeeID:   user.EmployeeID,
		DisplayName:  user.DisplayName,
		PasswordHash: user.PasswordHash,
		Role:         domain.Role(user.Role),
		Active:       user.Active,
		CreatedAt:    user.CreatedAt.Time,
		UpdatedAt:    user.UpdatedAt.Time,
	}, nil
}

func parseUUID(value, field string) (pgtype.UUID, error) {
	parsed, err := uuid.Parse(value)
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("parse %s ID: %w", field, err)
	}
	return pgtype.UUID{Bytes: parsed, Valid: true}, nil
}

func mapStoredUseCase(row dbgen.UseCase) (domain.StoredUseCase, error) {
	if !row.ID.Valid || !row.UpdatedAt.Valid {
		return domain.StoredUseCase{}, errors.New("use case contains invalid database values")
	}
	return domain.StoredUseCase{
		ID:        uuid.UUID(row.ID.Bytes).String(),
		SourceKey: row.SourceKey,
		Name:      row.Name,
		Status:    domain.UseCaseStatus(row.Status),
		UpdatedAt: row.UpdatedAt.Time,
	}, nil
}

func mapStoredExecution(row dbgen.Execution) (domain.StoredExecution, error) {
	if !row.ID.Valid || !row.UseCaseID.Valid || !row.OccurredAt.Valid || !row.CreatedAt.Valid {
		return domain.StoredExecution{}, errors.New("execution contains invalid database values")
	}
	return domain.StoredExecution{
		ID:         uuid.UUID(row.ID.Bytes).String(),
		UseCaseID:  uuid.UUID(row.UseCaseID.Bytes).String(),
		OccurredAt: row.OccurredAt.Time,
		Outcome:    domain.ExecutionOutcome(row.Outcome),
		SourceRef:  row.SourceRef,
		CreatedAt:  row.CreatedAt.Time,
	}, nil
}

func mapExecutionError(row dbgen.ExecutionError) (domain.ExecutionError, error) {
	if !row.ID.Valid || !row.UseCaseID.Valid || !row.OccurredAt.Valid || !row.CreatedAt.Valid {
		return domain.ExecutionError{}, errors.New("execution error contains invalid database values")
	}
	var executionID *string
	if row.ExecutionID.Valid {
		value := uuid.UUID(row.ExecutionID.Bytes).String()
		executionID = &value
	}
	return domain.ExecutionError{
		ID:          uuid.UUID(row.ID.Bytes).String(),
		ExecutionID: executionID,
		UseCaseID:   uuid.UUID(row.UseCaseID.Bytes).String(),
		OccurredAt:  row.OccurredAt.Time,
		Code:        row.Code,
		Label:       row.Label,
		CreatedAt:   row.CreatedAt.Time,
	}, nil
}

func mapSyncRun(row dbgen.SyncRun) (domain.SyncRun, error) {
	if !row.ID.Valid || !row.StartedAt.Valid {
		return domain.SyncRun{}, errors.New("sync run contains invalid database values")
	}
	var finishedAt *time.Time
	if row.FinishedAt.Valid {
		value := row.FinishedAt.Time
		finishedAt = &value
	}
	var errorCode *string
	if row.ErrorCode.Valid {
		value := row.ErrorCode.String
		errorCode = &value
	}
	return domain.SyncRun{
		ID:          uuid.UUID(row.ID.Bytes).String(),
		StartedAt:   row.StartedAt.Time,
		FinishedAt:  finishedAt,
		Status:      domain.SyncRunStatus(row.Status),
		RowsRead:    row.RowsRead,
		RowsWritten: row.RowsWritten,
		ErrorCode:   errorCode,
	}, nil
}
