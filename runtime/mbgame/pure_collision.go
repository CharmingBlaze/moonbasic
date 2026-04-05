package mbgame

import "math"

// 2D / 3D collision and distance helpers (no Raylib — pure math).

func boxCollide2D(x1, y1, w1, h1, x2, y2, w2, h2 float64) bool {
	return x1 < x2+w2 && x1+w1 > x2 && y1 < y2+h2 && y1+h1 > y2
}

func circleCollide2D(x1, y1, r1, x2, y2, r2 float64) bool {
	dx := x2 - x1
	dy := y2 - y1
	rs := r1 + r2
	return dx*dx+dy*dy <= rs*rs
}

func pointInBox2D(px, py, bx, by, bw, bh float64) bool {
	return px >= bx && px <= bx+bw && py >= by && py <= by+bh
}

func pointInCircle2D(px, py, cx, cy, r float64) bool {
	dx := px - cx
	dy := py - cy
	return dx*dx+dy*dy <= r*r
}

func circleBoxCollide2D(cx, cy, cr, bx, by, bw, bh float64) bool {
	// Closest point on AABB to circle center
	nx := math.Max(bx, math.Min(cx, bx+bw))
	ny := math.Max(by, math.Min(cy, by+bh))
	dx := cx - nx
	dy := cy - ny
	return dx*dx+dy*dy <= cr*cr
}

func lineCollide2D(x1, y1, x2, y2, x3, y3, x4, y4 float64) bool {
	d := (x1-x2)*(y3-y4) - (y1-y2)*(x3-x4)
	if math.Abs(d) < 1e-12 {
		return false
	}
	t := ((x1-x3)*(y3-y4) - (y1-y3)*(x3-x4)) / d
	u := -((x1-x2)*(y1-y3) - (y1-y2)*(x1-x3)) / d
	return t >= 0 && t <= 1 && u >= 0 && u <= 1
}

func pointOnLine2D(px, py, lx1, ly1, lx2, ly2, threshold float64) bool {
	// Distance from point to segment
	dx := lx2 - lx1
	dy := ly2 - ly1
	lenSq := dx*dx + dy*dy
	if lenSq < 1e-18 {
		d := math.Hypot(px-lx1, py-ly1)
		return d <= threshold
	}
	t := math.Max(0, math.Min(1, ((px-lx1)*dx+(py-ly1)*dy)/lenSq))
	qx := lx1 + t*dx
	qy := ly1 + t*dy
	return math.Hypot(px-qx, py-qy) <= threshold
}

func sphereCollide3D(x1, y1, z1, r1, x2, y2, z2, r2 float64) bool {
	dx := x2 - x1
	dy := y2 - y1
	dz := z2 - z1
	rs := r1 + r2
	return dx*dx+dy*dy+dz*dz <= rs*rs
}

func aabbCollide3D(minx1, miny1, minz1, maxx1, maxy1, maxz1, minx2, miny2, minz2, maxx2, maxy2, maxz2 float64) bool {
	return minx1 < maxx2 && maxx1 > minx2 &&
		miny1 < maxy2 && maxy1 > miny2 &&
		minz1 < maxz2 && maxz1 > minz2
}

// Sphere vs AABB: box min corner (bx,by,bz) and size (bw,bh,bd).
func sphereBoxCollide3D(sx, sy, sz, sr, bx, by, bz, bw, bh, bd float64) bool {
	cx := math.Max(bx, math.Min(sx, bx+bw))
	cy := math.Max(by, math.Min(sy, by+bh))
	cz := math.Max(bz, math.Min(sz, bz+bd))
	dx := sx - cx
	dy := sy - cy
	dz := sz - cz
	return dx*dx+dy*dy+dz*dz <= sr*sr
}

func pointInAABB3D(px, py, pz, bx, by, bz, bw, bh, bd float64) bool {
	return px >= bx && px <= bx+bw && py >= by && py <= by+bh && pz >= bz && pz <= bz+bd
}

func distance2D(x1, y1, x2, y2 float64) float64 {
	return math.Hypot(x2-x1, y2-y1)
}

func distance3D(x1, y1, z1, x2, y2, z2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	dz := z2 - z1
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func distanceSq2D(x1, y1, x2, y2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	return dx*dx + dy*dy
}

func distanceSq3D(x1, y1, z1, x2, y2, z2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	dz := z2 - z1
	return dx*dx + dy*dy + dz*dz
}
