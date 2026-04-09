# Zero-CGO Raylib strategy (purego + sidecar shared library)

## Goal

Ship **`raylib.dll` / `libraylib.so` / `libraylib.dylib`** next to the binary and load symbols with **[ebitengine/purego](https://github.com/ebitengine/purego)** so **`CGO_ENABLED=0`** builds do not require a C compiler.

## Current repo snapshot (maintenance signal)

moonBASIC references **`rl.`** across **~100+** files under [runtime/](../../runtime/) (window, draw, camera, audio, texture, mbmodel3d, mbentity, input, etc.). A full port is a **large fork** of [gen2brain/raylib-go](https://github.com/gen2brain/raylib-go); treat it as a **multi-release** effort.

**Regenerate usage counts** (from repo root):

```bash
rg -c "rl\." runtime/ --glob "*.go" | wc -l
rg "rl\.[A-Za-z]+" runtime/ --glob "*_cgo.go" -o | sort -u
```

The second command lists **unique** `rl.` identifiers to prioritize for `RegisterLibFunc` binding.

## Spike code

See [internal/raylibpurego/](../../internal/raylibpurego/): minimal **dynamic load** + **`GetFrameTime`** registration to validate:

- DLL/SO discovery next to executable
- purego registration path
- no `import "C"` in that package

Expand symbol-by-symbol mirroring existing [raylib-go](https://github.com/gen2brain/raylib-go) signatures until parity with moonBASIC’s call surface.

## Risks

- Struct layout must match **C ABI** (padding, alignment).
- GLFW / platform threads: Raylib still expects **main-thread** graphics; purego does not remove that rule.
- **raygui** is additional surface area (see [runtime/mbgui/](runtime/mbgui/)).

## Symbol inventory

Generated list of unique `rl.*` suffixes: [docs/audit/raylib_symbol_gap.txt](../audit/raylib_symbol_gap.txt).

## CI

[.github/workflows/ci.yml](../../.github/workflows/ci.yml) **`windows_purego`** already runs **`go build` with `CGO_ENABLED=0`** (raylib-go purego path on Windows). **`stub-only`** and **`linux_cgo_zero_modernc`** exercise **`CGO_ENABLED=0`** test matrix pieces. Extend coverage as bindings grow.
