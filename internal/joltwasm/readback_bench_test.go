package joltwasm

import (
	"context"
	_ "embed"
	"os"
	"testing"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

//go:embed testdata/minimal.wasm
var minimalWasm []byte

// BenchmarkReadbackCopy simulates pulling a dense buffer from guest linear memory into a Go slice.
// Replace testdata/minimal.wasm with a Jolt wasm once the build pipeline exists; set JOLT_WASM to
// override the embedded fixture at test time.
func BenchmarkReadbackCopy(b *testing.B) {
	ctx := context.Background()
	r := wazero.NewRuntime(ctx)
	defer r.Close(ctx)
	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	wasm := minimalWasm
	if p := os.Getenv("JOLT_WASM"); p != "" {
		var err error
		wasm, err = os.ReadFile(p)
		if err != nil {
			b.Fatal(err)
		}
	}

	mod, err := r.Instantiate(ctx, wasm)
	if err != nil {
		b.Fatal(err)
	}
	mem := mod.Memory()
	if mem == nil || mem.Size() == 0 {
		b.Skip("fixture has no linear memory; use a wasm module that exports memory")
	}
	n := mem.Size()
	if n > 65536 {
		n = 65536
	}
	dst := make([]byte, n)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf, ok := mem.Read(0, n)
		if !ok {
			b.Fatal("read range")
		}
		copy(dst, buf)
	}
}

// BenchmarkReadbackView uses the wazero memory slice view only (no copy to host buffer).
// Hot path for physics SoA: iterate floats via unsafe slice from Read (see StateView.FloatsAfterHeader).
func BenchmarkReadbackView(b *testing.B) {
	ctx := context.Background()
	r := wazero.NewRuntime(ctx)
	defer r.Close(ctx)
	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	wasm := minimalWasm
	if p := os.Getenv("JOLT_WASM"); p != "" {
		var err error
		wasm, err = os.ReadFile(p)
		if err != nil {
			b.Fatal(err)
		}
	}

	mod, err := r.Instantiate(ctx, wasm)
	if err != nil {
		b.Fatal(err)
	}
	mem := mod.Memory()
	if mem == nil || mem.Size() == 0 {
		b.Skip("fixture has no linear memory")
	}
	n := mem.Size()
	if n > 65536 {
		n = 65536
	}

	var sink float32
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf, ok := mem.Read(0, n)
		if !ok {
			b.Fatal("read range")
		}
		f := FloatsView(buf)
		if len(f) > 0 {
			sink += f[0]
		}
	}
	_ = sink
}
