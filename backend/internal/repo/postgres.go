package repo

import (
	"context"
	"errors"
	"fmt"

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
