//go:build cgo || (windows && !cgo)

package mbentity

import (
	"fmt"

	texmod "moonbasic/runtime/texture"
	mbimage "moonbasic/runtime/mbimage"
	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerEntitySpriteAnimAPI(m *Module, r runtime.Registrar) {
	r.Register("ENTITY.SETSPRITEFRAME", "entity", runtime.AdaptLegacy(m.entSetSpriteFrame))
	r.Register("ENTITY.SETANIMATION", "entity", runtime.AdaptLegacy(m.entSetAnimation))
}

// DuplicateEntityAt copies a template entity and moves it to world (x,y,z). Used by TERRAIN.APPLYTILES.
func (m *Module) DuplicateEntityAt(templateID int64, x, y, z float32) (int64, error) {
	v, err := m.entCopy([]value.Value{value.FromInt(templateID)})
	if err != nil {
		return 0, err
	}
	nid, ok := v.ToInt()
	if !ok || nid < 1 {
		return 0, fmt.Errorf("DuplicateEntityAt: bad new id")
	}
	_, err = m.entSetPosition([]value.Value{
		value.FromInt(nid),
		value.FromFloat(float64(x)),
		value.FromFloat(float64(y)),
		value.FromFloat(float64(z)),
	})
	if err != nil {
		return 0, err
	}
	return nid, nil
}

func (m *Module) entSetSpriteFrame(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("ENTITY.SETSPRITEFRAME: heap not bound")
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.SETSPRITEFRAME expects (entity, frameIndex)")
	}
	eid, ok := m.entID(args[0])
	if !ok || eid < 1 {
		return value.Nil, fmt.Errorf("ENTITY.SETSPRITEFRAME: invalid entity")
	}
	e := m.store().ents[eid]
	if e == nil || e.ext == nil || !e.ext.isSprite || e.texHandle == 0 {
		return value.Nil, fmt.Errorf("ENTITY.SETSPRITEFRAME: not a billboard with texture")
	}
	fi, ok2 := args[1].ToInt()
	if !ok2 {
		return value.Nil, fmt.Errorf("ENTITY.SETSPRITEFRAME: frameIndex must be numeric")
	}
	to, err := heap.Cast[*texmod.TextureObject](m.h, e.texHandle)
	if err != nil {
		return value.Nil, fmt.Errorf("ENTITY.SETSPRITEFRAME: entity texture is not an atlas (use TEXTURE.LOAD / TEXTURE.LOADANIM)")
	}
	// to.mu.Lock()
	to.FrameIndex = int32(fi) // Warning: FrameIndex concurrent assignment unprotected
	// to.mu.Unlock()
	return value.Nil, nil
}

func (m *Module) entSetAnimation(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("ENTITY.SETANIMATION: heap not bound")
	}
	if len(args) != 3 && len(args) != 4 {
		return value.Nil, fmt.Errorf("ENTITY.SETANIMATION expects (entity, imageSequenceHandle, fps [, loop])")
	}
	eid, ok := m.entID(args[0])
	if !ok || eid < 1 {
		return value.Nil, fmt.Errorf("ENTITY.SETANIMATION: invalid entity")
	}
	e := m.store().ents[eid]
	if e == nil || e.ext == nil || !e.ext.isSprite {
		return value.Nil, fmt.Errorf("ENTITY.SETANIMATION: not a sprite/billboard")
	}
	if args[1].Kind == value.KindInt && args[1].IVal == 0 {
		ext := e.getExt()
		ext.seqH = 0
		ext.seqFPS = 0
		ext.seqTime = 0
		ext.seqLoop = true
		return value.Nil, nil
	}
	if args[1].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("ENTITY.SETANIMATION: image sequence handle (or 0 to clear)")
	}
	if args[1].IVal == 0 {
		ext := e.getExt()
		ext.seqH = 0
		ext.seqFPS = 0
		ext.seqTime = 0
		ext.seqLoop = true
		return value.Nil, nil
	}
	_, err := heap.Cast[*mbimage.ImageSequence](m.h, heap.Handle(args[1].IVal))
	if err != nil {
		return value.Nil, fmt.Errorf("ENTITY.SETANIMATION: handle must be IMAGE.LOADSEQUENCE / IMAGE.LOADGIF")
	}
	fps, okf := args[2].ToFloat()
	if !okf {
		if fi, ok := args[2].ToInt(); ok {
			fps = float64(fi)
			okf = true
		}
	}
	if !okf || fps <= 0 {
		return value.Nil, fmt.Errorf("ENTITY.SETANIMATION: fps must be > 0")
	}
	loop := true
	if len(args) == 4 {
		if args[3].Kind == value.KindBool {
			loop = args[3].IVal != 0
		} else if bi, ok := args[3].ToInt(); ok {
			loop = bi != 0
		}
	}
	ext := e.getExt()
	ext.seqH = heap.Handle(args[1].IVal)
	ext.seqFPS = float32(fps)
	ext.seqTime = 0
	ext.seqLoop = loop
	return value.Nil, nil
}

func (m *Module) advanceSpriteImageSequences(dt float32) {
	if m.h == nil {
		return
	}
	st := m.store()
	for _, e := range st.ents {
		if e == nil || e.ext == nil || !e.ext.isSprite || e.ext.seqH == 0 || e.texHandle == 0 {
			continue
		}
		ext := e.ext
		seq, err := heap.Cast[*mbimage.ImageSequence](m.h, ext.seqH)
		if err != nil || seq.FrameCount() == 0 {
			continue
		}
		to, err := heap.Cast[*texmod.TextureObject](m.h, e.texHandle)
		if err != nil {
			continue
		}
		n := seq.FrameCount()
		ext.seqTime += dt
		frameDur := float32(1.0) / ext.seqFPS
		if frameDur <= 0 {
			frameDur = 0.1
		}
		idx := int(ext.seqTime / frameDur)
		if ext.seqLoop {
			idx %= n
		} else if idx >= n {
			idx = n - 1
		}
		im := seq.Frame(idx)
		if im == nil {
			continue
		}
		_ = texmod.SyncTextureFromImage(to, im)
	}
}
