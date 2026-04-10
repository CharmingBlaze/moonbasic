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
	entKindCone
	entKindMesh
	entKindModel
)

// ent is a lightweight Blitz-style object (integer id, optional Raylib model mesh).
type ent struct {
	id   int64
	kind entKind

	pos              rl.Vector3
	pitch, yaw, roll float32
	scale            rl.Vector3
	r, g, b          uint8
	alpha            float32
	hidden           bool
	shininess        float32
	texHandle        heap.Handle
	fxFlags          int32
	blendMode        int32 // -1 = default; else rl.BlendMode when drawing
	drawOrder        int32
	cullMode         int32 // 0=Auto (Frustum), 1=Force Visible, 2=Force Hidden

	w, h, d    float32
	radius     float32
	segH, segV int32
	cylH       float32

	useSphere bool
	static    bool
	vel       rl.Vector3
	gravity   float32
	onGround  bool
	// groundCoyoteLeft counts extra frames where jump is still allowed after losing floor contact.
	groundCoyoteLeft int32
	// jumpGrounded is updated at the end of ENTITY.UPDATE (strict floor + coyote); see entGrounded.
	jumpGrounded bool
	mass      float32
	friction  float32
	bounce    float32
	slide     bool
	pickMode  int32

	collided            bool
	otherID             int64
	hits                []int64
	hitPos              []rl.Vector3 // parallel to hits: contact point for rule-based collisions
	hitN                []rl.Vector3 // parallel to hits: world normal into source entity
	collType            int32
	hasHit              bool
	hitX, hitY, hitZ    float32
	hitNX, hitNY, hitNZ float32

	isSprite   bool
	spriteMode int32 // 1=y-billboard, 2=full-billboard, 3=static
	parentID   int64
	name       string
	// blenderTag: optional "tag" string from glTF extras (see MATERIAL.BULKASSIGN).
	blenderTag string

	rlModel    rl.Model
	hasRLModel bool
	loadPath   string
	// Loaded via rl.LoadModelAnimations (not embedded on rl.Model in raylib-go 0.56+).
	modelAnims []rl.ModelAnimation

	animIndex int32
	animTime  float32
	animSpeed float32
	animMode  int32 // 0–1=loop, 2=ping-pong, 3+=clamp/hold (legacy clamp used mode 1; use 3 now)
	animLen   float32

	// physBufIndex: index into PHYSICS3D shared matrix buffer (16 floats per body), or -1 if unset.
	physBufIndex int
	// collisionLayer: reserved for future Jolt object-layer / bitmask filtering (0–31); not yet applied to simulation.
	collisionLayer uint8

	// Skeletal bone socket (FindBone): host animated model + bone index; boneWorld is full world matrix each frame.
	boneHostID      int64
	boneIndex       int32
	boneWorld       rl.Matrix
	boneWorldValid  bool

	// Animation clip range (inclusive frame indices); -1 = use full animation.
	animClip0 int32
	animClip1 int32

	brushH heap.Handle // optional TagBrush; 0 = entity color/texture only

	// shadowCast: 0 = default (participate in shadow pass when renderer supports it), 1 = force cast, 2 = never cast (reserved).
	shadowCast int32

	// procMeshH: TagMeshBuilder for ENTITY.CREATEMESH procedural geometry (AddVertex / UpdateMesh).
	procMeshH heap.Handle

	// tween*: ENTITY.ANIMATETOWARD linear world-space lerp (advanced in ENTITY.UPDATE).
	tweenActive              bool
	tweenSX, tweenSY, tweenSZ float32
	tweenTX, tweenTY, tweenTZ float32
	tweenElapsed, tweenDuration float32
}

func newDefaultEnt(id int64) *ent {
	return &ent{
		id:    id,
		kind:  entKindEmpty,
		scale: rl.Vector3{X: 1, Y: 1, Z: 1},
		alpha: 1,
		r:     200, g: 200, b: 255,
		blendMode: -1,
		segH:      16, segV: 16,
		mass:         1,
		friction:     0.9,
		bounce:       0,
		animSpeed:    0,
		physBufIndex: -1,
		boneIndex:    -1,
		animClip0:    -1,
		animClip1:    -1,
	}
}
