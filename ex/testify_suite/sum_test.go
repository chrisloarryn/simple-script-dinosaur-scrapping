package testify_suite_ex

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// SumTestSuite demonstrates testify's suite-style tests.
type SumTestSuite struct {
	suite.Suite
}

func (s *SumTestSuite) TestSum() {
	s.Equal(5, Sum(2, 3))
}

func TestSumSuite(t *testing.T) {
	suite.Run(t, new(SumTestSuite))
}
