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

func NewDashboardSummary(source Source) *DashboardSummary {
	return &DashboardSummary{source: source, now: time.Now}
}

func (uc *DashboardSummary) Execute(ctx context.Context, input DashboardInput) (domain.DashboardSummary, error) {
	if !input.Principal.Role.CanRead() {
		return domain.DashboardSummary{}, domain.NewError(domain.KindForbidden, nil)
	}

	period, err := resolvePeriod(input.From, input.To, uc.now().UTC())
	if err != nil {
		return domain.DashboardSummary{}, err
	}

	snapshot, err := uc.source.Snapshot(ctx)
	if err != nil {
		if domain.ErrorKindOf(err) == domain.KindSourceUnavailable {
			return domain.DashboardSummary{}, err
		}
		return domain.DashboardSummary{}, domain.NewError(domain.KindInternal, err)
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
