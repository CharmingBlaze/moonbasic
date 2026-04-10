package mbgrid

// findPath is 4-connected BFS on walkable cells (cell >= 0). Returns linear indices or nil.
func findPath(g *gridObj, six, siz, tix, tiz int) []int {
	if g == nil || !g.contains(six, siz) || !g.contains(tix, tiz) {
		return nil
	}
	if g.cells[g.idx(six, siz)] < 0 || g.cells[g.idx(tix, tiz)] < 0 {
		return nil
	}
	start := g.idx(six, siz)
	goal := g.idx(tix, tiz)
	if start == goal {
		return []int{start}
	}
	visited := make([]bool, g.gw*g.gd)
	parent := make([]int, g.gw*g.gd)
	for i := range parent {
		parent[i] = -1
	}
	type qn struct{ ix, iz int }
	q := make([]qn, 0, 64)
	q = append(q, qn{six, siz})
	visited[start] = true
	head := 0
	dirs := [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for head < len(q) {
		cur := q[head]
		head++
		ci := g.idx(cur.ix, cur.iz)
		for _, d := range dirs {
			nix, niz := cur.ix+d[0], cur.iz+d[1]
			if !g.contains(nix, niz) {
				continue
			}
			ni := g.idx(nix, niz)
			if g.cells[ni] < 0 || visited[ni] {
				continue
			}
			parent[ni] = ci
			if ni == goal {
				var chain []int
				for i := goal; i >= 0; {
					chain = append(chain, i)
					if i == start {
						break
					}
					i = parent[i]
				}
				for i, j := 0, len(chain)-1; i < j; i, j = i+1, j-1 {
					chain[i], chain[j] = chain[j], chain[i]
				}
				return chain
			}
			visited[ni] = true
			q = append(q, qn{nix, niz})
		}
	}
	return nil
}
