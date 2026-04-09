# Jolt 3D via WASM (prototype notes)

This repo’s production 3D path uses [jolt-go](https://github.com/bbitechnologies/jolt-go) with CGO on Linux (`runtime/physics3d/jolt_*_linux.go`, build tag **`linux && cgo`**). A zero-CGO alternative is:

1. **Compile** a minimal Jolt (or a thin C ABI façade) to **WebAssembly** with fixed calling conventions and exported **linear memory** for body transforms (SoA layout aligns with `docs/architecture/PHYSICS_ARENA_AND_INTERPOLATION.md`).
2. **Host** the module with [wazero](https://github.com/tetratelabs/wazero) (`internal/joltwasm`): instantiate once, step from Go, **copy** (or carefully alias) transform buffers from guest memory each frame.
3. **Benchmark** host readback with `go test ./internal/joltwasm -bench=. -benchmem` and optional `JOLT_WASM=/path/to/custom.wasm`.

## Physics state buffer (host contract)

Guest linear memory begins with a fixed **16-byte little-endian header** (`PhysicsStateHeader` in [`internal/joltwasm/state_view.go`](../../internal/joltwasm/state_view.go)): `version`, `bodyCount`, `stride` (bytes per body in the SoA block), `reserved`. Packed `float32` body data follow immediately. The host uses [`StateView`](../../internal/joltwasm/state_view.go) and `Memory.Read` to obtain a **slice view** (no per-frame heap allocation); compare `BenchmarkReadbackView` vs `BenchmarkReadbackCopy`.

Pin Wasm memory **min=max pages** (or avoid `memory.grow` during reads) so views stay valid; see wazero `api.Memory` documentation.

Expectations: call and copy overhead dominate for small scenes; profile before committing to this backend on all tier-1 platforms.
