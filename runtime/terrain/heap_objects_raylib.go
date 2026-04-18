//go:build (cgo || (windows && !cgo)) && (!windows || !gopls_stub)

package terrain

import (
	"math"
	"sync/atomic"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/vm/heap"
)

// TerrainObject is a heightfield with optional chunk meshes for rendering.
type TerrainObject struct {
	WorldW    int
	WorldH    int
	CellSize  float32
	ChunkSize int
	Heights   []float32

	ChunkW int
	ChunkH int

	PX, PY, PZ float32
	RotX, RotY, RotZ float32

	// ScaleX/ScaleZ stretch world spacing per height sample; ScaleY scales vertical range (default 1).
	ScaleX, ScaleY, ScaleZ float32
	// DetailFactor in (0,1]: lower values build coarser chunk meshes (LOD / performance).
	DetailFactor float32

	DiffuseTex   rl.Texture2D
	DiffuseLoaded bool
	SplatImg     *rl.Image // retained for TERRAIN.GETSPLAT (diffuse/splat CPU sample); optional

	Chunks []chunkSlot

	StreamEnabled bool
	LoadDist      float32
	UnloadDist    float32
	CenterX       float32
	CenterZ       float32

	MaxHeight float32

	// Async chunk mesh build (mesh_jobs.go): goroutine posts chunkMeshJob; main thread drains meshJobs.
	meshJobs       chan chunkMeshJob
	meshJobsInflight atomic.Int32
	// MeshBuildBudgetPerTick limits rebuilds per WORLD.UPDATE tick (0 = unlimited).
	MeshBuildBudgetPerTick int
	MeshBuildAsync         bool

	freed bool
}

// chunkSlot owns GPU mesh and default material for one terrain chunk (rebuildChunkMesh).
type chunkSlot struct {
	Mesh       rl.Mesh
	Mat        rl.Material
	Loaded     bool
	Dirty      bool
	CX, CZ     int
	LastUpload int64
	// MinH/MaxH are heightfield samples (before TerrainObject.PY); set in rebuildChunkMesh.
	MinH, MaxH float32
	BoundsValid bool
	// PendingAsync: CPU heightmap prep in flight; main thread must drain before rebuilding.
	PendingAsync bool
}

func (t *TerrainObject) TypeName() string { return "Terrain" }
func (t *TerrainObject) TypeTag() uint16  { return heap.TagTerrain }

// Free releases GPU meshes and height data.
func (t *TerrainObject) Free() {
	if t.freed {
		return
	}
	t.freed = true
	t.shutdownMeshJobQueue()
	for i := range t.Chunks {
		ch := &t.Chunks[i]
		if ch.Loaded {
			rl.UnloadMaterial(ch.Mat)
			rl.UnloadMesh(&ch.Mesh)
			ch.Loaded = false
		}
	}
	if t.DiffuseLoaded {
		rl.UnloadTexture(t.DiffuseTex)
		t.DiffuseLoaded = false
	}
	if t.SplatImg != nil {
		rl.UnloadImage(t.SplatImg)
		t.SplatImg = nil
	}
	t.Chunks = nil
	t.Heights = nil
}

func max32Terrain(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func idx2(t *TerrainObject, cx, cz int) int {
	return cz*t.ChunkW + cx
}

func cellIndex(t *TerrainObject, wx, wz int) int {
	return wz*t.WorldW + wx
}

func (t *TerrainObject) heightAtCell(wx, wz int) float32 {
	if wx < 0 || wz < 0 || wx >= t.WorldW || wz >= t.WorldH {
		return 0
	}
	return t.Heights[cellIndex(t, wx, wz)]
}

func (t *TerrainObject) scaleXEff() float32 {
	if t.ScaleX <= 0 {
		return 1
	}
	return t.ScaleX
}
func (t *TerrainObject) scaleYEff() float32 {
	if t.ScaleY <= 0 {
		return 1
	}
	return t.ScaleY
}
func (t *TerrainObject) scaleZEff() float32 {
	if t.ScaleZ <= 0 {
		return 1
	}
	return t.ScaleZ
}

// HeightWorld returns bilinear height at world xz (relative to origin px,pz).
func (t *TerrainObject) HeightWorld(x, z float32) float32 {
	if t.WorldW < 2 || t.WorldH < 2 {
		return 0
	}
	sx := t.scaleXEff()
	sz := t.scaleZEff()
	sy := t.scaleYEff()
	lx := (x - t.PX) / (t.CellSize * sx)
	lz := (z - t.PZ) / (t.CellSize * sz)
	if lx < 0 || lz < 0 || lx >= float32(t.WorldW-1) || lz >= float32(t.WorldH-1) {
		return 0
	}
	x0 := int(lx)
	z0 := int(lz)
	fx := lx - float32(x0)
	fz := lz - float32(z0)
	h00 := t.heightAtCell(x0, z0)
	h10 := t.heightAtCell(x0+1, z0)
	h01 := t.heightAtCell(x0, z0+1)
	h11 := t.heightAtCell(x0+1, z0+1)
	a := h00*(1-fx) + h10*fx
	b := h01*(1-fx) + h11*fx
	raw := a*(1-fz) + b*fz
	return raw*sy + t.PY
}

// GridXZ returns fractional grid coordinates used for height/splat (same as HeightWorld).
func (t *TerrainObject) GridXZ(x, z float32) (lx, lz float32, ok bool) {
	if t.WorldW < 2 || t.WorldH < 2 {
		return 0, 0, false
	}
	sx := t.scaleXEff()
	sz := t.scaleZEff()
	lx = (x - t.PX) / (t.CellSize * sx)
	lz = (z - t.PZ) / (t.CellSize * sz)
	if lx < 0 || lz < 0 || lx >= float32(t.WorldW-1) || lz >= float32(t.WorldH-1) {
		return lx, lz, false
	}
	return lx, lz, true
}

// NormalWorld returns a unit up-ish normal for slope/tilt (Y typically positive).
func (t *TerrainObject) NormalWorld(x, z float32) (nx, ny, nz float32) {
	dx := t.CellSize * t.scaleXEff() * 0.5
	dz := t.CellSize * t.scaleZEff() * 0.5
	if dx < 1e-4 {
		dx = 1
	}
	if dz < 1e-4 {
		dz = 1
	}
	hL := t.HeightWorld(x-dx, z)
	hR := t.HeightWorld(x+dx, z)
	hD := t.HeightWorld(x, z-dz)
	hU := t.HeightWorld(x, z+dz)
	dhdx := (hR - hL) / (2 * dx)
	dhdz := (hU - hD) / (2 * dz)
	nx = -dhdx
	ny = 1
	nz = -dhdz
	inv := float32(1.0 / math.Sqrt(float64(nx*nx+ny*ny+nz*nz)))
	return nx * inv, ny * inv, nz * inv
}

// SplatAt returns a surface id 0..255 from the splat/diffuse image red channel (or -1 if none).
func (t *TerrainObject) SplatAt(x, z float32) int32 {
	if t.SplatImg == nil || t.SplatImg.Data == nil {
		return -1
	}
	lx, lz, ok := t.GridXZ(x, z)
	if !ok {
		return -1
	}
	w := float32(t.SplatImg.Width - 1)
	h := float32(t.SplatImg.Height - 1)
	if w <= 0 || h <= 0 {
		return -1
	}
	u := lx / float32(t.WorldW-1)
	v := lz / float32(t.WorldH-1)
	if u < 0 {
		u = 0
	}
	if v < 0 {
		v = 0
	}
	if u > 1 {
		u = 1
	}
	if v > 1 {
		v = 1
	}
	px := u * w
	py := v * h
	x0 := int(px)
	y0 := int(py)
	x1 := x0 + 1
	y1 := y0 + 1
	if x1 >= int(t.SplatImg.Width) {
		x1 = int(t.SplatImg.Width) - 1
	}
	if y1 >= int(t.SplatImg.Height) {
		y1 = int(t.SplatImg.Height) - 1
	}
	fx := px - float32(x0)
	fy := py - float32(y0)
	c00 := rl.GetImageColor(*t.SplatImg, int32(x0), int32(y0))
	c10 := rl.GetImageColor(*t.SplatImg, int32(x1), int32(y0))
	c01 := rl.GetImageColor(*t.SplatImg, int32(x0), int32(y1))
	c11 := rl.GetImageColor(*t.SplatImg, int32(x1), int32(y1))
	r00 := float32(c00.R)
	r10 := float32(c10.R)
	r01 := float32(c01.R)
	r11 := float32(c11.R)
	a := r00*(1-fx) + r10*fx
	b := r01*(1-fx) + r11*fx
	r := a*(1-fy) + b*fy
	return int32(r + 0.5)
}

// RaycastTerrain intersects a ray with the heightfield surface; returns hit and world position.
func (t *TerrainObject) RaycastTerrain(sx, sy, sz, dx, dy, dz, maxDist float32) (hit bool, hx, hy, hz float32) {
	if maxDist <= 0 {
		maxDist = 1e5
	}
	ln := float32(math.Sqrt(float64(dx*dx + dy*dy + dz*dz)))
	if ln < 1e-8 {
		return false, 0, 0, 0
	}
	dx /= ln
	dy /= ln
	dz /= ln
	step := t.CellSize * 0.25
	if sxe := t.scaleXEff(); sxe > 0 {
		step *= sxe
	}
	if step < 0.05 {
		step = 0.05
	}
	if step > 5 {
		step = 5
	}
	var prevDist float32 = -1
	for dist := float32(0); dist <= maxDist; dist += step {
		px := sx + dx*dist
		py := sy + dy*dist
		pz := sz + dz*dist
		g := t.HeightWorld(px, pz)
		if prevDist >= 0 {
			ppx := sx + dx*prevDist
			ppy := sy + dy*prevDist
			ppz := sz + dz*prevDist
			pg := t.HeightWorld(ppx, ppz)
			if ppy > pg && py <= g {
				t0 := prevDist
				t1 := dist
				for i := 0; i < 10; i++ {
					tm := (t0 + t1) * 0.5
					mx := sx + dx*tm
					my := sy + dy*tm
					mz := sz + dz*tm
					mg := t.HeightWorld(mx, mz)
					if my > mg {
						t0 = tm
					} else {
						t1 = tm
					}
				}
				tHit := (t0 + t1) * 0.5
				hx = sx + dx*tHit
				hz = sz + dz*tHit
				hy = t.HeightWorld(hx, hz)
				return true, hx, hy, hz
			}
		}
		prevDist = dist
	}
	return false, 0, 0, 0
}

// SlopeDeg approximate slope angle at world position (degrees).
func (t *TerrainObject) SlopeDeg(x, z float32) float32 {
	d := t.CellSize * 0.5 * max32Terrain(t.scaleXEff(), t.scaleZEff())
	if d < 1e-4 {
		d = 1
	}
	hL := t.HeightWorld(x-d, z)
	hR := t.HeightWorld(x+d, z)
	hD := t.HeightWorld(x, z-d)
	hU := t.HeightWorld(x, z+d)
	dhdx := (hR - hL) / (2 * d)
	dhdz := (hU - hD) / (2 * d)
	grad := math.Sqrt(float64(dhdx*dhdx + dhdz*dhdz))
	return float32(math.Atan(grad) * 180 / math.Pi)
}
