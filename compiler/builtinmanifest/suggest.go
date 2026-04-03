package builtinmanifest

import (
	"sort"
	"strings"
)

// EditDistance returns the Levenshtein distance between a and b.
func EditDistance(a, b string) int {
	return levenshtein(a, b)
}

// levenshtein computes edit distance (small strings only; manifest keys are short).
func levenshtein(a, b string) int {
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
			row[j] = min3(prev+1, row[j-1]+1, prev+cost)
			prev = cur
		}
	}
	return row[len(b)]
}

func min3(a, b, c int) int {
	if a <= b && a <= c {
		return a
	}
	if b <= c {
		return b
	}
	return c
}

// KeysWithNamespacePrefix returns sorted manifest keys starting with "NS." (namespace already uppercased).
func (t *Table) KeysWithNamespacePrefix(ns string) []string {
	if t == nil || t.Commands == nil {
		return nil
	}
	p := ns + "."
	var out []string
	for k := range t.Commands {
		if strings.HasPrefix(k, p) {
			out = append(out, k)
		}
	}
	sort.Strings(out)
	return out
}

// BestSimilarKey picks the manifest key with smallest Levenshtein distance to target.
// maxDist caps suggestions (e.g. 4); returns "", false if no candidate within maxDist.
func (t *Table) BestSimilarKey(target string, maxDist int) (key string, ok bool) {
	if t == nil || t.Commands == nil {
		return "", false
	}
	best := maxDist + 1
	var bestKey string
	for k := range t.Commands {
		d := EditDistance(target, k)
		if d < best {
			best = d
			bestKey = k
		}
	}
	if best <= maxDist {
		return bestKey, true
	}
	return "", false
}

// FormatNamespaceListing returns wrapped lines of "NS.*" keys for error hints (max ~72 chars per line).
func FormatNamespaceListing(ns string, keys []string, maxLine int) string {
	if len(keys) == 0 {
		return ""
	}
	prefix := "Available: "
	var lines []string
	var b strings.Builder
	b.WriteString(prefix)
	first := true
	for _, k := range keys {
		sep := " "
		if first {
			sep = ""
			first = false
		}
		if b.Len()+len(sep)+len(k) > maxLine && b.Len() > len(prefix) {
			lines = append(lines, b.String())
			b.Reset()
			b.WriteString(strings.Repeat(" ", len(prefix)))
		}
		b.WriteString(sep)
		b.WriteString(k)
	}
	if b.Len() > 0 {
		lines = append(lines, b.String())
	}
	return strings.Join(lines, "\n")
}
