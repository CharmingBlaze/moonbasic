package include

import (
	"path/filepath"
	"strings"
)

func stackContains(stack []string, abs string) bool {
	for _, p := range stack {
		if p == abs {
			return true
		}
	}
	return false
}

func formatCircularChain(stack []string, closesTo string) string {
	parts := make([]string, 0, len(stack)+1)
	for _, p := range stack {
		parts = append(parts, filepath.Base(p))
	}
	parts = append(parts, filepath.Base(closesTo))
	return strings.Join(parts, " → ")
}
