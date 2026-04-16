# MoonBASIC API Standardization Directive

**TO:** Compiler Engineering Team  
**FROM:** Language Design Authority  
**RE:** Complete API Consistency Overhaul  
**PRIORITY:** HIGH — Foundation for v1.0

**Related:** [STYLE_GUIDE.md](../STYLE_GUIDE.md) (normative style), [reference/API_CONVENTIONS.md](reference/API_CONVENTIONS.md), [EASY_MODE.md](EASY_MODE.md), [AGENTS.md](../AGENTS.md).  
**Note:** Examples below use **plain identifiers** only — moonBASIC does **not** use Blitz-style `#` / `$` / `?` / `%` name suffixes in source (registry may still list legacy `$` keys as deprecated aliases).

---

## Executive decision: API philosophy

moonBASIC is **standardized around the `Namespace.Method` pattern** with these principles:

### Approved design decisions

1. **Primary API style:** `Namespace.Method()` (e.g. `CAMERA.CREATE()`, `MODEL.SETPOS()`).
2. **Creation verb:** `CREATE` is standard (`MAKE` remains as deprecated manifest alias; warnings when tooling supports it).
3. **Position method:** `SETPOS` is canonical (`SETPOSITION` is a deprecated alias where registered).
4. **Easy Mode:** Convenience layer only — not primary documentation ([EASY_MODE.md](EASY_MODE.md)).
5. **Universal methods:** Spatial handles implement `.pos()`, `.rot()`, `.scale()` (and color/alpha where applicable) per [reference/UNIVERSAL_HANDLE_METHODS.md](reference/UNIVERSAL_HANDLE_METHODS.md).
6. **Method design:** Prefer **chainable, simple methods** over long positional argument lists.
7. **Symmetry (target):** Prefer **GET\*** for every **SET\*** where the surface is stable; every **CREATE** has **FREE** where applicable.

---

## Part 1 — Global naming standards

### Creation pattern (required)

```basic
' Correct — CREATE for instantiation
camera = CAMERA.CREATE()
model = MODEL.LOAD("player.glb")
light = LIGHT.CREATEPOINT()
body = BODY3D.CREATE()

' Deprecated — MAKE aliases (manifest + runtime)
camera = CAMERA.MAKE()
```

**Actions:** Keep `*.MAKE` as deprecated manifest aliases (see `compiler/builtinmanifest/api_standardization_test.go`); prefer `*.CREATE` in new code and examples; [MIGRATION_CREATE_FROM_MAKE.md](MIGRATION_CREATE_FROM_MAKE.md); run `go run ./tools/apidoc` when manifest changes.

### Position / transform pattern (required)

```basic
CAMERA.SETPOS(cam, x, y, z)
MODEL.SETROT(model, pitch, yaw, roll)
SPRITE.SETSCALE(sprite, sx, sy, sz)
```

`SETPOSITION` exists as a deprecated alias mapping to `SETPOS` where registered.

---

## Part 2 — Universal handle methods

Mandatory behaviors are implemented in the VM via `handleCallDispatch` / `handleCallBuiltin` and documented in [reference/UNIVERSAL_HANDLE_METHODS.md](reference/UNIVERSAL_HANDLE_METHODS.md).

**Minimum expectations**

- **Position:** `handle.pos(x, y, z)` / `handle.pos()` (getter → `GETPOS` or equivalent).
- **Rotation:** `handle.rot(...)` where the type supports rotation (2D vs 3D semantics differ).
- **Scale:** `handle.scale(...)` where applicable.
- **Color / alpha:** `.col()`, `.alpha()` for renderables where implemented.
- **Cleanup:** `handle.free()` → `NAMESPACE.FREE`.

### Coverage matrix (target — audit ground truth)

Implementation varies by type. Use this as a **review checklist**, not a guarantee every cell is implemented today:

| Handle area | .pos | .rot | .scale | .col / .alpha | .free | Notes |
|-------------|------|------|--------|---------------|-------|--------|
| CAMERA / CAMERA2D | ✓ | ✓* | — | — | ✓ | 2D may use single-axis rotation |
| MODEL / SPRITE | ✓ | ✓* | ✓ | ✓ | ✓ | |
| LIGHT | ✓ | ✓ | — | ✓ | ✓ | |
| BODY3D / BODY2D | ✓ | ✓* | ✓ | — | ✓ | Physics-specific helpers separate |
| PARTICLE / DECAL / TERRAIN / WATER / NAVAGENT | ✓ | ✓ | varies | varies | ✓ | |

\*2D-style types may use a single rotation value where documented.

**Actions:** Audit `vm/handlecall.go` and the manifest; extend dispatch before claiming parity with Blitz-style megafunctions.

---

## Part 3 — Ergonomics: reduce argument count

Prefer **CREATE + chainable setters** over long positional calls:

```basic
cam = CAMERA.CREATE()
cam.pos(0, 10, 20)
CAMERA.SETTARGET(cam, 0, 0, 0)
CAMERA.SETFOV(cam, 45)
```

Setters used in method-chaining should return the receiver handle where the VM implements it ([vm/vm_dispatch.go](../vm/vm_dispatch.go)).

---

## Part 4 — Namespace organization

Standard verbs:

| Verb | Role |
|------|------|
| `CREATE` / `CREATE*` | Construct |
| `LOAD` | File-backed assets |
| `FREE` | Release handle |
| `SETPOS`, `SETROT`, `SETSCALE`, … | Canonical setters |
| `GETPOS`, `GETROT`, `GETSCALE`, … | Canonical getters |

See [API_CONSISTENCY.md](API_CONSISTENCY.md) (generated from `compiler/builtinmanifest/commands.json`). Per-namespace sketches (CAMERA, MODEL, LIGHT, BODY3D, …) belong in **reference docs** and [STYLE_GUIDE.md](../STYLE_GUIDE.md); avoid duplicating full API lists here.

---

## Part 5 — Easy Mode

Easy Mode maps global helpers to `Namespace.Method` commands — **syntax sugar only**. Mappings and policy: [EASY_MODE.md](EASY_MODE.md). New examples and libraries should prefer **`CAMERA.*`**, **`ENTITY.*`**, etc.

---

## Part 6 — Documentation standards

| Document | Role |
|----------|------|
| [STYLE_GUIDE.md](../STYLE_GUIDE.md) | Normative style (CREATE, SETPOS, chaining) |
| [reference/API_CONVENTIONS.md](reference/API_CONVENTIONS.md) | Verb table and naming |
| [API_CONSISTENCY.md](API_CONSISTENCY.md) | Machine-checked command list |
| Per-namespace `docs/reference/*.md` | Type-specific behavior |

Namespace reference pages should lead with **`NAMESPACE.CREATE` / `LOAD` / `FREE`**, then **GET\*/SET\***, then handle methods.

---

## Part 7 — Implementation checklist

### Phase 1 — Manifest / compiler

- [x] `*.MAKE` rows remain **deprecated aliases** beside `*.CREATE` (`compiler/builtinmanifest/api_standardization_test.go`)
- [x] `SETPOSITION` aliases beside `SETPOS` where applicable (`TestSetPositionDeprecationMatchesSetPosKeys`)
- [x] Core string / file / GUI helpers: canonical non-`$` names + `deprecated_of` where applied; `go run ./tools/apidoc` after manifest edits
- [x] Duplicate **AUDIO** manifest rows (arity) — guarded by `TestNoDuplicateAudioManifestArgSignatures` in [`api_standardization_test.go`](../compiler/builtinmanifest/api_standardization_test.go); fix manifest if the test fails
- [x] Representative `go run . --check` on migrated `testdata/*.mb` (ongoing when touching builtins)

### Phase 2 — Runtime / VM

- [ ] Close gaps in universal handle paths (`vm/handlecall.go`) vs matrix above
- [ ] Setter chaining returns handle where specified (`value.KindNil` → receiver)
- [x] Deprecation for `MAKE` / `SETPOSITION` aliases: semantic notices + stderr warnings on compile/check ([`compiler/pipeline/compile.go`](../compiler/pipeline/compile.go)); `--strict-deprecated` on `moonbasic` / `go run .`; LSP diagnostic `deprecated-api` ([`lsp/server.go`](../lsp/server.go))

### Phase 3 — Documentation

- [x] [STYLE_GUIDE.md](../STYLE_GUIDE.md)
- [x] This directive + cross-links (STYLE_GUIDE, API_CONVENTIONS, EASY_MODE, AGENTS, UNIVERSAL_HANDLE_METHODS, DIRECTIVE_BASELINE_REPORT, CONTRIBUTING); stub [STYLE_GUIDE.md](STYLE_GUIDE.md) → root [STYLE_GUIDE.md](../STYLE_GUIDE.md); [GETTING_STARTED.md](GETTING_STARTED.md) intro + first-window / game-loop / key-constants + Modern Blitz window/input lines use **`NAMESPACE.ACTION`**
- [~] Namespace reference pages — incremental: [reference/PHYSICS3D.md](reference/PHYSICS3D.md) + [moonbasic-command-set/physics-3d.md](reference/moonbasic-command-set/physics-3d.md) (**`PHYSICS3D.UPDATE`** vs **`STEP`**); **Falling Cube** / **CHARCONTROLLER**; [reference/TRANSFORM.md](reference/TRANSFORM.md) **spinning cube**; [reference/MODEL.md](reference/MODEL.md) (**`RENDER.Begin3D`**); [reference/SPRITE.md](reference/SPRITE.md), [reference/PHYSICS2D.md](reference/PHYSICS2D.md) (**`PHYSICS2D.START`** / **`STEP`** workflow), [reference/AUDIO.md](reference/AUDIO.md); [reference/IMAGE.md](reference/IMAGE.md), [reference/WINDOW.md](reference/WINDOW.md) (**`WINDOW.*`**, Easy Mode map), [reference/TIME.md](reference/TIME.md) (**`TIME.DELTA`**), [reference/FONT.md](reference/FONT.md) (**`DRAW.TEXTEX` / `DRAW.TEXTFONT`** in manifest), [reference/MATH.md](reference/MATH.md), [reference/RENDER.md](reference/RENDER.md), [reference/TILEMAP.md](reference/TILEMAP.md), [reference/NETWORK.md](reference/NETWORK.md) (client + server: numeric **`EVENT.TYPE`** **1/2/3**), [reference/TEXTURE.md](reference/TEXTURE.md), [reference/PARTICLES.md](reference/PARTICLES.md), [reference/DRAW2D.md](reference/DRAW2D.md) (**`RENDER.DRAWFPS`**), [reference/DRAW3D.md](reference/DRAW3D.md) (**`RENDER.BEGIN3D`**), [reference/TERRAIN.md](reference/TERRAIN.md) (**`TERRAIN.*`**), [reference/CAMERA.md](reference/CAMERA.md) (**`CULL.*`**, **`RENDER.BEGIN3D`**), [reference/BLITZ3D.md](reference/BLITZ3D.md) (Raylib pipeline table), [reference/WATER.md](reference/WATER.md) (**`WATER.*`**), [reference/VEHICLE.md](reference/VEHICLE.md), [reference/INPUT.md](reference/INPUT.md) (**`INPUT.*`**, **`CURSOR.*`**), [reference/SCATTER.md](reference/SCATTER.md) (**`SCATTER.*`** / **`PROP.*`** runtime), [reference/SCENE.md](reference/SCENE.md) (**`SCENE.*`** with parentheses), [reference/PLAYER.md](reference/PLAYER.md), [reference/CHARACTER.md](reference/CHARACTER.md); [reference/API_CONVENTIONS.md](reference/API_CONVENTIONS.md); [reference/ENTITY.md](reference/ENTITY.md) **Creating and Moving**; [reference/GAMEHELPERS.md](reference/GAMEHELPERS.md) / [reference/GAMEPLAY_HELPERS.md](reference/GAMEPLAY_HELPERS.md) (**`TIME.*`**, **`INPUT.*`**, **`CAMERA.*`**); [reference/DEBUG.md](reference/DEBUG.md) (**`TEXTURE.LOAD`** in snippets)
- [~] [moonbasic-command-set/](reference/moonbasic-command-set/README.md) bridge tables — [input.md](reference/moonbasic-command-set/input.md), [graphics.md](reference/moonbasic-command-set/graphics.md), [runtime.md](reference/moonbasic-command-set/runtime.md) aligned to **`INPUT.*`**, **`WINDOW.*`**, **`RENDER.*`**, **`TIME.DELTA`**
- [x] [MIGRATION_CREATE_FROM_MAKE.md](MIGRATION_CREATE_FROM_MAKE.md)

### Phase 4 — Tooling

- [x] LSP namespace completions: sort **CREATE** before **MAKE**, **SETPOS** before **SETPOSITION** (`completionMethodPriority`); omit manifest rows with `deprecated_of` from completion list (`IsDeprecatedAlias`)
- [~] VS Code snippets in [`.vscode/moonbasic-snippets.code-snippets`](../.vscode/moonbasic-snippets.code-snippets): **`mbcam`**, **`mbcreate`**, **`mbsetpos`**, **`mbframe3d`** (clear / Begin3D / draw / End3D / Frame), **`mbphys3d`** (`PHYSICS3D.UPDATE`)
- [ ] Optional: linter rules for style (future)

### Phase 5 — Examples

- [x] Namespace-first position reads: `BODY3D.GETPOS`, `CHARACTERREF.GETPOSITION`, `ENTITY.GETPOS`, `NAVAGENT.GETPOS` (replacing handle `.X()/.Y()/.Z()`) in `demo_vehicle`, `demo_jet`, `demo_boat`, `demo_platformer`, `kcc_*.mb`, `sphere_drop/main.mb`, `mario64/main_easymode`; `sphere_drop` uses `name AS HANDLE(dim)` and locals `be`/`bh` before `CLEARPHYSBUFFER`/`FREE` on array slots
- [~] Showcase / template: **`examples/high_fidelity/modern_template.mb`** uses **`WINDOW.SETSIZE`**, **`WINDOW.SETMSAA`** (manifest entry added), **`WINDOW.OPEN`** / **`INPUT` / `RENDER`**; still **Easy Mode** where no namespace twin exists: **`EntityPBR`**, **`CreatePointLight`** (parents to entity vs raw **`LIGHT.CREATEPOINT`**), **`CameraSmoothFollow`** (no **`CAMERA.*`** smooth-follow in manifest). **`PHYSICS3D.UPDATE`** = **`STEP`**; **`modern_kcc_hop.mb`** uses **`PHYSICS3D.UPDATE()`**
- [x] **Terrain / entity demos — namespace `ENTITY.*`:** `terrain_colored`, `terrain_chase`, `terrain_async`, `terrain_stress` (`ENTITY.SETPOS`, `SCALE`, `COLOR`, `GETPOS` / `GETXZ`, `DRAWALL`, `FREE` / `FREEENTITIES`); `mario64/main_entities.mb` (`COPY`, `TYPE`, `RADIUS`, `MOVECAMERARELATIVE`, `GROUNDED`, `ADDFORCE`, …); `spin_cube` (`ENTITY.DRAWALL`); `high_fidelity/modern_template.mb` (`ENTITY.DRAWALL`, `ENTITY.ADDFORCE`); `mario64/main_easymode` (`ENTITY.COLOR` for collectibles)

---

## Part 8 — Breaking change timeline

| Milestone | Behavior |
|-----------|----------|
| **Current (0.9.x)** | Deprecated aliases work; tooling surfaces canonical names first; optional strict deprecation later |
| **Minor** | Louder warnings for `MAKE` / `SETPOSITION` / legacy patterns; examples migrated |
| **v1.0 (goal)** | Remove deprecated surface only after a documented migration window |

---

## Part 9 — Success criteria

1. **Naming:** New docs and examples use `CREATE` and `SETPOS`.
2. **Universal methods:** Dispatch + manifest aligned for spatial types (matrix audited).
3. **Documentation:** STYLE_GUIDE + API_CONVENTIONS + generated API_CONSISTENCY stay consistent.
4. **Easy Mode:** Clearly secondary ([EASY_MODE.md](EASY_MODE.md)).
5. **Ergonomics:** Chaining where the runtime supports it; avoid unmaintainable 8+ positional “god” calls in new API.
6. **Tooling:** LSP steers authors toward canonical names.

---

## Appendix — Quick reference

```basic
obj = NAMESPACE.CREATE()
obj2 = NAMESPACE.LOAD("path/to/asset")

obj.pos(x, y, z)
p = obj.pos()

obj.rot(pitch, yaw, roll)
r = obj.rot()

obj.scale(sx, sy, sz)
s = obj.scale()

obj.free()

NAMESPACE.SETPOS(obj, x, y, z)
NAMESPACE.FREE(obj)
```

---

**End of directive** — phases are iterative; backward compatibility is preserved while deprecated aliases remain in the manifest.
