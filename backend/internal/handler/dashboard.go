package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/mohammadalviyan/rpmp/backend/internal/domain"
	"github.com/mohammadalviyan/rpmp/backend/internal/middleware"
	"github.com/mohammadalviyan/rpmp/backend/internal/usecase"
)

type DashboardUsecase interface {
	Execute(context.Context, usecase.DashboardInput) (domain.DashboardSummary, error)
}

type Dashboard struct {
	summary DashboardUsecase
}

func NewDashboard(summary DashboardUsecase) *Dashboard {
	return &Dashboard{summary: summary}
}

func (h *Dashboard) Mount(mux *http.ServeMux, verifier middleware.TokenVerifier) {
	mux.Handle("GET /api/v1/dashboard/summary", middleware.RequireAuthentication(verifier, http.HandlerFunc(h.Summary)))
}

func (h *Dashboard) Summary(w http.ResponseWriter, r *http.Request) {
	principal, ok := middleware.PrincipalFromContext(r.Context())
	if !ok {
		writeError(w, domain.NewError(domain.KindUnauthenticated, nil))
		return
	}
	query := r.URL.Query()
	summary, err := h.summary.Execute(r.Context(), usecase.DashboardInput{
		Principal: principal,
		From:      query.Get("from"),
		To:        query.Get("to"),
	})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"period": map[string]string{
			"from":     summary.Period.From.UTC().Format(time.RFC3339),
			"to":       summary.Period.To.UTC().Format(time.RFC3339),
			"timezone": summary.Period.Timezone,
		},
		"freshness": map[string]string{
			"status":                     string(summary.Freshness.Status),
			"last_successful_refresh_at": summary.Freshness.LastSuccessfulRefreshAt.UTC().Format(time.RFC3339),
		},
		"kpis": map[string]any{
			"total_use_cases":   summary.KPIs.TotalUseCases,
			"active_use_cases":  summary.KPIs.ActiveUseCases,
			"execution_volume":  summary.KPIs.ExecutionVolume,
			"success_rate":      summary.KPIs.SuccessRate,
			"failed_executions": summary.KPIs.FailedExecutions,
		},
	})
}
