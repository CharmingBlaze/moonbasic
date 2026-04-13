//go:build cgo || (windows && !cgo)

package mbmatrix

import (
	"moonbasic/hal"
	"moonbasic/vm/heap"
)

// matObj holds a 4×4 matrix (same heap tag as CAMERA.GETMATRIX / RAY.HITMESH transform handles).
// No native resource — pure value; Free is a no-op (idempotent).
type matObj struct {
	m hal.Matrix
}

func (o *matObj) TypeName() string { return "Matrix4" }
func (o *matObj) TypeTag() uint16 { return heap.TagMatrix }
func (o *matObj) Free() {}

// vec3Obj is a heap-wrapped Vector3; no native resource.
type vec3Obj struct {
	v hal.V3
}

func (o *vec3Obj) TypeName() string { return "Vec3" }
func (o *vec3Obj) TypeTag() uint16 { return heap.TagVec3 }
func (o *vec3Obj) Free() {}

// vec2Obj is a heap-wrapped Vector2; no native resource.
type vec2Obj struct {
	v hal.V2
}

func (o *vec2Obj) TypeName() string { return "Vec2" }
func (o *vec2Obj) TypeTag() uint16 { return heap.TagVec2 }
func (o *vec2Obj) Free() {}

// quatObj is a heap-wrapped Quaternion (V4); no native resource.
type quatObj struct {
	q hal.V4
}

func (o *quatObj) TypeName() string { return "Quaternion" }
func (o *quatObj) TypeTag() uint16 { return heap.TagQuaternion }
func (o *quatObj) Free() {}
