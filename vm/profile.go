package vm

import (
	"sort"

	"moonbasic/lineprof"
)

var _ lineprof.LineProfiler = (*ProfileRecorder)(nil)

// ProfileRecorder counts VM instructions executed per 1-based source line.
type ProfileRecorder struct {
	LineHits map[int]uint64
}

// NewProfileRecorder returns an empty recorder.
func NewProfileRecorder() *ProfileRecorder {
	return &ProfileRecorder{LineHits: make(map[int]uint64)}
}

// RecordLine increments the counter for a source line (ignored if line < 1).
func (p *ProfileRecorder) RecordLine(line int) {
	if p == nil || line < 1 {
		return
	}
	if p.LineHits == nil {
		p.LineHits = make(map[int]uint64)
	}
	p.LineHits[line]++
}

// TopProfileLines returns up to n source lines with the highest hit counts (stable tie-break by line).
func TopProfileLines(p *ProfileRecorder, n int) []struct {
	Line  int
	Count uint64
} {
	if p == nil || len(p.LineHits) == 0 {
		return nil
	}
	type pair struct {
		line  int
		count uint64
	}
	var ps []pair
	for ln, c := range p.LineHits {
		ps = append(ps, pair{ln, c})
	}
	sort.Slice(ps, func(i, j int) bool {
		if ps[i].count == ps[j].count {
			return ps[i].line < ps[j].line
		}
		return ps[i].count > ps[j].count
	})
	if n > len(ps) {
		n = len(ps)
	}
	out := make([]struct {
		Line  int
		Count uint64
	}, n)
	for i := 0; i < n; i++ {
		out[i].Line = ps[i].line
		out[i].Count = ps[i].count
	}
	return out
}
