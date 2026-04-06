package strmod

import "testing"

func TestRuneIndex(t *testing.T) {
	if got := runeIndex("hello", "ll", 0); got != 2 {
		t.Fatalf("runeIndex(hello, ll, 0) = %d want 2", got)
	}
	if got := runeIndex("hello", "ll", 3); got != -1 {
		t.Fatalf("runeIndex(hello, ll, 3) = %d want -1", got)
	}
}
