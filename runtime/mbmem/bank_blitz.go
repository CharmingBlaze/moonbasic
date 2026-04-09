package mbmem

import (
	"moonbasic/runtime"
	"moonbasic/vm/value"
)

// registerBankBlitzAliases registers Blitz-style bank names (same handles as MEM.*).
func registerBankBlitzAliases(m *Module, r runtime.Registrar) {
	r.Register("CreateBank", "mem", runtime.AdaptLegacy(m.memMake))
	r.Register("FreeBank", "mem", runtime.AdaptLegacy(m.memFree))
	r.Register("BankSize", "mem", runtime.AdaptLegacy(m.memSize))
	r.Register("ResizeBank", "mem", runtime.AdaptLegacy(m.memResize))
	r.Register("CopyBank", "mem", runtime.AdaptLegacy(m.bankCopyBlitzOrder))
	r.Register("PeekByte", "mem", runtime.AdaptLegacy(m.memGetByte))
	r.Register("PokeByte", "mem", runtime.AdaptLegacy(m.memSetByte))
	r.Register("PeekShort", "mem", runtime.AdaptLegacy(m.memGetWord))
	r.Register("PokeShort", "mem", runtime.AdaptLegacy(m.memSetWord))
	r.Register("PeekInt", "mem", runtime.AdaptLegacy(m.memGetDword))
	r.Register("PokeInt", "mem", runtime.AdaptLegacy(m.memSetDword))
	r.Register("PeekFloat", "mem", runtime.AdaptLegacy(m.memGetDouble))
	r.Register("PokeFloat", "mem", runtime.AdaptLegacy(m.memSetDouble))
}

// bankCopyBlitzOrder is CopyBank(src, srcOffset, dest, destOffset, count) — forwards to MEM.COPY(src, dst, ...).
func (m *Module) bankCopyBlitzOrder(args []value.Value) (value.Value, error) {
	if len(args) != 5 {
		return value.Nil, runtime.Errorf("CopyBank expects (src, srcOffset, dest, destOffset, count)")
	}
	// MEM.COPY(src, dst, srcOff, dstOff, size)
	reordered := []value.Value{args[0], args[2], args[1], args[3], args[4]}
	return m.memCopy(reordered)
}
