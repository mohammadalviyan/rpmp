package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mohammadalviyan/rpmp/backend/internal/domain"
	"github.com/mohammadalviyan/rpmp/backend/internal/repo/dbgen"
)

const (
	syncAdvisoryLockID int64 = 72706
	syncCleanupTimeout       = 5 * time.Second
	txStatusIdle       byte  = 'I'
)

// errSyncSessionUnusable rejects work on a session whose connection was
// returned to the pool or destroyed after unconfirmed transaction cleanup.
var errSyncSessionUnusable = errors.New("sync session connection is unusable")

type SyncPostgres struct {
	pool *pgxpool.Pool
}

func NewSyncPostgres(pool *pgxpool.Pool) *SyncPostgres {
	return &SyncPostgres{pool: pool}
}

func (r *SyncPostgres) TryAcquireSyncSession(ctx context.Context) (domain.SyncSession, bool, error) {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return nil, false, fmt.Errorf("acquire sync connection: %w", err)
	}
	var acquired bool
	if err := conn.QueryRow(ctx, "SELECT pg_try_advisory_lock($1)", syncAdvisoryLockID).Scan(&acquired); err != nil {
		conn.Release()
		return nil, false, fmt.Errorf("acquire sync lock: %w", err)
	}
	if !acquired {
		conn.Release()
		return nil, false, nil
	}
	return &syncSession{conn: conn, queries: dbgen.New(conn)}, true, nil
}

type syncSession struct {
	conn     *pgxpool.Conn
	queries  *dbgen.Queries
	done     bool
	poisoned bool
}

func (s *syncSession) usable() error {
	if s.poisoned || s.done {
		return errSyncSessionUnusable
	}
	return nil
}

func (s *syncSession) StartSyncRun(ctx context.Context, start domain.SyncRunStart) (domain.SyncRun, error) {
	if err := s.usable(); err != nil {
		return domain.SyncRun{}, err
	}
	id, err := parseUUID(start.ID, "sync run")
	if err != nil {
		return domain.SyncRun{}, err
	}
	row, err := s.queries.StartSyncRun(ctx, dbgen.StartSyncRunParams{
		ID:        id,
		StartedAt: pgtype.Timestamptz{Time: start.StartedAt, Valid: true},
	})
	if err != nil {
		return domain.SyncRun{}, mapQueryError(err)
	}
	return mapSyncRun(row)
}

func (s *syncSession) FinishSyncRun(ctx context.Context, finish domain.SyncRunFinish) (domain.SyncRun, error) {
	if err := s.usable(); err != nil {
		return domain.SyncRun{}, err
	}
	id, err := parseUUID(finish.ID, "sync run")
	if err != nil {
		return domain.SyncRun{}, err
	}
	var errorCode pgtype.Text
	if finish.ErrorCode != nil {
		errorCode = pgtype.Text{String: *finish.ErrorCode, Valid: true}
	}
	row, err := s.queries.FinishSyncRun(ctx, dbgen.FinishSyncRunParams{
		ID: id, FinishedAt: pgtype.Timestamptz{Time: finish.FinishedAt, Valid: true},
		Status: string(finish.Status), RowsRead: finish.RowsRead,
		RowsWritten: finish.RowsWritten, ErrorCode: errorCode,
	})
	if err != nil {
		return domain.SyncRun{}, mapQueryError(err)
	}
	return mapSyncRun(row)
}

func (s *syncSession) PublishSnapshot(
	ctx context.Context,
	syncRunID string,
	snapshot domain.AggregateSnapshot,
) (result domain.SnapshotPublishResult, err error) {
	if err = s.usable(); err != nil {
		return result, err
	}
	runID, err := parseUUID(syncRunID, "sync run")
	if err != nil {
		return result, err
	}
	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return result, fmt.Errorf("begin snapshot transaction: %w", err)
	}
	committed := false
	defer func() {
		if committed {
			return
		}
		if cleanupErr := s.rollback(tx); cleanupErr != nil {
			result = domain.SnapshotPublishResult{}
			err = errors.Join(err, cleanupErr)
		}
	}()
	queries := s.queries.WithTx(tx)
	snapshotID := uuid.New()
	_, err = queries.InsertAggregateSnapshot(ctx, dbgen.InsertAggregateSnapshotParams{
		ID: pgtype.UUID{Bytes: snapshotID, Valid: true}, SyncRunID: runID,
		SourceSnapshotKey: snapshot.SourceSnapshotKey,
		ImportedAt:        pgtype.Timestamptz{Time: snapshot.ImportedAt, Valid: true},
		SourceRowCount:    int32(len(snapshot.Rows)),
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.SnapshotPublishResult{Idempotent: true}, nil
	}
	if err != nil {
		return result, fmt.Errorf("insert aggregate snapshot: %w", err)
	}

	useCaseIDs := make(map[string]pgtype.UUID)
	for _, row := range snapshot.Rows {
		if _, exists := useCaseIDs[row.UseCaseSourceKey]; exists {
			continue
		}
		useCase, queryErr := queries.UpsertUseCase(ctx, dbgen.UpsertUseCaseParams{
			ID: pgtype.UUID{Bytes: uuid.New(), Valid: true}, SourceKey: row.UseCaseSourceKey,
			Name: row.UseCaseName, Status: string(domain.UseCaseActive),
			UpdatedAt: pgtype.Timestamptz{Time: snapshot.ImportedAt, Valid: true},
		})
		if queryErr != nil {
			return result, fmt.Errorf("upsert use case: %w", queryErr)
		}
		useCaseIDs[row.UseCaseSourceKey] = useCase.ID
	}

	params := make([]dbgen.InsertProcessAggregateRowsParams, len(snapshot.Rows))
	for i, row := range snapshot.Rows {
		params[i] = dbgen.InsertProcessAggregateRowsParams{
			ID:                  pgtype.UUID{Bytes: uuid.New(), Valid: true},
			AggregateSnapshotID: pgtype.UUID{Bytes: snapshotID, Valid: true},
			SyncRunID:           runID, UseCaseID: useCaseIDs[row.UseCaseSourceKey],
			SourceProcessKey: row.SourceProcessKey, ProcessName: row.ProcessName,
			PackageName: row.PackageName, EnvironmentName: textValue(row.EnvironmentName),
			ExecutingCount: row.ExecutingCount, PendingCount: row.PendingCount,
			SuspendedCount: row.SuspendedCount, ResumedCount: row.ResumedCount,
			SuccessfulCount: row.SuccessfulCount, ErrorCount: row.ErrorCount,
			StoppedCount:           row.StoppedCount,
			AverageDurationSeconds: floatValue(row.AverageDurationSeconds),
			AveragePendingSeconds:  floatValue(row.AveragePendingSeconds),
			SourceTotalRows:        row.SourceTotalRows, SourceEntityKey: row.SourceEntityKey,
			ImportedAt: pgtype.Timestamptz{Time: snapshot.ImportedAt, Valid: true},
		}
	}
	written, err := queries.InsertProcessAggregateRows(ctx, params)
	if err != nil {
		return result, fmt.Errorf("insert process aggregate rows: %w", err)
	}
	if err = tx.Commit(ctx); err != nil {
		return result, fmt.Errorf("commit snapshot transaction: %w", err)
	}
	committed = true
	return domain.SnapshotPublishResult{RowsWritten: int32(written)}, nil
}

// rollback ends the snapshot transaction on a context the caller cannot
// cancel. A rollback that cannot be confirmed destroys the connection so no
// later statement runs inside a leftover transaction.
func (s *syncSession) rollback(tx pgx.Tx) error {
	if s.poisoned || s.done {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), syncCleanupTimeout)
	defer cancel()
	if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
		return errors.Join(fmt.Errorf("roll back snapshot transaction: %w", err), s.discard())
	}
	conn := s.conn.Conn()
	if conn.IsClosed() {
		return errors.Join(errors.New("snapshot transaction cleanup closed the connection"), s.discard())
	}
	if status := conn.PgConn().TxStatus(); status != txStatusIdle {
		return errors.Join(
			fmt.Errorf("snapshot transaction cleanup left transaction status %q", string(status)),
			s.discard(),
		)
	}
	return nil
}

// discard removes the connection from the pool instead of returning a
// connection whose transaction state is unknown.
func (s *syncSession) discard() error {
	if s.poisoned || s.done {
		s.poisoned = true
		return nil
	}
	s.poisoned = true
	s.done = true
	ctx, cancel := context.WithTimeout(context.Background(), syncCleanupTimeout)
	defer cancel()
	if err := s.conn.Hijack().Close(ctx); err != nil {
		return fmt.Errorf("close unusable sync connection: %w", err)
	}
	return nil
}

func (s *syncSession) Release(ctx context.Context) error {
	if s.done {
		return nil
	}
	var unlocked bool
	if err := s.conn.QueryRow(ctx, "SELECT pg_advisory_unlock($1)", syncAdvisoryLockID).Scan(&unlocked); err != nil {
		return errors.Join(fmt.Errorf("release sync lock: %w", err), s.discard())
	}
	s.done = true
	s.conn.Release()
	if !unlocked {
		return errors.New("sync advisory lock was not held")
	}
	return nil
}

func textValue(value *string) pgtype.Text {
	if value == nil {
		return pgtype.Text{}
	}
	return pgtype.Text{String: *value, Valid: true}
}

func floatValue(value *float64) pgtype.Float8 {
	if value == nil {
		return pgtype.Float8{}
	}
	return pgtype.Float8{Float64: *value, Valid: true}
}
