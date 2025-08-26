package table_subtests_ex

import "testing"

// Demonstrates table-driven tests with subtests (t.Run).
func TestSum_TableDriven(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"both positive", 2, 3, 5},
		{"with zero a", 0, 7, 7},
		{"with zero b", 9, 0, 9},
		{"negatives", -2, -3, -5},
		{"mixed", -2, 5, 3},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Sum(tc.a, tc.b)
			if got != tc.want {
				t.Fatalf("%s: Sum(%d,%d)=%d; want %d", tc.name, tc.a, tc.b, got, tc.want)
			}
		})
	}
}
