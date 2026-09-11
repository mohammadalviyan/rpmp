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

type fakeDashboardRepository struct {
	summary domain.DashboardSummaryAggregate
	trend   []domain.ExecutionTrendAggregate
	groups  []domain.ErrorGroup
	err     error
}

func (f fakeDashboardRepository) AggregateDashboardSummary(
	context.Context,
	domain.Period,
) (domain.DashboardSummaryAggregate, error) {
	return f.summary, f.err
}

func (f fakeDashboardRepository) AggregateDashboardExecutionTrend(
	context.Context,
	domain.Period,
) ([]domain.ExecutionTrendAggregate, error) {
	return f.trend, f.err
}

func (f fakeDashboardRepository) AggregateDashboardErrors(
	context.Context,
	domain.Period,
) ([]domain.ErrorGroup, error) {
	return f.groups, f.err
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

func TestDashboardExecutionTrendFromStubFixtures(t *testing.T) {
	source, err := stub.New()
	if err != nil {
		t.Fatal(err)
	}
	trend, err := NewDashboardExecutionTrend(source).Execute(context.Background(), DashboardInput{
		Principal: domain.Principal{UserID: "viewer-1", Role: domain.RoleViewer},
		From:      "2026-08-07T00:00:00Z",
		To:        "2026-09-06T00:00:00Z",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(trend.Points) != 2 {
		t.Fatalf("points = %#v", trend.Points)
	}
	august := trend.Points[0]
	if !august.Bucket.Equal(time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)) ||
		august.Label != "Aug" || august.Success != 94 || august.Failure != 6 {
		t.Fatalf("august point = %#v", august)
	}
	september := trend.Points[1]
	if !september.Bucket.Equal(time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)) ||
		september.Label != "Sep" || september.Success != 0 || september.Failure != 0 {
		t.Fatalf("september point = %#v", september)
	}
}

func TestDashboardErrorsFromStubFixtures(t *testing.T) {
	source, err := stub.New()
	if err != nil {
		t.Fatal(err)
	}
	result, err := NewDashboardErrors(source).Execute(context.Background(), DashboardInput{
		Principal: domain.Principal{UserID: "admin-1", Role: domain.RoleAdmin},
		From:      "2026-08-07T00:00:00Z",
		To:        "2026-09-06T00:00:00Z",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Groups) != 2 ||
		result.Groups[0] != (domain.ErrorGroup{Code: "faulted", Label: "Faulted", Count: 4}) ||
		result.Groups[1] != (domain.ErrorGroup{Code: "stopped", Label: "Stopped", Count: 2}) {
		t.Fatalf("groups = %#v", result.Groups)
	}
}

func TestDashboardChartPeriodRules(t *testing.T) {
	now := time.Date(2026, 9, 6, 0, 0, 0, 0, time.UTC)
	trend := NewDashboardExecutionTrend(fakeSource{})
	trend.now = func() time.Time { return now }
	errorsUsecase := NewDashboardErrors(fakeSource{})
	errorsUsecase.now = func() time.Time { return now }

	defaultTrend, err := trend.Execute(context.Background(), DashboardInput{
		Principal: domain.Principal{Role: domain.RoleViewer},
	})
	if err != nil {
		t.Fatal(err)
	}
	defaultErrors, err := errorsUsecase.Execute(context.Background(), DashboardInput{
		Principal: domain.Principal{Role: domain.RoleViewer},
	})
	if err != nil {
		t.Fatal(err)
	}
	wantFrom := now.AddDate(0, 0, -30)
	if !defaultTrend.Period.From.Equal(wantFrom) || !defaultTrend.Period.To.Equal(now) ||
		!defaultErrors.Period.From.Equal(wantFrom) || !defaultErrors.Period.To.Equal(now) {
		t.Fatalf("default periods = %#v %#v", defaultTrend.Period, defaultErrors.Period)
	}

	invalidInputs := []DashboardInput{
		{Principal: domain.Principal{Role: domain.RoleViewer}, From: "2026-08-07T00:00:00Z"},
		{Principal: domain.Principal{Role: domain.RoleViewer}, To: "2026-09-06T00:00:00Z"},
		{Principal: domain.Principal{Role: domain.RoleViewer}, From: "invalid", To: "2026-09-06T00:00:00Z"},
		{Principal: domain.Principal{Role: domain.RoleViewer}, From: "2026-09-06T00:00:00Z", To: "2026-09-06T00:00:00Z"},
	}
	for _, input := range invalidInputs {
		if _, err := trend.Execute(context.Background(), input); domain.ErrorKindOf(err) != domain.KindInvalidPeriod {
			t.Fatalf("trend input %#v err = %v", input, err)
		}
		if _, err := errorsUsecase.Execute(context.Background(), input); domain.ErrorKindOf(err) != domain.KindInvalidPeriod {
			t.Fatalf("errors input %#v err = %v", input, err)
		}
	}
}

func TestDashboardChartsForbiddenAndSourceUnavailable(t *testing.T) {
	source := fakeSource{err: domain.NewError(domain.KindSourceUnavailable, errors.New("adapter details"))}
	trend := NewDashboardExecutionTrend(source)
	errorsUsecase := NewDashboardErrors(source)
	forbidden := DashboardInput{Principal: domain.Principal{Role: domain.Role("operator")}}
	if _, err := trend.Execute(context.Background(), forbidden); domain.ErrorKindOf(err) != domain.KindForbidden {
		t.Fatalf("trend forbidden err = %v", err)
	}
	if _, err := errorsUsecase.Execute(context.Background(), forbidden); domain.ErrorKindOf(err) != domain.KindForbidden {
		t.Fatalf("errors forbidden err = %v", err)
	}
	viewer := DashboardInput{Principal: domain.Principal{Role: domain.RoleViewer}}
	if _, err := trend.Execute(context.Background(), viewer); domain.ErrorKindOf(err) != domain.KindSourceUnavailable {
		t.Fatalf("trend source err = %v", err)
	}
	if _, err := errorsUsecase.Execute(context.Background(), viewer); domain.ErrorKindOf(err) != domain.KindSourceUnavailable {
		t.Fatalf("errors source err = %v", err)
	}
}

func TestDashboardUsecasesFromRepositoryAggregates(t *testing.T) {
	from := "2026-01-15T00:00:00Z"
	to := "2026-04-01T00:00:00Z"
	freshAt := time.Date(2026, 4, 1, 1, 0, 0, 0, time.UTC)
	repository := fakeDashboardRepository{
		summary: domain.DashboardSummaryAggregate{
			TotalUseCases: 8, ActiveUseCases: 6, ExecutionVolume: 4,
			SuccessfulCount: 3, FailedExecutions: 1,
			Freshness: domain.Freshness{
				Status: domain.FreshnessSyncFailed, LastSuccessfulRefreshAt: freshAt,
			},
		},
		trend: []domain.ExecutionTrendAggregate{
			{Bucket: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), Success: 2, Failure: 1},
			{Bucket: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC), Success: 1},
		},
		groups: []domain.ErrorGroup{
			{Code: "timeout", Label: "Timeout", Count: 2},
			{Code: "validation", Label: "Validation", Count: 1},
		},
	}
	input := DashboardInput{
		Principal: domain.Principal{Role: domain.RoleViewer},
		From:      from,
		To:        to,
	}

	summary, err := NewDashboardSummaryFromRepository(repository).Execute(context.Background(), input)
	if err != nil {
		t.Fatal(err)
	}
	if summary.KPIs.TotalUseCases != 8 || summary.KPIs.ActiveUseCases != 6 ||
		summary.KPIs.ExecutionVolume != 4 || summary.KPIs.FailedExecutions != 1 ||
		summary.KPIs.SuccessRate == nil || *summary.KPIs.SuccessRate != 75 ||
		summary.Freshness.Status != domain.FreshnessSyncFailed ||
		!summary.Freshness.LastSuccessfulRefreshAt.Equal(freshAt) {
		t.Fatalf("repository summary = %#v", summary)
	}

	trend, err := NewDashboardExecutionTrendFromRepository(repository).Execute(context.Background(), input)
	if err != nil {
		t.Fatal(err)
	}
	if len(trend.Points) != 3 ||
		trend.Points[0].Success != 2 || trend.Points[0].Failure != 1 ||
		trend.Points[1].Success != 0 || trend.Points[1].Failure != 0 ||
		trend.Points[2].Success != 1 || trend.Points[2].Failure != 0 {
		t.Fatalf("repository trend = %#v", trend.Points)
	}

	groups, err := NewDashboardErrorsFromRepository(repository).Execute(context.Background(), input)
	if err != nil {
		t.Fatal(err)
	}
	if len(groups.Groups) != 2 || groups.Groups[0].Code != "timeout" || groups.Groups[1].Code != "validation" {
		t.Fatalf("repository errors = %#v", groups.Groups)
	}
}
