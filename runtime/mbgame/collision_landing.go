package mbgame

import "math"

// SphereFeetInBoxTopSupport reports whether (px,py,pz,pr) is in the same horizontal / feet-height
// band as BoxTopLandSnap (ignoring vertical velocity). Used to detect resting contact so we can
// skip gravity without fighting the snap pass.
func SphereFeetInBoxTopSupport(px, py, pz, pr, bx, by, bz, bw, bh, bd float64) bool {
	halfW := bw*0.5 + pr
	halfD := bd*0.5 + pr
	if math.Abs(px-bx) > halfW || math.Abs(pz-bz) > halfD {
		return false
	}
	top := by + bh*0.5
	feet := py - pr
	return feet <= top+0.22 && feet >= top-1.25
}

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
	if feet <= top+0.22 && feet >= top-1.25 {
		return top + pr
	}
	return 0
}
