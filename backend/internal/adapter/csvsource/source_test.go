package csvsource

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const testHeader = "ProcessId,ProcessName,PackageName,EnvironmentName,FullyQualifiedName,CountExecuting,CountPending,CountSuspended,CountResumed,CountSuccessful,CountErrors,CountStopped,AverageDurationInSeconds,AveragePendingTimeInSeconds,TotalRows,EntityId\n"

func TestParsePreservesFailureOnlyRowsAndGroupsByFullName(t *testing.T) {
	fixture := testHeader +
		"10,Error Process,error.pkg,NULL,Area/Full Name,0,0,0,0,0,23,0,162,0,2,110\n" +
		"11,Stopped Process,stopped.pkg,Production,Area/Full Name,0,0,0,0,0,0,7,NULL,NULL,2,111\n"
	importedAt := time.Date(2026, 9, 9, 12, 0, 0, 0, time.UTC)

	snapshot, err := parse([]byte(fixture), importedAt)
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.Rows) != 2 {
		t.Fatalf("rows = %d, want 2", len(snapshot.Rows))
	}
	if snapshot.SourceSnapshotKey == "" || !snapshot.ImportedAt.Equal(importedAt) {
		t.Fatalf("unexpected snapshot metadata: %#v", snapshot)
	}
	errorOnly := snapshot.Rows[0]
	if errorOnly.UseCaseSourceKey != "Area/Full Name" ||
		errorOnly.UseCaseName != "Area/Full Name" ||
		errorOnly.SuccessfulCount != 0 ||
		errorOnly.ErrorCount != 23 ||
		errorOnly.StoppedCount != 0 ||
		errorOnly.FailureCount() != 23 {
		t.Fatalf("error-only row = %#v", errorOnly)
	}
	stoppedOnly := snapshot.Rows[1]
	if stoppedOnly.SuccessfulCount != 0 ||
		stoppedOnly.ErrorCount != 0 ||
		stoppedOnly.StoppedCount != 7 ||
		stoppedOnly.FailureCount() != 7 {
		t.Fatalf("stopped-only row = %#v", stoppedOnly)
	}
}

func TestParseUsesStableContentHash(t *testing.T) {
	fixture := testHeader +
		"10,Process,pkg,NULL,Area/Name,0,0,0,0,1,0,0,NULL,NULL,1,110\n"

	first, err := parse([]byte(fixture), time.Now())
	if err != nil {
		t.Fatal(err)
	}
	second, err := parse([]byte(fixture), time.Now().Add(time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if first.SourceSnapshotKey != second.SourceSnapshotKey {
		t.Fatalf("hash changed: %q != %q", first.SourceSnapshotKey, second.SourceSnapshotKey)
	}
}

func TestParseRejectsInvalidInput(t *testing.T) {
	tests := map[string]string{
		"missing required header": strings.Replace(testHeader, "ProcessName", "Name", 1) +
			"10,Process,pkg,NULL,Area/Name,0,0,0,0,1,0,0,NULL,NULL,1,110\n",
		"negative count": testHeader +
			"10,Process,pkg,NULL,Area/Name,0,0,0,0,-1,0,0,NULL,NULL,1,110\n",
		"empty group": testHeader +
			"10,Process,pkg,NULL,,0,0,0,0,1,0,0,NULL,NULL,1,110\n",
		"invalid ID": testHeader +
			"zero,Process,pkg,NULL,Area/Name,0,0,0,0,1,0,0,NULL,NULL,1,110\n",
	}
	for name, fixture := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := parse([]byte(fixture), time.Now()); err == nil {
				t.Fatal("invalid CSV accepted")
			}
		})
	}
}

func TestAuthoritativeCSVInventoryWhenPresent(t *testing.T) {
	path := filepath.Join("..", "..", "..", "..", "docs", "report", "DataSourceOrchestrator.csv")
	snapshot, err := New(path).Read(context.Background())
	if errors.Is(err, os.ErrNotExist) {
		t.Skip("authoritative CSV is not present in this checkout")
	}
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.Rows) != 213 {
		t.Fatalf("rows = %d, want 213", len(snapshot.Rows))
	}
	for index, row := range snapshot.Rows {
		if row.SourceTotalRows != 213 {
			t.Fatalf("row %d source total = %d, want 213", index+2, row.SourceTotalRows)
		}
	}
}
