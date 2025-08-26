package standard_lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSum(t *testing.T) {
	result := Sum(2, 3)
	expected := 5

	if result != expected {
		t.Errorf("Sum(2, 3) = %d; want %d", result, expected)
	}

	// assert equality
	assert.Equal(t, expected, result, "Sum(2, 3) = %d; want %d", result, expected)
}
