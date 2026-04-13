package vm

import (
	"testing"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/opcode"
)

// BenchmarkVM_RegisterArithmetic measures the execute loop for a tight chain of OpAdd + OpHalt.
// Profiling: go test -bench=BenchmarkVM_RegisterArithmetic -benchmem -cpuprofile=vm.prof ./vm/
// Then: go tool pprof -web vm.prof
func BenchmarkVM_RegisterArithmetic(b *testing.B) {
	const nOps = 512
	p := opcode.NewProgram()
	c := p.Main
	c.Emit(opcode.OpPushInt, 0, 0, 0, c.AddInt(7), 1)
	c.Emit(opcode.OpPushInt, 1, 0, 0, c.AddInt(3), 1)
	for i := 0; i < nOps; i++ {
		c.Emit(opcode.OpAdd, 2, 0, 1, 0, 1)
	}
	c.Emit(opcode.OpHalt, 0, 0, 0, 0, 1)

	h := heap.New()
	reg := runtime.NewRegistryHeadless(h)
	reg.InitCore()
	v := New(reg, h)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := v.Execute(p); err != nil {
			b.Fatal(err)
		}
	}
	reg.Shutdown()
}
