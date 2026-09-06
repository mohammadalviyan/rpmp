package domain

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrorKindOfWrappedTypedError(t *testing.T) {
	cause := errors.New("private detail")
	err := fmt.Errorf("operation failed: %w", NewError(KindForbidden, cause))
	if got := ErrorKindOf(err); got != KindForbidden {
		t.Fatalf("kind = %q, want %q", got, KindForbidden)
	}
	if !errors.Is(err, cause) {
		t.Fatal("typed error did not preserve its cause")
	}
}

func TestRoleCanRead(t *testing.T) {
	if !RoleViewer.CanRead() || !RoleAdmin.CanRead() {
		t.Fatal("viewer or admin cannot read")
	}
	if Role("operator").CanRead() {
		t.Fatal("unknown role can read")
	}
}

func TestExecutionStatusFailureGrouping(t *testing.T) {
	if !ExecutionFailure.IsFailure() || !ExecutionException.IsFailure() {
		t.Fatal("failure grouping missed a stub failure status")
	}
	if ExecutionSuccess.IsFailure() {
		t.Fatal("success counted as failure")
	}
}
