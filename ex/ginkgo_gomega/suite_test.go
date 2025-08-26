package ginkgo_gomega_ex

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGinkgoGomega(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GinkgoGomega Suite")
}
