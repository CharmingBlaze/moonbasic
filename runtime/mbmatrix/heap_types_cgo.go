//go:build cgo

package mbmatrix

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/vm/heap"
)

// matObj holds a 4×4 matrix (same heap tag as CAMERA.GETMATRIX / RAY.HITMESH transform handles).
// No Raylib unload — pure value; Free is a no-op (idempotent).
type matObj struct {
	m rl.Matrix
}

func (o *matObj) TypeName() string { return "Matrix4" }

func (o *matObj) TypeTag() uint16 { return heap.TagMatrix }

func (o *matObj) Free() {}

// vec3Obj is a heap-wrapped Vector3; no native resource.
type vec3Obj struct {
	v rl.Vector3
}

func (o *vec3Obj) TypeName() string { return "Vec3" }

func (o *vec3Obj) TypeTag() uint16 { return heap.TagVec3 }

func (o *vec3Obj) Free() {}

type vec2Obj struct {
	v rl.Vector2
}

func (o *vec2Obj) TypeName() string { return "Vec2" }

func (o *vec2Obj) TypeTag() uint16 { return heap.TagVec2 }

func (o *vec2Obj) Free() {}

type quatObj struct {
	q rl.Quaternion
}

func (o *quatObj) TypeName() string { return "Quaternion" }

func (o *quatObj) TypeTag() uint16 { return heap.TagQuaternion }

func (o *quatObj) Free() {}
