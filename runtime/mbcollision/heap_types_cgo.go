//go:build cgo

package mbcollision

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/vm/heap"
)

type rayObj struct {
	r rl.Ray
}

func (o *rayObj) TypeName() string { return "Ray" }

func (o *rayObj) TypeTag() uint16 { return heap.TagRay }

func (o *rayObj) Free() {}

type bboxObj struct {
	box rl.BoundingBox
}

func (o *bboxObj) TypeName() string { return "BoundingBox" }

func (o *bboxObj) TypeTag() uint16 { return heap.TagBBox }

func (o *bboxObj) Free() {}

type bsphereObj struct {
	center rl.Vector3
	radius float32
}

func (o *bsphereObj) TypeName() string { return "BoundingSphere" }

func (o *bsphereObj) TypeTag() uint16 { return heap.TagBSphere }

func (o *bsphereObj) Free() {}
