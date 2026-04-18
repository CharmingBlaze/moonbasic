//go:build cgo || (windows && !cgo)

package water

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/vm/heap"
)

// WaterObject is a finite water plane with simple wave animation state.
// Owns: Mesh (GenMeshPlane), Mat (LoadMaterialDefault) — both released in Free.
type WaterObject struct {
	Mesh     rl.Mesh
	Mat      rl.Material
	Width    float32
	Depth    float32
	PX, PY, PZ float32
	RotX, RotY, RotZ float32
	ScaleX, ScaleY, ScaleZ float32
	WaveT    float32
	WaveAmp  float32
	WaveFreq float32
	Shallow  rl.Color
	Deep     rl.Color
	BedY     float32
	freed    bool
}

func (w *WaterObject) TypeName() string { return "Water" }
func (w *WaterObject) TypeTag() uint16  { return heap.TagWater }

func (w *WaterObject) Free() {
	if w.freed {
		return
	}
	if w.Mesh.TriangleCount > 0 {
		rl.UnloadMesh(&w.Mesh)
	}
	rl.UnloadMaterial(w.Mat)
	w.Mesh = rl.Mesh{}
	w.freed = true
}
