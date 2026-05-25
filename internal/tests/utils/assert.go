// Package utils holds small assertion helpers shared across the per-domain
// test packages under internal/tests.
package utils

import (
	"testing"

	"profconnect-api/internal/domain"
)

// AssertAppErr fails the test unless err is a *domain.AppError of the given
// type. Used everywhere we need to verify that a usecase mapped a precondition
// to the right HTTP-equivalent error category.
func AssertAppErr(t *testing.T, err error, want domain.ErrorType) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error of type %q, got nil", want)
	}
	appErr, ok := domain.IsAppError(err)
	if !ok {
		t.Fatalf("expected AppError, got %T: %v", err, err)
	}
	if appErr.Type != want {
		t.Fatalf("expected error type %q, got %q (msg=%q)", want, appErr.Type, appErr.Message)
	}
}
