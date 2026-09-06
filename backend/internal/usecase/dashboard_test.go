package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mohammadalviyan/rpmp/backend/internal/adapter/stub"
	"github.com/mohammadalviyan/rpmp/backend/internal/domain"
)

type fakeSource struct {
	snapshot domain.SourceSnapshot
	err      error
}

func (f fakeSource) Snapshot(context.Context) (domain.SourceSnapshot, error) {
	return f.snapshot, f.err
}

func TestDashboardSummaryFromStubFixtures(t *testing.T) {
	source, err := stub.New()
	if err != nil {
		t.Fatal(err)
	}
	uc := NewDashboardSummary(source)
	uc.now = func() time.Time { return time.Date(2026, 9, 6, 0, 0, 0, 0, time.UTC) }

	summary, err := uc.Execute(context.Background(), DashboardInput{
		Principal: domain.Principal{UserID: "user-1", Role: domain.RoleViewer},
	})
	if err != nil {
		t.Fatal(err)
	}
	assertFrozenKPIs(t, summary)
	if summary.Period.Timezone != "UTC" ||
		!summary.Period.From.Equal(time.Date(2026, 8, 7, 0, 0, 0, 0, time.UTC)) ||
		!summary.Period.To.Equal(time.Date(2026, 9, 6, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("default period = %#v", summary.Period)
	}
	if !summary.Freshness.LastSuccessfulRefreshAt.Equal(time.Date(2026, 9, 6, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("freshness used request time: %#v", summary.Freshness)
	}
}

func TestDashboardSummaryAdminAndExplicitPeriod(t *testing.T) {
	source, err := stub.New()
	if err != nil {
		t.Fatal(err)
	}
	uc := NewDashboardSummary(source)
	summary, err := uc.Execute(context.Background(), DashboardInput{
		Principal: domain.Principal{UserID: "admin-1", Role: domain.RoleAdmin},
		From:      "2026-08-07T00:00:00Z",
		To:        "2026-09-06T00:00:00Z",
	})
	if err != nil {
		t.Fatal(err)
	}
	assertFrozenKPIs(t, summary)
}

func TestDashboardPeriodValidation(t *testing.T) {
	uc := NewDashboardSummary(fakeSource{})
	cases := []DashboardInput{
		{Principal: domain.Principal{Role: domain.RoleViewer}, From: "2026-08-07T00:00:00Z"},
		{Principal: domain.Principal{Role: domain.RoleViewer}, To: "2026-09-06T00:00:00Z"},
		{Principal: domain.Principal{Role: domain.RoleViewer}, From: "not-a-time", To: "2026-09-06T00:00:00Z"},
		{Principal: domain.Principal{Role: domain.RoleViewer}, From: "2026-09-06T00:00:00Z", To: "2026-08-07T00:00:00Z"},
		{Principal: domain.Principal{Role: domain.RoleViewer}, From: "2026-09-06T00:00:00Z", To: "2026-09-06T00:00:00Z"},
	}
	for _, input := range cases {
		_, err := uc.Execute(context.Background(), input)
		if domain.ErrorKindOf(err) != domain.KindInvalidPeriod {
			t.Fatalf("input %#v err = %v", input, err)
		}
	}
}

func TestDashboardZeroVolumeNullRateAndFailureGrouping(t *testing.T) {
	from := time.Date(2026, 8, 7, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 9, 6, 0, 0, 0, 0, time.UTC)
	uc := NewDashboardSummary(fakeSource{snapshot: domain.SourceSnapshot{
		UseCases: []domain.UseCase{
			{ID: "a", Status: domain.UseCaseActive},
			{ID: "b", Status: domain.UseCaseInactive},
		},
		Executions: []domain.Execution{
			{ID: "1", StartedAt: from, Status: domain.ExecutionSuccess},
			{ID: "2", StartedAt: from.Add(time.Hour), Status: domain.ExecutionFailure},
			{ID: "3", StartedAt: from.Add(2 * time.Hour), Status: domain.ExecutionException},
			{ID: "4", StartedAt: to, Status: domain.ExecutionSuccess},
			{ID: "5", StartedAt: from.Add(-time.Second), Status: domain.ExecutionFailure},
		},
		Freshness: domain.Freshness{Status: domain.FreshnessFresh, LastSuccessfulRefreshAt: to},
	}})

	summary, err := uc.Execute(context.Background(), DashboardInput{
		Principal: domain.Principal{Role: domain.RoleViewer},
		From:      "2026-08-07T00:00:00Z",
		To:        "2026-09-06T00:00:00Z",
	})
	if err != nil {
		t.Fatal(err)
	}
	if summary.KPIs.TotalUseCases != 2 || summary.KPIs.ActiveUseCases != 1 {
		t.Fatalf("use cases = %#v", summary.KPIs)
	}
	if summary.KPIs.ExecutionVolume != 3 || summary.KPIs.FailedExecutions != 2 || summary.KPIs.SuccessRate == nil || *summary.KPIs.SuccessRate != 33.33 {
		t.Fatalf("period kpis = %#v", summary.KPIs)
	}

	empty, err := uc.Execute(context.Background(), DashboardInput{
		Principal: domain.Principal{Role: domain.RoleViewer},
		From:      "2026-01-01T00:00:00Z",
		To:        "2026-01-02T00:00:00Z",
	})
	if err != nil {
		t.Fatal(err)
	}
	if empty.KPIs.ExecutionVolume != 0 || empty.KPIs.FailedExecutions != 0 || empty.KPIs.SuccessRate != nil {
		t.Fatalf("empty kpis = %#v", empty.KPIs)
	}
}

func TestDashboardForbiddenAndSourceUnavailable(t *testing.T) {
	uc := NewDashboardSummary(fakeSource{err: domain.NewError(domain.KindSourceUnavailable, errors.New("vendor timeout"))})
	_, err := uc.Execute(context.Background(), DashboardInput{
		Principal: domain.Principal{Role: domain.Role("operator")},
	})
	if domain.ErrorKindOf(err) != domain.KindForbidden {
		t.Fatalf("authz err = %v", err)
	}

	_, err = uc.Execute(context.Background(), DashboardInput{
		Principal: domain.Principal{Role: domain.RoleViewer},
	})
	if domain.ErrorKindOf(err) != domain.KindSourceUnavailable {
		t.Fatalf("source err = %v", err)
	}
}

func assertFrozenKPIs(t *testing.T, summary domain.DashboardSummary) {
	t.Helper()
	if summary.KPIs.TotalUseCases != 12 || summary.KPIs.ActiveUseCases != 9 ||
		summary.KPIs.ExecutionVolume != 100 || summary.KPIs.FailedExecutions != 6 ||
		summary.KPIs.SuccessRate == nil || *summary.KPIs.SuccessRate != 94 {
		t.Fatalf("kpis = %#v", summary.KPIs)
	}
}
