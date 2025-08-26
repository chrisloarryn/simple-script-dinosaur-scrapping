package gomock_ex

import (
	"testing"

	"github.com/golang/mock/gomock"
)

func TestComputeWithMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockAdder(ctrl)
	m.EXPECT().Sum(2, 3).Return(5)

	got := Compute(m, 2, 3)
	if got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}
}
