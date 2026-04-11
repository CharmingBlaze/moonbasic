//go:build cgo || (windows && !cgo)

package mbentity

import (
	"moonbasic/runtime"
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
	entKindCapsule // MODEL.CREATECAPSULE — visual capsule (physics already used Jolt capsule)
)
	
type complexTween struct {
	prop     string
	start    float32
	target   float32
	elapsed  float32
	duration float32
	ease     string
}


// ent is a lightweight Blitz-style object (Essentails: ID, Transform, Appearance, Physics).
type ent struct {
	id   int64
	kind entKind

	// Transform (Backing SoA in entityStore)
	spatial *runtime.SpatialBuffer
	scale   rl.Vector3

	// Appearance & Rendering
	r, g, b          uint8
	alpha            float32
	hidden           bool
	shininess        float32
	texHandle        heap.Handle
	fxFlags          int32
	blendMode        int32 // -1 = default; else rl.BlendMode when drawing
	drawOrder        int32
	cullMode         int32 // 0=Auto (Frustum), 1=Force Visible, 2=Force Hidden

	// Geometry Bounds (Used for fast culling/picking)
	w, h, d    float32
	radius     float32
	segH, segV int32
	cylH       float32

	// Legacy Physics & State
	useSphere bool
	static    bool
	vel       rl.Vector3
	gravity   float32
	onGround  bool
	groundCoyoteLeft int32
	jumpGrounded bool
	mass      float32
	friction  float32
	bounce    float32
	pickMode  int32

	// Core Handles
	rlModel    rl.Model
	hasRLModel bool
	physBufIndex int
	// physicsDriven: Jolt-linked body owns motion; ENTITY.UPDATE must not apply scripted vel/gravity.
	physicsDriven bool
	collisionLayer uint8

	ext *entExt // Modular extensions (AI, Tweens, Animation State, Collision Results)
}

type entExt struct {
	// Metadata & Hierarchy
	name       string
	blenderTag string
	parentID   int64
	brushH     heap.Handle 
	procMeshH  heap.Handle

	// Animation State (Moved from core)
	loadPath   string
	modelAnims []rl.ModelAnimation
	animIndex  int32
	animTime   float32
	animSpeed  float32
	animMode   int32 
	animLen    float32
	animClip0  int32
	animClip1  int32

	// Bone / Attachment Data
	boneHostID      int64
	boneIndex       int32
	boneWorld       rl.Matrix
	boneWorldValid  bool

	// Collision Results (Bookkeeping)
	collided            bool
	otherID             int64
	hits                []int64
	hitPos              []rl.Vector3 
	hitN                []rl.Vector3 
	collType            int32
	hasHit              bool
	hitX, hitY, hitZ    float32
	hitNX, hitNY, hitNZ float32
	slide               bool

	// Tweens
	tweenActive              bool
	tweenSX, tweenSY, tweenSZ float32
	tweenTX, tweenTY, tweenTZ float32
	tweenElapsed, tweenDuration float32
	tweenFading              bool
	tweenAlphaStart          float32
	tweenAlphaEnd            float32
	tweenTurning             bool
	turnTargetX, turnTargetZ float32
	turnSpeed                float32
	tweenPulsing             bool
	pulseR1, pulseG1, pulseB1 uint8
	pulseR2, pulseG2, pulseB2 uint8
	pulseSpeed               float32
	pulseT                   float32
	complexTweens            []complexTween

	// Image Sequence / Sprite
	isSprite   bool
	spriteMode int32 
	seqH       heap.Handle
	seqFPS     float32
	seqTime    float32
	seqLoop    bool

	// AI State
	aiMode       string
	aiTarget     int64
	aiWaypoints  []rl.Vector3
	aiSpeed      float32
	aiIndex      int
	aiWanderCenter rl.Vector3
	aiWanderRadius float32
	onHitAction  string
	ghostMode    bool
	ghostTimer   float32

	// Polish
	outlineThickness float32
	outlineColor     rl.Color
	shadowCast       int32
}

func (e *ent) getExt() *entExt {
	if e.ext == nil {
		e.ext = &entExt{
			outlineColor: rl.Black,
		}
	}
	return e.ext
}

// parentID returns the parent entity id, or 0 if none. Does not allocate ext.
func (e *ent) parentID() int64 {
	if e.ext == nil {
		return 0
	}
	return e.ext.parentID
}

// boneWorldValid reports whether a bone-socket world matrix is current.
func (e *ent) boneWorldValid() bool {
	return e.ext != nil && e.ext.boneWorldValid
}

func (e *ent) getPos() rl.Vector3 {
	if e.spatial == nil { return rl.Vector3{} }
	return rl.Vector3{X: e.spatial.X[e.id], Y: e.spatial.Y[e.id], Z: e.spatial.Z[e.id]}
}

func (e *ent) setPos(v rl.Vector3) {
	if e.spatial == nil { return }
	e.spatial.X[e.id] = v.X
	e.spatial.Y[e.id] = v.Y
	e.spatial.Z[e.id] = v.Z
}

func (e *ent) getRot() (p, w, r float32) {
	if e.spatial == nil { return 0, 0, 0 }
	return e.spatial.P[e.id], e.spatial.W[e.id], e.spatial.R[e.id]
}

func (e *ent) setRot(p, w, r float32) {
	if e.spatial == nil { return }
	e.spatial.P[e.id] = p
	e.spatial.W[e.id] = w
	e.spatial.R[e.id] = r
}

func newDefaultEnt(id int64, s *runtime.SpatialBuffer) *ent {
	e := &ent{
		id:      id,
		spatial: s,
		kind:    entKindEmpty,
		scale:   rl.Vector3{X: 1, Y: 1, Z: 1},
		alpha:   1,
		r:       200, g: 200, b: 255,
		blendMode: -1,
		segH:      16, segV: 16,
		mass:         1,
		friction:     0.9,
		bounce:       0,
		physBufIndex: -1,
	}
	// Default transform in SoA
	e.setPos(rl.Vector3{})
	e.setRot(0, 0, 0)
	return e
}
