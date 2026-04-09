//go:build cgo || (windows && !cgo)

package mbentity

import (
	"fmt"
	"math"
	"strings"
	"unsafe"

	"moonbasic/runtime"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// animMode: 0–1 = loop, 2 = ping-pong, 3+ = clamp/hold at end (legacy scripts used 1 for clamp; use 3 now).
func pickAnimFrame(e *ent, anim rl.ModelAnimation) int32 {
	if anim.FrameCount <= 0 {
		return 0
	}
	fc := anim.FrameCount
	lo, hi := int32(0), fc-1
	if e.animClip0 >= 0 {
		lo = e.animClip0
		hi = e.animClip1
		if hi < lo {
			hi = lo
		}
		if lo < 0 {
			lo = 0
		}
		if hi >= fc {
			hi = fc - 1
		}
	}
	span := hi - lo + 1
	if span <= 1 {
		return lo
	}

	mode := e.animMode
	switch {
	case mode == 2:
		// Ping-pong over [lo..hi]
		period := float32(2 * (span - 1))
		if period <= 0 {
			return lo
		}
		phase := float32(math.Mod(float64(e.animTime), float64(period)))
		if phase < 0 {
			phase += period
		}
		var local int32
		if phase >= float32(span-1) {
			local = int32(period - phase)
		} else {
			local = int32(phase)
		}
		if local < 0 {
			local = 0
		}
		if local >= span {
			local = span - 1
		}
		return lo + local
	case mode >= 3:
		off := int32(e.animTime)
		if off < 0 {
			off = 0
		}
		if off >= span {
			off = span - 1
		}
		return lo + off
	default:
		// 0,1: loop
		spanF := float32(span)
		t := float32(math.Mod(float64(e.animTime), float64(spanF)))
		if t < 0 {
			t += spanF
		}
		off := int32(t)
		if off >= span {
			off = span - 1
		}
		return lo + off
	}
}

func boneNameStr(name [32]int8) string {
	b := make([]byte, 0, 32)
	for i := 0; i < len(name) && name[i] != 0; i++ {
		b = append(b, byte(name[i]))
	}
	return string(b)
}

func (m *Module) syncBoneSockets() {
	st := m.store()
	for _, e := range st.ents {
		if e == nil || e.boneIndex < 0 || e.boneHostID < 1 {
			continue
		}
		host := st.ents[e.boneHostID]
		if host == nil || !host.hasRLModel {
			e.boneWorldValid = false
			continue
		}
		meshes := host.rlModel.GetMeshes()
		if len(meshes) == 0 {
			e.boneWorldValid = false
			continue
		}
		mesh := meshes[0]
		bi := int(e.boneIndex)
		if mesh.BoneMatrices == nil || bi < 0 || int(mesh.BoneCount) <= bi {
			e.boneWorldValid = false
			continue
		}
		bm := unsafe.Slice(mesh.BoneMatrices, mesh.BoneCount)[bi]
		hw := m.worldMatrix(host)
		// bone matrix is in model space; host world × model-space bone = world bone transform
		e.boneWorld = rl.MatrixMultiply(hw, bm)
		e.boneWorldValid = true
	}
}

func (m *Module) entFindBone(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("FindBone expects (entity#, name$)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("FindBone: invalid entity")
	}
	host := m.store().ents[id]
	if host == nil || !host.hasRLModel {
		return value.Nil, fmt.Errorf("FindBone: entity has no model")
	}
	name, err := rt.ArgString(args, 1)
	if err != nil || name == "" {
		return value.Nil, fmt.Errorf("FindBone: name required")
	}
	bones := host.rlModel.GetBones()
	var bi int32 = -1
	for i := range bones {
		if strings.EqualFold(boneNameStr(bones[i].Name), name) {
			bi = int32(i)
			break
		}
	}
	if bi < 0 {
		return value.Nil, fmt.Errorf("FindBone: no bone %q", name)
	}
	st := m.store()
	nid := st.nextID
	st.nextID++
	e := newDefaultEnt(nid)
	e.kind = entKindEmpty
	e.hidden = true
	e.static = true
	e.boneHostID = id
	e.boneIndex = bi
	e.boneWorldValid = false
	st.ents[nid] = e
	return value.FromInt(nid), nil
}

func (m *Module) entExtractAnimSeq(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("ENTITY.EXTRACTANIMSEQ expects (entity#, startFrame#, endFrame#)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	s0, ok1 := args[1].ToInt()
	s1, ok2 := args[2].ToInt()
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("frames must be numeric")
	}
	e.animClip0 = int32(s0)
	e.animClip1 = int32(s1)
	return value.Nil, nil
}

func (m *Module) entSetAnimIndex(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.SETANIMINDEX expects (entity#, animIndex#)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	ai, ok := args[1].ToInt()
	if !ok || ai < 0 {
		return value.Nil, fmt.Errorf("anim index must be >= 0")
	}
	e.animIndex = int32(ai)
	e.animTime = 0
	return value.Nil, nil
}
