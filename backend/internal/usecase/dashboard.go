package usecase

import (
	"context"
	"math"
	"time"

	"github.com/mohammadalviyan/rpmp/backend/internal/domain"
)

type Source interface {
	Snapshot(context.Context) (domain.SourceSnapshot, error)
}

type DashboardInput struct {
	Principal domain.Principal
	From      string
	To        string
}

type DashboardSummary struct {
	source Source
	now    func() time.Time
}

type DashboardExecutionTrend struct {
	source Source
	now    func() time.Time
}

type DashboardErrors struct {
	source Source
	now    func() time.Time
}

func NewDashboardSummary(source Source) *DashboardSummary {
	return &DashboardSummary{source: source, now: time.Now}
}

func NewDashboardExecutionTrend(source Source) *DashboardExecutionTrend {
	return &DashboardExecutionTrend{source: source, now: time.Now}
}

func NewDashboardErrors(source Source) *DashboardErrors {
	return &DashboardErrors{source: source, now: time.Now}
}

func (uc *DashboardSummary) Execute(ctx context.Context, input DashboardInput) (domain.DashboardSummary, error) {
	if !input.Principal.Role.CanRead() {
		return domain.DashboardSummary{}, domain.NewError(domain.KindForbidden, nil)
	}

	period, err := resolvePeriod(input.From, input.To, uc.now().UTC())
	if err != nil {
		return domain.DashboardSummary{}, err
	}

	snapshot, err := loadSnapshot(ctx, uc.source)
	if err != nil {
		return domain.DashboardSummary{}, err
	}

	active := 0
	for _, useCase := range snapshot.UseCases {
		if useCase.Status == domain.UseCaseActive {
			active++
		}
	}

	volume := 0
	successful := 0
	failed := 0
	for _, execution := range snapshot.Executions {
		if !inPeriod(execution.StartedAt, period) {
			continue
		}
		volume++
		switch {
		case execution.Status == domain.ExecutionSuccess:
			successful++
		case execution.Status.IsFailure():
			failed++
		}
	}

	kpis := domain.KPIs{
		TotalUseCases:    len(snapshot.UseCases),
		ActiveUseCases:   active,
		ExecutionVolume:  volume,
		FailedExecutions: failed,
	}
	if volume > 0 {
		rate := math.Round((float64(successful)/float64(volume))*100*100) / 100
		kpis.SuccessRate = &rate
	}

	return domain.DashboardSummary{
		Period:    period,
		Freshness: snapshot.Freshness,
		KPIs:      kpis,
	}, nil
}

func (uc *DashboardExecutionTrend) Execute(ctx context.Context, input DashboardInput) (domain.ExecutionTrend, error) {
	if !input.Principal.Role.CanRead() {
		return domain.ExecutionTrend{}, domain.NewError(domain.KindForbidden, nil)
	}

	period, err := resolvePeriod(input.From, input.To, uc.now().UTC())
	if err != nil {
		return domain.ExecutionTrend{}, err
	}
	snapshot, err := loadSnapshot(ctx, uc.source)
	if err != nil {
		return domain.ExecutionTrend{}, err
	}

	firstBucket := time.Date(period.From.Year(), period.From.Month(), 1, 0, 0, 0, 0, time.UTC)
	points := make([]domain.ExecutionTrendPoint, 0)
	pointByBucket := make(map[time.Time]int)
	for bucket := firstBucket; bucket.Before(period.To); bucket = bucket.AddDate(0, 1, 0) {
		pointByBucket[bucket] = len(points)
		points = append(points, domain.ExecutionTrendPoint{
			Bucket: bucket,
			Label:  bucket.Format("Jan"),
		})
	}

	for _, execution := range snapshot.Executions {
		if !inPeriod(execution.StartedAt, period) {
			continue
		}
		bucket := time.Date(
			execution.StartedAt.UTC().Year(),
			execution.StartedAt.UTC().Month(),
			1, 0, 0, 0, 0, time.UTC,
		)
		index, ok := pointByBucket[bucket]
		if !ok {
			continue
		}
		switch {
		case execution.Status == domain.ExecutionSuccess:
			points[index].Success++
		case execution.Status.IsFailure():
			points[index].Failure++
		}
	}

	return domain.ExecutionTrend{Period: period, Points: points}, nil
}

func (uc *DashboardErrors) Execute(ctx context.Context, input DashboardInput) (domain.DashboardErrors, error) {
	if !input.Principal.Role.CanRead() {
		return domain.DashboardErrors{}, domain.NewError(domain.KindForbidden, nil)
	}

	period, err := resolvePeriod(input.From, input.To, uc.now().UTC())
	if err != nil {
		return domain.DashboardErrors{}, err
	}
	snapshot, err := loadSnapshot(ctx, uc.source)
	if err != nil {
		return domain.DashboardErrors{}, err
	}

	groups := []domain.ErrorGroup{
		{Code: "faulted", Label: "Faulted"},
		{Code: "stopped", Label: "Stopped"},
	}
	for _, execution := range snapshot.Executions {
		if !inPeriod(execution.StartedAt, period) {
			continue
		}
		switch execution.Status {
		case domain.ExecutionFailure:
			groups[0].Count++
		case domain.ExecutionException:
			groups[1].Count++
		}
	}

	return domain.DashboardErrors{Period: period, Groups: groups}, nil
}

func loadSnapshot(ctx context.Context, source Source) (domain.SourceSnapshot, error) {
	snapshot, err := source.Snapshot(ctx)
	if err == nil {
		return snapshot, nil
	}
	if domain.ErrorKindOf(err) == domain.KindSourceUnavailable {
		return domain.SourceSnapshot{}, err
	}
	return domain.SourceSnapshot{}, domain.NewError(domain.KindInternal, err)
}

func resolvePeriod(fromRaw, toRaw string, now time.Time) (domain.Period, error) {
	if fromRaw == "" && toRaw == "" {
		return domain.Period{
			From:     now.AddDate(0, 0, -30),
			To:       now,
			Timezone: "UTC",
		}, nil
	}
	if fromRaw == "" || toRaw == "" {
		return domain.Period{}, domain.NewError(domain.KindInvalidPeriod, nil)
	}
	from, err := time.Parse(time.RFC3339, fromRaw)
	if err != nil {
		return domain.Period{}, domain.NewError(domain.KindInvalidPeriod, err)
	}
	to, err := time.Parse(time.RFC3339, toRaw)
	if err != nil {
		return domain.Period{}, domain.NewError(domain.KindInvalidPeriod, err)
	}
	from = from.UTC()
	to = to.UTC()
	if !from.Before(to) {
		return domain.Period{}, domain.NewError(domain.KindInvalidPeriod, nil)
	}
	return domain.Period{From: from, To: to, Timezone: "UTC"}, nil
}

func inPeriod(startedAt time.Time, period domain.Period) bool {
	startedAt = startedAt.UTC()
	return !startedAt.Before(period.From) && startedAt.Before(period.To)
}
