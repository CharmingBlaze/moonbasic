package mbgame

import "math"

// BoxTopLandSnap returns sphere-centre Y to snap to when landing on the top of an AABB box,
// or 0 when there is no landing. Same rules as BOXTOPLAND.
func BoxTopLandSnap(px, py, pz, pvy, pr, bx, by, bz, bw, bh, bd float64) float64 {
	if pvy > 0 {
		return 0
	}
	halfW := bw*0.5 + pr
	halfD := bd*0.5 + pr
	if math.Abs(px-bx) > halfW || math.Abs(pz-bz) > halfD {
		return 0
	}
	top := by + bh*0.5
	feet := py - pr
	if feet <= top+0.12 && feet >= top-0.55 {
		return top + pr
	}
	return 0
}
