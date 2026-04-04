package runtime

import (
	"fmt"
	"strings"
)

// BestSimilarCommand picks the candidate closest to target by Levenshtein distance (uppercase compare).
func BestSimilarCommand(target string, candidates []string, maxDist int) (string, bool) {
	tu := strings.ToUpper(strings.TrimSpace(target))
	best := maxDist + 1
	var bk string
	for _, c := range candidates {
		cu := strings.ToUpper(strings.TrimSpace(c))
		d := levenshteinDist(tu, cu)
		if d < best {
			best = d
			bk = cu
		}
	}
	if best <= maxDist && bk != "" {
		return bk, true
	}
	return "", false
}

// FormatUnknownRegistryCommand builds a registry miss error with optional did-you-mean.
func FormatUnknownRegistryCommand(key string, candidates []string) error {
	if alt, ok := BestSimilarCommand(key, candidates, 3); ok {
		return fmt.Errorf("unknown command %q\n  Did you mean %s?\n  Hint: Built-ins use dotted names like CAMERA.SETPOS (see docs/API_CONSISTENCY.md).", key, alt)
	}
	return fmt.Errorf("unknown command %q\n  Hint: Check namespace and spelling (e.g. TIME.GETFPS, CAMERA.SETPOS). See docs/API_CONSISTENCY.md.", key)
}

func levenshteinDist(a, b string) int {
	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}
	if len(a) > len(b) {
		a, b = b, a
	}
	row := make([]int, len(b)+1)
	for j := range row {
		row[j] = j
	}
	for i := 1; i <= len(a); i++ {
		prev := row[0]
		row[0] = i
		for j := 1; j <= len(b); j++ {
			cur := row[j]
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			row[j] = min3int(prev+1, row[j-1]+1, prev+cost)
			prev = cur
		}
	}
	return row[len(b)]
}

func min3int(a, b, c int) int {
	if a <= b && a <= c {
		return a
	}
	if b <= c {
		return b
	}
	return c
}
