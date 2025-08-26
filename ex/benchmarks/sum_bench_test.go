package benchmarks_ex

import "testing"

// BenchmarkSum measures performance of Sum.
func BenchmarkSum(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Sum(i, i+1)
	}
}
