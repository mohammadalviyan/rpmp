package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/mohammadalviyan/rpmp/backend/internal/domain"
	"github.com/mohammadalviyan/rpmp/backend/internal/middleware"
	"github.com/mohammadalviyan/rpmp/backend/internal/usecase"
)

type ListUseCasesUsecase interface {
	Execute(context.Context, usecase.ListUseCasesInput) (domain.UseCaseList, error)
}

type GetUseCaseUsecase interface {
	Execute(context.Context, usecase.GetUseCaseInput) (domain.UseCaseDetail, error)
}

type UseCases struct {
	list ListUseCasesUsecase
	get  GetUseCaseUsecase
}

func NewUseCases(list ListUseCasesUsecase, get GetUseCaseUsecase) *UseCases {
	return &UseCases{list: list, get: get}
}

func (h *UseCases) Mount(mux *http.ServeMux, verifier middleware.TokenVerifier) {
	mux.Handle("GET /api/v1/use-cases", middleware.RequireAuthentication(verifier, http.HandlerFunc(h.List)))
	mux.Handle("GET /api/v1/use-cases/{id}", middleware.RequireAuthentication(verifier, http.HandlerFunc(h.Get)))
}

func (h *UseCases) List(w http.ResponseWriter, r *http.Request) {
	principal, ok := middleware.PrincipalFromContext(r.Context())
	if !ok {
		writeError(w, domain.NewError(domain.KindUnauthenticated, nil))
		return
	}
	result, err := h.list.Execute(r.Context(), usecase.ListUseCasesInput{
		Principal: principal,
		Query:     r.URL.Query().Get("q"),
		Status:    r.URL.Query().Get("status"),
	})
	if err != nil {
		writeError(w, err)
		return
	}
	items := make([]useCaseListJSON, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, listUseCaseResponse(item))
	}
	writeJSON(w, http.StatusOK, struct {
		Freshness freshnessJSON     `json:"freshness"`
		Items     []useCaseListJSON `json:"items"`
	}{
		Freshness: useCaseFreshnessResponse(result.Freshness),
		Items:     items,
	})
}

func (h *UseCases) Get(w http.ResponseWriter, r *http.Request) {
	principal, ok := middleware.PrincipalFromContext(r.Context())
	if !ok {
		writeError(w, domain.NewError(domain.KindUnauthenticated, nil))
		return
	}
	result, err := h.get.Execute(r.Context(), usecase.GetUseCaseInput{
		Principal: principal,
		ID:        r.PathValue("id"),
	})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, struct {
		Freshness freshnessJSON     `json:"freshness"`
		UseCase   useCaseDetailJSON `json:"use_case"`
	}{
		Freshness: useCaseFreshnessResponse(result.Freshness),
		UseCase:   detailUseCaseResponse(result.UseCase),
	})
}

type freshnessJSON struct {
	Status                  domain.FreshnessStatus `json:"status"`
	LastSuccessfulRefreshAt string                 `json:"last_successful_refresh_at,omitempty"`
}

type useCaseListJSON struct {
	ID              string               `json:"id"`
	Name            string               `json:"name"`
	Status          domain.UseCaseStatus `json:"status"`
	ProcessCount    int                  `json:"process_count"`
	SuccessfulCount int                  `json:"successful_count"`
	ErrorCount      int                  `json:"error_count"`
	StoppedCount    int                  `json:"stopped_count"`
	FailedCount     int                  `json:"failed_count"`
	ExecutionVolume int                  `json:"execution_volume"`
	SuccessRate     *float64             `json:"success_rate"`
	Environments    []string             `json:"environments"`
}

type useCaseDetailJSON struct {
	ID              string               `json:"id"`
	Name            string               `json:"name"`
	Status          domain.UseCaseStatus `json:"status"`
	SourceKey       string               `json:"source_key"`
	UpdatedAt       string               `json:"updated_at"`
	ProcessCount    int                  `json:"process_count"`
	SuccessfulCount int                  `json:"successful_count"`
	ErrorCount      int                  `json:"error_count"`
	StoppedCount    int                  `json:"stopped_count"`
	FailedCount     int                  `json:"failed_count"`
	ExecutionVolume int                  `json:"execution_volume"`
	SuccessRate     *float64             `json:"success_rate"`
	Environments    []string             `json:"environments"`
	Processes       []useCaseProcessJSON `json:"processes"`
}

type useCaseProcessJSON struct {
	SourceProcessKey string  `json:"source_process_key"`
	ProcessName      string  `json:"process_name"`
	PackageName      string  `json:"package_name"`
	EnvironmentName  *string `json:"environment_name"`
	SuccessfulCount  int     `json:"successful_count"`
	ErrorCount       int     `json:"error_count"`
	StoppedCount     int     `json:"stopped_count"`
	FailedCount      int     `json:"failed_count"`
	ExecutingCount   int     `json:"executing_count"`
	PendingCount     int     `json:"pending_count"`
	SuspendedCount   int     `json:"suspended_count"`
	ResumedCount     int     `json:"resumed_count"`
}

func useCaseFreshnessResponse(freshness domain.Freshness) freshnessJSON {
	response := freshnessJSON{Status: freshness.Status}
	if !freshness.LastSuccessfulRefreshAt.IsZero() {
		response.LastSuccessfulRefreshAt = freshness.LastSuccessfulRefreshAt.UTC().Format(time.RFC3339)
	}
	return response
}

func listUseCaseResponse(item domain.UseCaseMetrics) useCaseListJSON {
	return useCaseListJSON{
		ID: item.ID, Name: item.Name, Status: item.Status, ProcessCount: item.ProcessCount,
		SuccessfulCount: item.SuccessfulCount, ErrorCount: item.ErrorCount,
		StoppedCount: item.StoppedCount, FailedCount: item.FailedCount,
		ExecutionVolume: item.ExecutionVolume, SuccessRate: item.SuccessRate,
		Environments: item.Environments,
	}
}

func detailUseCaseResponse(item domain.UseCaseMetrics) useCaseDetailJSON {
	processes := make([]useCaseProcessJSON, 0, len(item.Processes))
	for _, process := range item.Processes {
		processes = append(processes, useCaseProcessJSON{
			SourceProcessKey: process.SourceProcessKey, ProcessName: process.ProcessName,
			PackageName: process.PackageName, EnvironmentName: process.EnvironmentName,
			SuccessfulCount: process.SuccessfulCount, ErrorCount: process.ErrorCount,
			StoppedCount: process.StoppedCount, FailedCount: process.ErrorCount + process.StoppedCount,
			ExecutingCount: process.ExecutingCount, PendingCount: process.PendingCount,
			SuspendedCount: process.SuspendedCount, ResumedCount: process.ResumedCount,
		})
	}
	return useCaseDetailJSON{
		ID: item.ID, Name: item.Name, Status: item.Status, SourceKey: item.SourceKey,
		UpdatedAt: item.UpdatedAt.UTC().Format(time.RFC3339), ProcessCount: item.ProcessCount,
		SuccessfulCount: item.SuccessfulCount, ErrorCount: item.ErrorCount,
		StoppedCount: item.StoppedCount, FailedCount: item.FailedCount,
		ExecutionVolume: item.ExecutionVolume, SuccessRate: item.SuccessRate,
		Environments: item.Environments, Processes: processes,
	}
}
