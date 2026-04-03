package errors

import "testing"

func TestFormat(t *testing.T) {
	e := NewParseError("game.mbc", 14, 8, "Expected '(' after 'RENDERFRAME'", "RENDERFRAME", "All commands require parentheses: RENDERFRAME()")
	s := Format(e)
	if len(s) < 20 {
		t.Fatal(s)
	}
}
