package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/mohammadalviyan/rpmp/backend/internal/domain"
	"github.com/mohammadalviyan/rpmp/backend/internal/middleware"
	"github.com/mohammadalviyan/rpmp/backend/internal/usecase"
)

type SyncStatusUsecase interface {
	Execute(context.Context, domain.Principal) (usecase.SyncStatus, error)
}

type StartSyncUsecase interface {
	Execute(context.Context, domain.Principal) (domain.SyncRun, error)
}

type Sync struct {
	status SyncStatusUsecase
	start  StartSyncUsecase
}

func NewSync(status SyncStatusUsecase, start StartSyncUsecase) *Sync {
	return &Sync{status: status, start: start}
}

func (h *Sync) Mount(mux *http.ServeMux, verifier middleware.TokenVerifier, allowedOrigin string) {
	mux.Handle(
		"GET /api/v1/sync/status",
		middleware.RequireAuthentication(verifier, http.HandlerFunc(h.Status)),
	)
	start := middleware.RequireOrigin(allowedOrigin, http.HandlerFunc(h.Start))
	mux.Handle("POST /api/v1/sync", middleware.RequireAuthentication(verifier, start))
}

func (h *Sync) Status(w http.ResponseWriter, r *http.Request) {
	principal, ok := middleware.PrincipalFromContext(r.Context())
	if !ok {
		writeError(w, domain.NewError(domain.KindUnauthenticated, nil))
		return
	}
	status, err := h.status.Execute(r.Context(), principal)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, syncStatusResponse(status))
}

func (h *Sync) Start(w http.ResponseWriter, r *http.Request) {
	principal, ok := middleware.PrincipalFromContext(r.Context())
	if !ok {
		writeError(w, domain.NewError(domain.KindUnauthenticated, nil))
		return
	}
	run, err := h.start.Execute(r.Context(), principal)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusAccepted, map[string]string{
		"run_id": run.ID,
		"status": string(domain.SyncRunRunning),
	})
}

type syncStatusJSON struct {
	RunID       *string              `json:"run_id"`
	Status      domain.SyncRunStatus `json:"status"`
	StartedAt   *string              `json:"started_at"`
	FinishedAt  *string              `json:"finished_at"`
	RowsWritten int32                `json:"rows_written"`
}

func syncStatusResponse(status usecase.SyncStatus) syncStatusJSON {
	return syncStatusJSON{
		RunID:       status.RunID,
		Status:      status.Status,
		StartedAt:   formatOptionalTime(status.StartedAt),
		FinishedAt:  formatOptionalTime(status.FinishedAt),
		RowsWritten: status.RowsWritten,
	}
}

func formatOptionalTime(value *time.Time) *string {
	if value == nil {
		return nil
	}
	formatted := value.UTC().Format(time.RFC3339)
	return &formatted
}
