//go:build cgo

package mbmatrix

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/vm/heap"
)

// matObj holds a 4×4 matrix (same heap tag as CAMERA.GETMATRIX / RAY.HITMESH transform handles).
type matObj struct {
	m rl.Matrix
}

func (o *matObj) TypeName() string { return "Matrix4" }

func (o *matObj) TypeTag() uint16 { return heap.TagMatrix }

func (o *matObj) Free() {}

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
