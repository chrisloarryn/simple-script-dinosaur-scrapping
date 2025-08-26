package fuzz_ex

import "testing"

// FuzzSumCommutative fuzz-tests the commutativity property of Sum.
func FuzzSumCommutative(f *testing.F) {
	// Seed corpus
	f.Add(int(0), int(0))
	f.Add(1, 2)
	f.Add(-1, 5)

	f.Fuzz(func(t *testing.T, a, b int) {
		if Sum(a, b) != Sum(b, a) {
			t.Fatalf("Sum not commutative for (%d,%d)", a, b)
		}
	})
}
