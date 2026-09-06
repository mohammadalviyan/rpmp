package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/mohammadalviyan/rpmp/backend/internal/adapter/stub"
	"github.com/mohammadalviyan/rpmp/backend/internal/auth"
	"github.com/mohammadalviyan/rpmp/backend/internal/domain"
	"github.com/mohammadalviyan/rpmp/backend/internal/usecase"
)

type dashboardStub struct {
	summary domain.DashboardSummary
	err     error
}

func (s dashboardStub) Execute(context.Context, usecase.DashboardInput) (domain.DashboardSummary, error) {
	return s.summary, s.err
}

func TestDashboardSummaryContractForViewerAndAdmin(t *testing.T) {
	source, err := stub.New()
	if err != nil {
		t.Fatal(err)
	}
	summary := usecase.NewDashboardSummary(source)
	tokens, err := auth.NewTokens("rpmp", "01234567890123456789012345678901", 30*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	router := NewRouter(
		NewAuth(loginStub{}, currentStub{}, &logoutStub{}, CookieConfig{TTL: 30 * time.Minute, Secure: true}),
		NewDashboard(summary),
		tokens,
		"https://rpmp.example",
	)

	for _, role := range []domain.Role{domain.RoleViewer, domain.RoleAdmin} {
		token, err := tokens.Issue("user-1", role)
		if err != nil {
			t.Fatal(err)
		}
		request := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard/summary?from=2026-08-07T00:00:00Z&to=2026-09-06T00:00:00Z", nil)
		request.AddCookie(&http.Cookie{Name: "rpmp_access", Value: token})
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		if response.Code != http.StatusOK {
			t.Fatalf("role %s status = %d %s", role, response.Code, response.Body.String())
		}
		var body map[string]any
		if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
			t.Fatal(err)
		}
		period := body["period"].(map[string]any)
		if period["from"] != "2026-08-07T00:00:00Z" || period["to"] != "2026-09-06T00:00:00Z" || period["timezone"] != "UTC" {
			t.Fatalf("period = %#v", period)
		}
		freshness := body["freshness"].(map[string]any)
		if freshness["status"] != "fresh" || freshness["last_successful_refresh_at"] != "2026-09-06T00:00:00Z" {
			t.Fatalf("freshness = %#v", freshness)
		}
		kpis := body["kpis"].(map[string]any)
		if kpis["total_use_cases"] != float64(12) || kpis["active_use_cases"] != float64(9) ||
			kpis["execution_volume"] != float64(100) || kpis["success_rate"] != float64(94) ||
			kpis["failed_executions"] != float64(6) {
			t.Fatalf("kpis = %#v", kpis)
		}
	}
}

func TestDashboardSummaryRejectsUnauthenticatedAndInvalidPeriod(t *testing.T) {
	router := NewRouter(
		NewAuth(loginStub{}, currentStub{}, &logoutStub{}, CookieConfig{TTL: 30 * time.Minute}),
		NewDashboard(dashboardStub{}),
		verifierStub{},
		"https://rpmp.example",
	)

	request := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard/summary", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusUnauthorized || !strings.Contains(response.Body.String(), `"code":"unauthenticated"`) {
		t.Fatalf("unauthenticated = %d %s", response.Code, response.Body.String())
	}

	protected := NewRouter(
		NewAuth(loginStub{}, currentStub{}, &logoutStub{}, CookieConfig{TTL: 30 * time.Minute}),
		NewDashboard(dashboardStub{err: domain.NewError(domain.KindInvalidPeriod, nil)}),
		verifierStub{principal: domain.Principal{UserID: "user-1", Role: domain.RoleViewer}},
		"https://rpmp.example",
	)
	request = httptest.NewRequest(http.MethodGet, "/api/v1/dashboard/summary?from=2026-09-06T00:00:00Z", nil)
	request.AddCookie(&http.Cookie{Name: "rpmp_access", Value: "valid"})
	response = httptest.NewRecorder()
	protected.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest ||
		response.Body.String() != "{\"code\":\"invalid_period\",\"message\":\"The requested period is invalid.\"}\n" {
		t.Fatalf("invalid period = %d %q", response.Code, response.Body.String())
	}
}

func TestDashboardSummaryHidesSourceFailureDetails(t *testing.T) {
	router := NewRouter(
		NewAuth(loginStub{}, currentStub{}, &logoutStub{}, CookieConfig{TTL: 30 * time.Minute}),
		NewDashboard(dashboardStub{err: domain.NewError(domain.KindSourceUnavailable, errors.New("JobState timeout"))}),
		verifierStub{principal: domain.Principal{UserID: "user-1", Role: domain.RoleViewer}},
		"https://rpmp.example",
	)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard/summary", nil)
	request.AddCookie(&http.Cookie{Name: "rpmp_access", Value: "valid"})
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d", response.Code)
	}
	body := response.Body.String()
	if !strings.Contains(body, `"code":"source_unavailable"`) || strings.Contains(body, "JobState") || strings.Contains(body, "timeout") {
		t.Fatalf("leaked source error: %s", body)
	}
}

func TestDashboardZeroVolumeEncodesNullSuccessRate(t *testing.T) {
	router := NewRouter(
		NewAuth(loginStub{}, currentStub{}, &logoutStub{}, CookieConfig{TTL: 30 * time.Minute}),
		NewDashboard(dashboardStub{summary: domain.DashboardSummary{
			Period: domain.Period{
				From:     time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
				To:       time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
				Timezone: "UTC",
			},
			Freshness: domain.Freshness{
				Status:                  domain.FreshnessFresh,
				LastSuccessfulRefreshAt: time.Date(2026, 9, 6, 0, 0, 0, 0, time.UTC),
			},
			KPIs: domain.KPIs{TotalUseCases: 12, ActiveUseCases: 9},
		}}),
		verifierStub{principal: domain.Principal{UserID: "user-1", Role: domain.RoleViewer}},
		"https://rpmp.example",
	)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard/summary?from=2026-01-01T00:00:00Z&to=2026-01-02T00:00:00Z", nil)
	request.AddCookie(&http.Cookie{Name: "rpmp_access", Value: "valid"})
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if !strings.Contains(response.Body.String(), `"success_rate":null`) {
		t.Fatalf("body = %s", response.Body.String())
	}
}
