package stub

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/mohammadalviyan/rpmp/backend/internal/domain"
)

func TestSnapshotNormalizesVendorFixture(t *testing.T) {
	source, err := New()
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := source.Snapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.Freshness.Status != domain.FreshnessFresh {
		t.Fatalf("freshness = %#v", snapshot.Freshness)
	}
	wantRefresh := time.Date(2026, 9, 6, 0, 0, 0, 0, time.UTC)
	if !snapshot.Freshness.LastSuccessfulRefreshAt.Equal(wantRefresh) {
		t.Fatalf("refresh at = %s", snapshot.Freshness.LastSuccessfulRefreshAt)
	}
	if len(snapshot.UseCases) != 12 {
		t.Fatalf("use cases = %d", len(snapshot.UseCases))
	}
	active := 0
	for _, useCase := range snapshot.UseCases {
		if useCase.Status == domain.UseCaseActive {
			active++
		}
	}
	if active != 9 {
		t.Fatalf("active use cases = %d", active)
	}

	statuses := map[domain.ExecutionStatus]int{}
	for _, execution := range snapshot.Executions {
		statuses[execution.Status]++
	}
	if statuses[domain.ExecutionSuccess] != 106 || statuses[domain.ExecutionFailure] != 7 || statuses[domain.ExecutionException] != 2 {
		t.Fatalf("execution statuses = %#v", statuses)
	}
}

func TestSnapshotKeepsVendorTermsOutOfDomain(t *testing.T) {
	source, err := New()
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := source.Snapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	raw := defaultFixture
	if !strings.Contains(string(raw), `"JobState"`) || !strings.Contains(string(raw), `"EnvironmentStatus"`) {
		t.Fatal("fixture no longer contains vendor fields")
	}
	for _, useCase := range snapshot.UseCases {
		if useCase.Status != domain.UseCaseActive && useCase.Status != domain.UseCaseInactive {
			t.Fatalf("vendor status leaked: %q", useCase.Status)
		}
	}
	for _, execution := range snapshot.Executions {
		switch execution.Status {
		case domain.ExecutionSuccess, domain.ExecutionFailure, domain.ExecutionException:
		default:
			t.Fatalf("vendor job state leaked: %q", execution.Status)
		}
	}
}

func TestUnavailableSourceReturnsTypedError(t *testing.T) {
	_, err := NewUnavailable().Snapshot(context.Background())
	if domain.ErrorKindOf(err) != domain.KindSourceUnavailable {
		t.Fatalf("err = %v", err)
	}
	if !strings.Contains(err.Error(), string(domain.KindSourceUnavailable)) {
		t.Fatalf("typed error string = %q", err.Error())
	}
}

func TestUnknownVendorStateIsUnavailable(t *testing.T) {
	source, err := NewFromJSON([]byte(`{
		"Availability":"Available",
		"LastSuccessfulSync":"2026-09-06T00:00:00Z",
		"Releases":[{"Key":"uc-01","EnvironmentStatus":"Available"}],
		"Jobs":[{"ReleaseKey":"uc-01","CreationTime":"2026-08-20T00:00:00Z","JobState":99,"Count":1}]
	}`))
	if err != nil {
		t.Fatal(err)
	}
	_, snapshotErr := source.Snapshot(context.Background())
	if domain.ErrorKindOf(snapshotErr) != domain.KindSourceUnavailable {
		t.Fatalf("err = %v", snapshotErr)
	}
}
