package usecase

import (
	"context"
	"math"
	"sort"

	"github.com/mohammadalviyan/rpmp/backend/internal/domain"
)

type UseCaseRepository interface {
	ListUseCaseSnapshots(context.Context, string, string) ([]domain.UseCaseSnapshot, error)
	GetUseCaseSnapshot(context.Context, string) (domain.UseCaseSnapshot, error)
	UseCaseFreshness(context.Context) (domain.Freshness, error)
}

type ListUseCasesInput struct {
	Principal domain.Principal
	Query     string
	Status    string
}

type GetUseCaseInput struct {
	Principal domain.Principal
	ID        string
}

type ListUseCases struct {
	repository UseCaseRepository
}

type GetUseCase struct {
	repository UseCaseRepository
}

func NewListUseCases(repository UseCaseRepository) *ListUseCases {
	return &ListUseCases{repository: repository}
}

func NewGetUseCase(repository UseCaseRepository) *GetUseCase {
	return &GetUseCase{repository: repository}
}

func (uc *ListUseCases) Execute(ctx context.Context, input ListUseCasesInput) (domain.UseCaseList, error) {
	if !input.Principal.Role.CanRead() {
		return domain.UseCaseList{}, domain.NewError(domain.KindForbidden, nil)
	}
	if input.Status != "" && input.Status != string(domain.UseCaseActive) &&
		input.Status != string(domain.UseCaseInactive) {
		return domain.UseCaseList{}, domain.NewError(domain.KindInvalidRequest, nil)
	}

	snapshots, err := uc.repository.ListUseCaseSnapshots(ctx, input.Query, input.Status)
	if err != nil {
		return domain.UseCaseList{}, mapUseCaseDataError(err)
	}
	freshness, err := uc.repository.UseCaseFreshness(ctx)
	if err != nil {
		return domain.UseCaseList{}, mapUseCaseDataError(err)
	}
	items := make([]domain.UseCaseMetrics, 0, len(snapshots))
	for _, snapshot := range snapshots {
		items = append(items, calculateUseCaseMetrics(snapshot))
	}
	return domain.UseCaseList{Freshness: freshness, Items: items}, nil
}

func (uc *GetUseCase) Execute(ctx context.Context, input GetUseCaseInput) (domain.UseCaseDetail, error) {
	if !input.Principal.Role.CanRead() {
		return domain.UseCaseDetail{}, domain.NewError(domain.KindForbidden, nil)
	}
	snapshot, err := uc.repository.GetUseCaseSnapshot(ctx, input.ID)
	if err != nil {
		return domain.UseCaseDetail{}, mapUseCaseDataError(err)
	}
	freshness, err := uc.repository.UseCaseFreshness(ctx)
	if err != nil {
		return domain.UseCaseDetail{}, mapUseCaseDataError(err)
	}
	return domain.UseCaseDetail{
		Freshness: freshness,
		UseCase:   calculateUseCaseMetrics(snapshot),
	}, nil
}

func calculateUseCaseMetrics(snapshot domain.UseCaseSnapshot) domain.UseCaseMetrics {
	metrics := domain.UseCaseMetrics{
		ID:           snapshot.UseCase.ID,
		Name:         snapshot.UseCase.Name,
		Status:       snapshot.UseCase.Status,
		SourceKey:    snapshot.UseCase.SourceKey,
		UpdatedAt:    snapshot.UseCase.UpdatedAt,
		ProcessCount: len(snapshot.Processes),
		Environments: []string{},
		Processes:    snapshot.Processes,
	}
	environments := make(map[string]struct{})
	for _, process := range snapshot.Processes {
		metrics.SuccessfulCount += process.SuccessfulCount
		metrics.ErrorCount += process.ErrorCount
		metrics.StoppedCount += process.StoppedCount
		if process.EnvironmentName != nil && *process.EnvironmentName != "" {
			environments[*process.EnvironmentName] = struct{}{}
		}
	}
	metrics.FailedCount = metrics.ErrorCount + metrics.StoppedCount
	metrics.ExecutionVolume = metrics.SuccessfulCount + metrics.FailedCount
	if metrics.ExecutionVolume > 0 {
		rate := math.Round(
			(float64(metrics.SuccessfulCount)/float64(metrics.ExecutionVolume))*100*100,
		) / 100
		metrics.SuccessRate = &rate
	}
	for environment := range environments {
		metrics.Environments = append(metrics.Environments, environment)
	}
	sort.Strings(metrics.Environments)
	return metrics
}

func mapUseCaseDataError(err error) error {
	switch domain.ErrorKindOf(err) {
	case domain.KindNotFound, domain.KindSourceUnavailable:
		return err
	default:
		return domain.NewError(domain.KindInternal, err)
	}
}
