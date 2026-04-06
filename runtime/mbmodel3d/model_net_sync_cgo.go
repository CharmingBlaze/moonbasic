//go:build cgo || (windows && !cgo)

package mbmodel3d

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/vm/heap"
)

// ModelTranslationForSync returns world translation from a Model or LODModel (for network sync).
func ModelTranslationForSync(h *heap.Store, mh heap.Handle) (x, y, z float32, err error) {
	if h == nil {
		return 0, 0, 0, fmt.Errorf("ModelTranslationForSync: nil heap")
	}
	if mo, e := heap.Cast[*modelObj](h, mh); e == nil {
		t := mo.model.Transform
		return t.M12, t.M13, t.M14, nil
	}
	if lo, e := heap.Cast[*lodModelObj](h, mh); e == nil {
		t := lo.transform
		return t.M12, t.M13, t.M14, nil
	}
	return 0, 0, 0, fmt.Errorf("ModelTranslationForSync: handle must be Model or LODModel")
}

// SetModelTranslation applies translation to Model or LODModel (client-side sync apply).
func SetModelTranslation(h *heap.Store, mh heap.Handle, x, y, z float32) error {
	if h == nil {
		return fmt.Errorf("SetModelTranslation: nil heap")
	}
	if mo, e := heap.Cast[*modelObj](h, mh); e == nil {
		mo.model.Transform = rl.MatrixTranslate(x, y, z)
		return nil
	}
	if lo, e := heap.Cast[*lodModelObj](h, mh); e == nil {
		lo.transform = rl.MatrixTranslate(x, y, z)
		return nil
	}
	return fmt.Errorf("SetModelTranslation: handle must be Model or LODModel")
}
