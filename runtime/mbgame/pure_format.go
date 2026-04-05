package mbgame

import (
	"fmt"
	"strconv"
	"strings"
)

func formatInt(n int, digits int) string {
	if digits < 1 {
		digits = 1
	}
	s := strconv.Itoa(n)
	if len(s) >= digits {
		return s
	}
	return strings.Repeat("0", digits-len(s)) + s
}

func formatScore(n int64) string {
	s := strconv.FormatInt(n, 10)
	if len(s) <= 3 {
		return s
	}
	var b strings.Builder
	mod := len(s) % 3
	if mod > 0 {
		b.WriteString(s[:mod])
		if len(s) > mod {
			b.WriteString(",")
		}
	}
	for i := mod; i < len(s); i += 3 {
		if i > mod {
			b.WriteString(",")
		}
		b.WriteString(s[i : i+3])
	}
	return b.String()
}

func formatTime(seconds int) string {
	s := seconds % 60
	m := (seconds / 60) % 60
	h := seconds / 3600
	if h > 0 {
		return fmt.Sprintf("%d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%d:%02d", m, s)
}

func formatTime2(seconds int) string {
	s := seconds % 60
	m := (seconds / 60) % 60
	h := seconds / 3600
	return fmt.Sprintf("%d:%02d:%02d", h, m, s)
}
