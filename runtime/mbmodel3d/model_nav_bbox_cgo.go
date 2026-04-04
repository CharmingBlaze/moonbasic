//go:build cgo

package mbmodel3d

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/vm/heap"
)

// ModelBoundingBoxForNav returns a world-oriented bounding box for Model or LODModel handles
// (used by navigation / obstacles). InstancedModel is not supported here.
func ModelBoundingBoxForNav(h *heap.Store, mh heap.Handle) (rl.BoundingBox, error) {
	if h == nil {
		return rl.BoundingBox{}, fmt.Errorf("ModelBoundingBoxForNav: nil heap")
	}
	if mo, err := heap.Cast[*modelObj](h, mh); err == nil {
		return rl.GetModelBoundingBox(mo.model), nil
	}
	if lo, err := heap.Cast[*lodModelObj](h, mh); err == nil {
		if lo.models[0].MeshCount == 0 {
			return rl.BoundingBox{}, fmt.Errorf("ModelBoundingBoxForNav: LOD model has no meshes")
		}
		m := lo.models[0]
		m.Transform = lo.transform
		return rl.GetModelBoundingBox(m), nil
	}
	return rl.BoundingBox{}, fmt.Errorf("ModelBoundingBoxForNav: handle must be Model or LODModel")
}
