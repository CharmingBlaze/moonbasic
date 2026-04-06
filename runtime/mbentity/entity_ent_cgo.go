//go:build cgo || (windows && !cgo)

package mbentity

import (
	"moonbasic/vm/heap"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type entKind int8

const (
	entKindEmpty entKind = iota
	entKindBox
	entKindSphere
	entKindCylinder
	entKindPlane
	entKindMesh
	entKindModel
)

// ent is a lightweight Blitz-style object (integer id, optional Raylib model mesh).
type ent struct {
	id int64
	kind entKind

	pos               rl.Vector3
	pitch, yaw, roll float32
	scale             rl.Vector3
	r, g, b           uint8
	alpha             float32
	hidden            bool
	shininess         float32
	texHandle         heap.Handle
	fxFlags           int32
	blendMode         int32 // -1 = default; else rl.BlendMode when drawing
	drawOrder         int32

	w, h, d  float32
	radius   float32
	segH, segV int32
	cylH     float32

	useSphere bool
	static    bool
	vel       rl.Vector3
	gravity   float32
	onGround  bool
	mass      float32
	friction  float32
	bounce    float32
	slide     bool
	pickMode  int32

	collided  bool
	otherID   int64
	collType  int32
	hasHit    bool
	hitX, hitY, hitZ   float32
	hitNX, hitNY, hitNZ float32

	parentID int64
	name     string

	rlModel   rl.Model
	hasRLModel bool
	loadPath  string
	// Loaded via rl.LoadModelAnimations (not embedded on rl.Model in raylib-go 0.56+).
	modelAnims []rl.ModelAnimation

	animIndex int32
	animTime  float32
	animSpeed float32
	animMode  int32 // 0=loop, 1=one shot (simplified)
	animLen   float32
}

func newDefaultEnt(id int64) *ent {
	return &ent{
		id:        id,
		kind:      entKindEmpty,
		scale:     rl.Vector3{X: 1, Y: 1, Z: 1},
		alpha:     1,
		r:         200, g: 200, b: 255,
		blendMode: -1,
		segH:      16, segV: 16,
		mass:      1,
		friction:  0.9,
		bounce:    0,
		animSpeed: 0,
	}
}
