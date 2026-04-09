package vm

import (
	"bytes"
	"strings"
	"testing"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/opcode"
)

func TestTraceOpcodeSequence(t *testing.T) {
	var buf bytes.Buffer
	h := heap.New()
	reg := runtime.NewRegistry(h)
	reg.InitCore()
	v := New(reg, h)
	v.Trace = true
	v.TraceOut = &buf

	p := opcode.NewProgram()
	// Emit(op, dst, srcA, srcB, operand, line)
	p.Main.Emit(opcode.OpPushInt, 0, 0, 0, p.Main.AddInt(2), 1)
	p.Main.Emit(opcode.OpPushInt, 1, 0, 0, p.Main.AddInt(3), 1)
	p.Main.Emit(opcode.OpAdd, 2, 0, 1, 0, 1) // R2 = R0 + R1
	p.Main.Emit(opcode.OpHalt, 0, 0, 0, 0, 1)

	if err := v.Execute(p); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "[trace]") || !strings.Contains(out, "PUSH_INT") {
		t.Fatalf("expected trace lines, got:\n%s", out)
	}
	if strings.Count(out, "[trace]") < 4 {
		t.Fatalf("expected trace per opcode, got %d lines", strings.Count(out, "[trace]"))
	}
	reg.Shutdown()
}
