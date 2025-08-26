package testify_ex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSum(t *testing.T) {
	got := Sum(2, 3)
	assert.Equal(t, 5, got, "Sum should add two integers")
}
