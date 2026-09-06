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

type DashboardExecutionTrendUsecase interface {
	Execute(context.Context, usecase.DashboardInput) (domain.ExecutionTrend, error)
}

type DashboardErrorsUsecase interface {
	Execute(context.Context, usecase.DashboardInput) (domain.DashboardErrors, error)
}

type Dashboard struct {
	summary        DashboardUsecase
	executionTrend DashboardExecutionTrendUsecase
	errors         DashboardErrorsUsecase
}

func NewDashboard(
	summary DashboardUsecase,
	executionTrend DashboardExecutionTrendUsecase,
	errors DashboardErrorsUsecase,
) *Dashboard {
	return &Dashboard{summary: summary, executionTrend: executionTrend, errors: errors}
}

func (h *Dashboard) Mount(mux *http.ServeMux, verifier middleware.TokenVerifier) {
	mux.Handle("GET /api/v1/dashboard/summary", middleware.RequireAuthentication(verifier, http.HandlerFunc(h.Summary)))
	mux.Handle("GET /api/v1/dashboard/execution-trend", middleware.RequireAuthentication(verifier, http.HandlerFunc(h.ExecutionTrend)))
	mux.Handle("GET /api/v1/dashboard/errors", middleware.RequireAuthentication(verifier, http.HandlerFunc(h.Errors)))
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
		"period": periodResponse(summary.Period),
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

func (h *Dashboard) ExecutionTrend(w http.ResponseWriter, r *http.Request) {
	principal, ok := middleware.PrincipalFromContext(r.Context())
	if !ok {
		writeError(w, domain.NewError(domain.KindUnauthenticated, nil))
		return
	}
	query := r.URL.Query()
	trend, err := h.executionTrend.Execute(r.Context(), usecase.DashboardInput{
		Principal: principal,
		From:      query.Get("from"),
		To:        query.Get("to"),
	})
	if err != nil {
		writeError(w, err)
		return
	}

	points := make([]map[string]any, 0, len(trend.Points))
	for _, point := range trend.Points {
		points = append(points, map[string]any{
			"bucket":  point.Bucket.UTC().Format(time.RFC3339),
			"label":   point.Label,
			"success": point.Success,
			"failure": point.Failure,
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"period": periodResponse(trend.Period),
		"points": points,
	})
}

func (h *Dashboard) Errors(w http.ResponseWriter, r *http.Request) {
	principal, ok := middleware.PrincipalFromContext(r.Context())
	if !ok {
		writeError(w, domain.NewError(domain.KindUnauthenticated, nil))
		return
	}
	query := r.URL.Query()
	result, err := h.errors.Execute(r.Context(), usecase.DashboardInput{
		Principal: principal,
		From:      query.Get("from"),
		To:        query.Get("to"),
	})
	if err != nil {
		writeError(w, err)
		return
	}

	groups := make([]map[string]any, 0, len(result.Groups))
	for _, group := range result.Groups {
		groups = append(groups, map[string]any{
			"code":  group.Code,
			"label": group.Label,
			"count": group.Count,
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"period": periodResponse(result.Period),
		"groups": groups,
	})
}

func periodResponse(period domain.Period) map[string]string {
	return map[string]string{
		"from":     period.From.UTC().Format(time.RFC3339),
		"to":       period.To.UTC().Format(time.RFC3339),
		"timezone": period.Timezone,
	}
}
