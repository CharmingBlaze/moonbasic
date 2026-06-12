# moonBASIC — What to Add and Change to Make It Truly Great

Forward-looking companion to the release audit. Items are grouped by area, not severity — these are investments rather than bugs.

**Status legend:** ✅ Done · 🚧 Partial · ❌ Not started

---

## Language Features

| # | Feature | Status | Notes |
|---|---------|--------|-------|
| 1 | **First-class function references** | ✅ | **`@FuncName`**, anonymous **`FUNCTION() … ENDFUNCTION`**, **`cb(args)`** via `OpCallRef`; collision/tween/btree/event callbacks accept string or ref. |
| 2 | **String interpolation `$"…"`** | ✅ | Lexer + parser desugar to `STR`/`FORMAT` + concat. [LANGUAGE.md](LANGUAGE.md) |
| 3 | **Multi-return syntax** | ✅ | `RETURN a, b, c` + `x, y, z = fn()`. [LANGUAGE.md](LANGUAGE.md) |
| 4 | **Typed function signatures** | ✅ | `FUNCTION f(x AS FLOAT) AS INT` — parse, static checks with **variable inference**, LSP hover + signature help. |
| 5 | **`FOR EACH`** | ✅ | **`FOR EACH v IN arr`** and **`FOR e = EACH(Type)`** (VM type-instance registry). [LANGUAGE.md](LANGUAGE.md) |
| 6 | **`ENUM` / named integer sets** | ✅ | `ENUM State … ENDENUM`, `State.IDLE`. [LANGUAGE.md](LANGUAGE.md) |

---

## Developer Experience

| # | Feature | Status | Notes |
|---|---------|--------|-------|
| 7 | **Built-in game loop timing modes** | ✅ | **`WINDOW.SETLOOPMODE`**, **`TIME.PHYSICSSTEPS`**, **`TIME.PHYSICSSTEP`** for fixed-step accumulation. |
| 8 | **Package / module system** | ✅ | **`IMPORT "pkg"`**, [PACKAGES.md](PACKAGES.md), **`moonbasic install/list/publish`**, bundled default registry (`demo_extra`). |
| 9 | **Asset pipeline command** | ✅ | **`ASSET.PATH`**, **`ASSET.RESOLVE`**, loaders (**`MODEL`**, **`TEXTURE`**, **`TILEMAP`**); **`moonbasic pack`** bundles `.mbc`, assets, and runtime binary. |
| 10 | **Compiler error recovery** | ✅ | Parser collects multiple errors per pass; continues after sync-on-boundary. |
| 11 | **Stdlib game-math helpers** | ✅ | `MATH.LERPANGLE`, `MATH.PINGPONG`, `ARRAY.SORT`, `STRING.SPLIT`, etc. already in manifest. 2D overlap: use **`COLLISION.BOXOVERLAP2D`**. |

---

## Architecture / Runtime

| # | Feature | Status | Notes |
|---|---------|--------|-------|
| 12 | **Coroutines / fibers** | ✅ | **`YIELD`**, **`COROUTINE … ENDCOROUTINE`**, **`COROUTINE.*`**, auto-resume each frame. |
| 13 | **Tile-map end-to-end** | ✅ | `examples/tilemap/` — load TMX, draw, solid-layer collision. [TILEMAP.md](reference/TILEMAP.md) |
| 14 | **Gamepad support** | ✅ | `examples/gamepad/`; **`GAMEPAD_*`** constants; **`INPUT.GAMEPADCONNECTED`**; **`INPUT.ONGAMEPAD`** connect/disconnect events. |
| 15 | **Readable profiler output** | ✅ | **`--profile`** / **`--profile-out`** (lines + functions); function wall-time HTML via **`WriteProfileFlameHTML`**. |

---

## Ecosystem

| # | Feature | Status | Notes |
|---|---------|--------|-------|
| 16 | **Package registry** | ✅ | Default URL + bundled fallback; **`moonbasic list --remote`**. |
| 17 | **Interactive playground** | ✅ | **`moonbasic playground`** — compile, bytecode preview, **Run** (headless VM + PRINT); examples dropdown. |
| 18 | **Game jam mode** | ✅ | **`SPRITE.BUILTIN`**, **`SOUND.BUILTIN`**, **`FONT.BUILTIN`**; [examples/gamejam/](../examples/gamejam/). |
| 19 | **Visual command browser** | ✅ | Static searchable UI: [web/command-browser.html](../web/command-browser.html) (regenerate: `go run ./tools/cmdbrowser`). |
| 20 | **BlitzBASIC 3D porting guide** | ✅ | [BLITZ3D_PORTING.md](BLITZ3D_PORTING.md) — side-by-side Blitz vs moonBASIC; compat layer in `blitzengine`. |

---

## Priority order (solo dev)

Engineering order that maximizes user-visible return:

1. String interpolation — ✅  
2. Multi-return — ✅  
3. FOR EACH (array + type) — ✅  
4. ENUM — ✅  
5. Stdlib one-liners — ✅  
6. Tilemap example — ✅  
7. Gamepad example + constants — ✅  
8. Function references — ✅  
9. Coroutines — ✅  
10. Package system — ✅  

Items **16–20** can proceed in parallel with language work.

---

## Deliberately not planned

- **Classes and inheritance** — `TYPE` + handles + entity system is the model.  
- **Exceptions** — use error returns, `ON ERROR`, or clean runtime messages.  
- **Generics / templates** — optional typed signatures (#4) cover static checking without syntax overhead.  
- **Embedded script VM** — moonBASIC *is* the script language.  
- **Whitespace-sensitive blocks** — keywords (`THEN`, `WEND`, `ENDIF`) delimit blocks.

---

## Recently shipped (this cycle)

See [CHANGELOG.md](CHANGELOG.md) for detail.

- `$"…"` interpolation, multi-return, ENUM, FOR EACH (array + type)  
- `@Func` references + callback APIs (partial)  
- `IMPORT "package"` (basic package roots)  
- `examples/tilemap/`, `examples/gamepad/`  
- Docs: [LANGUAGE.md](LANGUAGE.md), reference updates, this roadmap  
