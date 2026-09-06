package stub

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/mohammadalviyan/rpmp/backend/internal/domain"
)

// Vendor JobState integers stay in this package. Successful=3, Faulted=4, Stopped=6.
const (
	vendorJobSuccessful = 3
	vendorJobFaulted    = 4
	vendorJobStopped    = 6
)

type vendorFixture struct {
	Availability       string          `json:"Availability"`
	LastSuccessfulSync string          `json:"LastSuccessfulSync"`
	Releases           []vendorRelease `json:"Releases"`
	Jobs               []vendorJob     `json:"Jobs"`
}

type vendorRelease struct {
	Key               string `json:"Key"`
	EnvironmentStatus string `json:"EnvironmentStatus"`
}

type vendorJob struct {
	ReleaseKey   string `json:"ReleaseKey"`
	CreationTime string `json:"CreationTime"`
	JobState     int    `json:"JobState"`
	Count        int    `json:"Count"`
}

type Source struct {
	raw []byte
}

func New() (*Source, error) {
	return NewFromJSON(defaultFixture)
}

func NewUnavailable() *Source {
	return &Source{}
}

func NewFromJSON(raw []byte) (*Source, error) {
	if len(raw) == 0 {
		return nil, errors.New("stub fixture is empty")
	}
	return &Source{raw: raw}, nil
}

func (s *Source) Snapshot(context.Context) (domain.SourceSnapshot, error) {
	if len(s.raw) == 0 {
		return domain.SourceSnapshot{}, domain.NewError(domain.KindSourceUnavailable, errors.New("vendor export timed out"))
	}

	var fixture vendorFixture
	if err := json.Unmarshal(s.raw, &fixture); err != nil {
		return domain.SourceSnapshot{}, domain.NewError(domain.KindSourceUnavailable, err)
	}
	if fixture.Availability != "Available" {
		return domain.SourceSnapshot{}, domain.NewError(domain.KindSourceUnavailable, fmt.Errorf("vendor availability %q", fixture.Availability))
	}

	refreshedAt, err := time.Parse(time.RFC3339, fixture.LastSuccessfulSync)
	if err != nil {
		return domain.SourceSnapshot{}, domain.NewError(domain.KindSourceUnavailable, err)
	}

	useCases := make([]domain.UseCase, 0, len(fixture.Releases))
	for _, release := range fixture.Releases {
		status := domain.UseCaseInactive
		if release.EnvironmentStatus == "Available" {
			status = domain.UseCaseActive
		}
		useCases = append(useCases, domain.UseCase{ID: release.Key, Status: status})
	}

	var executions []domain.Execution
	seq := 0
	for _, job := range fixture.Jobs {
		status, ok := mapJobState(job.JobState)
		if !ok {
			return domain.SourceSnapshot{}, domain.NewError(domain.KindSourceUnavailable, fmt.Errorf("unknown JobState %d", job.JobState))
		}
		startedAt, err := time.Parse(time.RFC3339, job.CreationTime)
		if err != nil {
			return domain.SourceSnapshot{}, domain.NewError(domain.KindSourceUnavailable, err)
		}
		count := job.Count
		if count <= 0 {
			count = 1
		}
		for i := 0; i < count; i++ {
			seq++
			executions = append(executions, domain.Execution{
				ID:        fmt.Sprintf("%s-%d", job.ReleaseKey, seq),
				UseCaseID: job.ReleaseKey,
				StartedAt: startedAt.UTC(),
				Status:    status,
			})
		}
	}

	return domain.SourceSnapshot{
		UseCases:   useCases,
		Executions: executions,
		Freshness: domain.Freshness{
			Status:                  domain.FreshnessFresh,
			LastSuccessfulRefreshAt: refreshedAt.UTC(),
		},
	}, nil
}

func mapJobState(state int) (domain.ExecutionStatus, bool) {
	switch state {
	case vendorJobSuccessful:
		return domain.ExecutionSuccess, true
	case vendorJobFaulted:
		return domain.ExecutionFailure, true
	case vendorJobStopped:
		return domain.ExecutionException, true
	default:
		return "", false
	}
}
