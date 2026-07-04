package postgres

import "testing"

func TestNullIfEmpty(t *testing.T) {
	t.Parallel()

	if got := nullIfEmpty(""); got != nil {
		t.Fatalf("empty: got %v, want nil", got)
	}
	if got := nullIfEmpty("x"); got != "x" {
		t.Fatalf("non-empty: got %v", got)
	}
}

func TestNullIfNonPositive(t *testing.T) {
	t.Parallel()

	if got := nullIfNonPositive(0); got != nil {
		t.Fatalf("zero: got %v, want nil", got)
	}
	if got := nullIfNonPositive(-1); got != nil {
		t.Fatalf("negative: got %v, want nil", got)
	}
	if got := nullIfNonPositive(42); got != 42 {
		t.Fatalf("positive: got %v", got)
	}
}
