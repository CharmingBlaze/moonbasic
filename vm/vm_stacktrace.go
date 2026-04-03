package vm

import (
	"fmt"
	"strings"
)

// FormatCallStack returns a multi-line description of the moonBASIC call stack
// (innermost frame first). Safe to call only while the VM is executing a program.
func (v *VM) FormatCallStack() string {
	if v.CallStack == nil || v.CallStack.Depth() == 0 {
		return "(empty call stack)\n"
	}
	var b strings.Builder
	frames := v.CallStack.FramesCopy()
	for i := len(frames) - 1; i >= 0; i-- {
		fr := frames[i]
		ip := fr.IP - 1
		if ip < 0 {
			ip = 0
		}
		line := -1
		if ip < len(fr.Chunk.SourceLines) {
			line = int(fr.Chunk.SourceLines[ip])
		}
		fmt.Fprintf(&b, "  %s line %d\n", fr.Chunk.Name, line)
	}
	return b.String()
}
