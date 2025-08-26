package testify_require_ex

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSumRequire(t *testing.T) {
	got := Sum(2, 3)
	require.Equal(t, 5, got, "Sum should add two integers")
}
