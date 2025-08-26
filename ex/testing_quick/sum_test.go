package testing_quick_ex

import (
	"testing"
	"testing/quick"
)

// Property-based test: Sum is commutative.
func TestSum_QuickCommutative(t *testing.T) {
	f := func(a, b int) bool {
		return Sum(a, b) == Sum(b, a)
	}
	if err := quick.Check(f, nil); err != nil {
		t.Fatalf("commutativity failed: %v", err)
	}
}
