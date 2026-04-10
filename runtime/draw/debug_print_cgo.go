//go:build cgo || (windows && !cgo)

package mbdraw

import (
	"fmt"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

var (
	debugPrintLastFrame uint64
	debugPrintY         float32 = 10
)

func registerDebugPrint(m *Module, r runtime.Registrar) {
	dbg := func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) < 1 {
			return value.Nil, fmt.Errorf("DEBUG.PRINT expects (template$, ...)")
		}
		tmpl := stringFromRT(rt, args[0])
		n := len(args) - 1
		if n > 10 {
			return value.Nil, fmt.Errorf("DEBUG.PRINT: at most 10 placeholders {0}..{9}")
		}
		pairs := make([]string, 0, n*2)
		for i := 0; i < n; i++ {
			pairs = append(pairs, fmt.Sprintf("{%d}", i))
			pairs = append(pairs, valueStringDebug(rt, args[i+1]))
		}
		msg := strings.NewReplacer(pairs...).Replace(tmpl)
		if rt.FrameCount != debugPrintLastFrame {
			debugPrintLastFrame = rt.FrameCount
			debugPrintY = 10
		}
		x := int32(12)
		y := int32(debugPrintY)
		sz := int32(14)
		rl.DrawText(msg, x, y, sz, rl.RayWhite)
		debugPrintY += float32(sz) + 4
		return value.Nil, nil
	}
	r.Register("DEBUG.PRINT", "draw", dbg)
}

func valueStringDebug(rt *runtime.Runtime, v value.Value) string {
	if v.Kind == value.KindString {
		return stringFromRT(rt, v)
	}
	return v.String()
}
