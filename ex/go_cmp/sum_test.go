package go_cmp_ex

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Demonstrates using go-cmp to compare values in tests.
func TestSum_WithGoCmp(t *testing.T) {
	got := Sum(2, 3)
	want := 5
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("Sum diff (-want +got):\n%s", diff)
	}
}
