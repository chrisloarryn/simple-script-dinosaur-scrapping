package gomock_ex

// Adder defines a dependency that can sum two integers.
type Adder interface {
	Sum(a, b int) int
}

// Compute delegates the sum operation to the provided Adder.
func Compute(adder Adder, a, b int) int {
	return adder.Sum(a, b)
}
