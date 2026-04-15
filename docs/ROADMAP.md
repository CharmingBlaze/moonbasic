# moonBASIC roadmap (engineering phases)

This file is **maintainer documentation** for the long-term plan: API polish, rendering, 2D/3D systems, tools, and language features. It is not a promise of delivery dates; priorities shift with contributors and dependencies.

## Right now — Phase 1 (polish)

| Track | Goal |
|--------|------|
| **1A** | **API consistency** — one naming story for spatial handles (`SetPos` + optional `SetPosition` alias), documented in [API_CONSISTENCY.md](./API_CONSISTENCY.md) (regenerate with `go run ./tools/apidoc`). |
| **1B** | **Error quality** — compile-time did-you-mean + arity hints; runtime file/line wrapping; heap `Cast` hints. See [ERROR_MESSAGES.md](./ERROR_MESSAGES.md). |
| **1C** | **Sensible defaults** — `CAMERA.CREATE`, `LIGHT.CREATE`, `BODY3D.CREATE`, materials, etc. (verify deprecated `*.MAKE` aliases and each module’s `LOAD` paths). |
| **1D** | **Debug overlay** — `DEBUG.WATCH` / `DEBUG.WATCHCLEAR` with on-screen panel when **CGO + Raylib** (`runtime/mbdebug/overlay_cgo.go`), hooked from the window frame path in `compiler/pipeline/pipeline.go`. |

## Next — Phases 2–4 (features)

Roughly: **PBR / shadows / instancing / particles / post** (Phase 2), **tilemap / atlas / anim / 2D lighting / transitions** (Phase 3), **scene / save / input map / pool / tween / events** (Phase 4). Details belong in design docs per subsystem when work starts.

## Later — Phases 5–10

Deferred rendering, advanced effects, compute, decals, AI/navigation, multiplayer, editor tooling, packaging, web export, and language features (lambdas, interpolation, optional chaining, expanded types). Track in issues or ADRs as scope firms up.

## Ship criteria (v1.0 direction)

Examples from the product brief: spinning cube demo solid; one complete 2D sample; one 3D sample; multiplayer sample; single-file packaging; documentation coverage; multiple finished example games. Treat as **release goals**, not a checklist for every intermediate commit.
