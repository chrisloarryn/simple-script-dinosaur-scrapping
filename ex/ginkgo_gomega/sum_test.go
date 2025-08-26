package ginkgo_gomega_ex

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sum", func() {
	It("adds two integers", func() {
		Expect(Sum(2, 3)).To(Equal(5))
	})
})
