package implementations

import (
	"testing"
)

func TestLockFreeQueue(t *testing.T) {
	q := NewLockFreeQueue(1024)
	go func() {
		for i := 0; i < 1e6; i++ {
			q.Put(int64(i))
		}
	}()
	for i := 0; i < 1e6; i++ {
		if g := q.Get(); g != int64(i) {
			t.Fatalf("got %d, want %d", g, i)
		}
	}
}
