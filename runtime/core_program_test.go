package runtime

import (
	"testing"
	"time"

	"moonbasic/vm/value"
)

func TestSleepDuration(t *testing.T) {
	d, err := sleepDuration(value.FromInt(100))
	if err != nil || d != 100*time.Millisecond {
		t.Fatalf("int ms: got %v %v", d, err)
	}
	d, err = sleepDuration(value.FromFloat(0.5))
	if err != nil || d != 500*time.Millisecond {
		t.Fatalf("float sec: got %v %v", d, err)
	}
	_, err = sleepDuration(value.Nil)
	if err == nil {
		t.Fatal("expected error for nil")
	}
}
