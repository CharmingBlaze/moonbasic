package heap

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestReleaseOnce_DoOnce(t *testing.T) {
	var r ReleaseOnce
	n := 0
	r.Do(func() { n++ })
	r.Do(func() { n++ })
	if n != 1 {
		t.Fatalf("expected fn once, got %d", n)
	}
}

func TestReleaseOnce_Concurrent(t *testing.T) {
	var r ReleaseOnce
	var wg sync.WaitGroup
	var n atomic.Int32
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r.Do(func() { n.Add(1) })
		}()
	}
	wg.Wait()
	if n.Load() != 1 {
		t.Fatalf("expected exactly one concurrent execution, got %d", n.Load())
	}
}
