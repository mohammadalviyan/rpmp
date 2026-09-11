package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mohammadalviyan/rpmp/backend/internal/domain"
)

type fakeUseCaseRepository struct {
	snapshots []domain.UseCaseSnapshot
	snapshot  domain.UseCaseSnapshot
	freshness domain.Freshness
	err       error
	query     string
	status    string
	id        string
}

func (f *fakeUseCaseRepository) ListUseCaseSnapshots(
	_ context.Context,
	query string,
	status string,
) ([]domain.UseCaseSnapshot, error) {
	f.query, f.status = query, status
	return f.snapshots, f.err
}

func (f *fakeUseCaseRepository) GetUseCaseSnapshot(
	_ context.Context,
	id string,
) (domain.UseCaseSnapshot, error) {
	f.id = id
	return f.snapshot, f.err
}

func (f *fakeUseCaseRepository) UseCaseFreshness(context.Context) (domain.Freshness, error) {
	return f.freshness, f.err
}

func TestListUseCasesCalculatesFrozenMetrics(t *testing.T) {
	prod := "PROD"
	empty := ""
	repository := &fakeUseCaseRepository{
		freshness: domain.Freshness{
			Status: domain.FreshnessFresh, LastSuccessfulRefreshAt: time.Date(2026, 9, 11, 2, 0, 0, 0, time.UTC),
		},
		snapshots: []domain.UseCaseSnapshot{
			{
				UseCase: domain.StoredUseCase{
					ID: "use-case-1", Name: "Finance/Invoice.Process", Status: domain.UseCaseActive,
				},
				Processes: []domain.UseCaseProcessSnapshot{
					{SourceProcessKey: "1", EnvironmentName: &prod, SuccessfulCount: 50, ErrorCount: 1},
					{SourceProcessKey: "2", EnvironmentName: &prod, SuccessfulCount: 50, ErrorCount: 2, StoppedCount: 1},
					{SourceProcessKey: "3", EnvironmentName: &empty},
				},
			},
			{
				UseCase:   domain.StoredUseCase{ID: "use-case-2", Name: "Zero", Status: domain.UseCaseInactive},
				Processes: []domain.UseCaseProcessSnapshot{},
			},
		},
	}

	result, err := NewListUseCases(repository).Execute(context.Background(), ListUseCasesInput{
		Principal: domain.Principal{Role: domain.RoleViewer},
		Query:     "invoice",
		Status:    "active",
	})
	if err != nil {
		t.Fatal(err)
	}
	if repository.query != "invoice" || repository.status != "active" {
		t.Fatalf("filters = %q %q", repository.query, repository.status)
	}
	item := result.Items[0]
	if item.ProcessCount != 3 || item.SuccessfulCount != 100 || item.ErrorCount != 3 ||
		item.StoppedCount != 1 || item.FailedCount != 4 || item.ExecutionVolume != 104 ||
		item.SuccessRate == nil || *item.SuccessRate != 96.15 ||
		len(item.Environments) != 1 || item.Environments[0] != "PROD" {
		t.Fatalf("metrics = %#v", item)
	}
	zero := result.Items[1]
	if zero.ProcessCount != 0 || zero.ExecutionVolume != 0 || zero.SuccessRate != nil ||
		len(zero.Environments) != 0 {
		t.Fatalf("zero metrics = %#v", zero)
	}
}

func TestUseCasesAuthorizationValidationAndErrors(t *testing.T) {
	repository := &fakeUseCaseRepository{}
	list := NewListUseCases(repository)
	for _, role := range []domain.Role{domain.RoleViewer, domain.RoleAdmin} {
		if _, err := list.Execute(context.Background(), ListUseCasesInput{
			Principal: domain.Principal{Role: role},
		}); err != nil {
			t.Fatalf("role %s err = %v", role, err)
		}
	}
	if _, err := list.Execute(context.Background(), ListUseCasesInput{
		Principal: domain.Principal{Role: domain.Role("operator")},
	}); domain.ErrorKindOf(err) != domain.KindForbidden {
		t.Fatalf("forbidden err = %v", err)
	}
	if _, err := list.Execute(context.Background(), ListUseCasesInput{
		Principal: domain.Principal{Role: domain.RoleViewer}, Status: "paused",
	}); domain.ErrorKindOf(err) != domain.KindInvalidRequest {
		t.Fatalf("invalid status err = %v", err)
	}

	repository.err = domain.NewError(domain.KindNotFound, errors.New("missing"))
	if _, err := NewGetUseCase(repository).Execute(context.Background(), GetUseCaseInput{
		Principal: domain.Principal{Role: domain.RoleAdmin}, ID: "missing",
	}); domain.ErrorKindOf(err) != domain.KindNotFound {
		t.Fatalf("not found err = %v", err)
	}
}

func TestGetUseCaseReturnsDetailProcesses(t *testing.T) {
	environment := "UAT"
	updatedAt := time.Date(2026, 9, 11, 2, 0, 0, 0, time.UTC)
	repository := &fakeUseCaseRepository{
		snapshot: domain.UseCaseSnapshot{
			UseCase: domain.StoredUseCase{
				ID: "use-case-1", SourceKey: "Finance/Invoice.Process",
				Name: "Finance/Invoice.Process", Status: domain.UseCaseActive, UpdatedAt: updatedAt,
			},
			Processes: []domain.UseCaseProcessSnapshot{
				{
					SourceProcessKey: "42", ProcessName: "Invoice", PackageName: "Invoice.1.0.0",
					EnvironmentName: &environment, SuccessfulCount: 2, ErrorCount: 1,
					ExecutingCount: 3, PendingCount: 4, SuspendedCount: 5, ResumedCount: 6,
				},
			},
		},
		freshness: domain.Freshness{Status: domain.FreshnessFresh, LastSuccessfulRefreshAt: updatedAt},
	}
	result, err := NewGetUseCase(repository).Execute(context.Background(), GetUseCaseInput{
		Principal: domain.Principal{Role: domain.RoleViewer}, ID: "use-case-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if repository.id != "use-case-1" || result.UseCase.SourceKey != "Finance/Invoice.Process" ||
		len(result.UseCase.Processes) != 1 || result.UseCase.Processes[0].PendingCount != 4 {
		t.Fatalf("detail = %#v", result)
	}
}
