package dmnet

import (
	"testing"
)

func TestAdd(t *testing.T) {
	const n = 51
	const m = 4
	if ans := Add(n, m); ans != 9 {
		t.Errorf("Error answer: %v", ans)
	}
}
