package gomega_t_ex

import (
	"testing"

	. "github.com/onsi/gomega"
)

// Demonstrates using Gomega with the standard testing.T (no Ginkgo).
func TestSum_WithGomegaT(t *testing.T) {
	g := NewWithT(t)
	g.Expect(Sum(2, 3)).To(Equal(5))
	g.Expect(Sum(-2, 5)).To(Equal(3))
}
