package water

import (
	"moonbasic/runtime"
)

func (m *Module) Register(r runtime.Registrar) {
	r.Register("WATER.MAKE", "water", runtime.AdaptLegacy(m.wMake))
	r.Register("WATER.CREATE", "water", runtime.AdaptLegacy(m.wCreate))
	r.Register("WATER.FREE", "water", runtime.AdaptLegacy(m.wFree))
	r.Register("WATER.SETPOS", "water", runtime.AdaptLegacy(m.wSetPos))
	r.Register("WATER.DRAW", "water", runtime.AdaptLegacy(m.wDraw))
	r.Register("WATER.UPDATE", "water", runtime.AdaptLegacy(m.wUpdate))
	r.Register("WATER.SETWAVEHEIGHT", "water", runtime.AdaptLegacy(m.wSetWaveHeight))
	r.Register("WATER.SETWAVE", "water", runtime.AdaptLegacy(m.wSetWave))
	r.Register("WATER.GETWAVEY", "water", runtime.AdaptLegacy(m.wGetWaveY))
	r.Register("WATER.GETDEPTH", "water", runtime.AdaptLegacy(m.wGetDepth))
	r.Register("WATER.ISUNDER", "water", runtime.AdaptLegacy(m.wIsUnder))
	r.Register("WATER.SETSHALLOWCOLOR", "water", runtime.AdaptLegacy(m.wSetShallow))
	r.Register("WATER.SETDEEPCOLOR", "water", runtime.AdaptLegacy(m.wSetDeep))
	r.Register("WATER.SETCOLOR", "water", runtime.AdaptLegacy(m.wSetColor))
}
