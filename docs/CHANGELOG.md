# Changelog

This document tracks the recent development history of moonBASIC.

---

## Unreleased (May 2026)

### Language

- **`$"..."` string interpolation** — `{expr}` and `{expr:fmt}` desugar to **`STR`** / **`FORMAT`** + concatenation. Documented in [LANGUAGE.md](LANGUAGE.md) and [STRING.md](reference/STRING.md).
- **Multi-value `RETURN`** — `RETURN a, b, c` with **`x, y, z = fn()`** unpacking (1-based temporary array; no caller **`ERASE`**). Replaces the old “pack into **`DIM`** and **`ERASE`**” workaround as the preferred style in [LANGUAGE.md](LANGUAGE.md).
- **`ENUM` … `ENDENUM`** — grouped integer constants; members as **`State.IDLE`** or **`STATE_IDLE`**. See [LANGUAGE.md](LANGUAGE.md).
- **`FOR EACH var IN array … NEXT`** — array iteration.
- **`FOR var = EACH(Type) … NEXT`** — iterate live **`NEW(Type)`** instances (VM type registry).
- **`@FunctionName`** and anonymous **`FUNCTION() … ENDFUNCTION`** — first-class refs; **`cb(args)`** calls stored refs.
- **`WINDOW.SETLOOPMODE(mode, hz)`** — variable / fixed / semi-fixed **`TIME.DELTA()`** behavior.
- **`TIME.PHYSICSSTEPS()` / `TIME.PHYSICSSTEP()`** — fixed-step accumulator for multi-update physics loops.
- **`INPUT.ONGAMEPAD(pad, @callback)`** — gamepad connect/disconnect events **`(padIndex, connected)`**.
- **`TEXTURE.LOAD` / `TILEMAP.LOAD`** — paths resolved via **`ASSET.RESOLVE`**.
- **Coroutine auto-tick** — **`COROUTINE.START`** coroutines resume each frame without a manual **`RESUME`** loop.
- **`COROUTINE … ENDCOROUTINE`** block syntax — auto-starts a coroutine handle variable.
- **`FUNCTION f(x AS FLOAT) AS INT`** — optional typed signatures with static arity/return checks and **variable type inference** across assignments and calls.
- **`SPRITE.BUILTIN(name$)`** — game-jam placeholder sprites (`player`, `enemy`, `bullet`, …).
- **Profiler function view** — wall time per user **`FUNCTION`** in **`--profile`** output.
- **`SOUND.BUILTIN` / `FONT.BUILTIN`** — game-jam synthesized SFX and default font handle.
- **`moonbasic playground`** — local compile-check web UI with **Run** (headless VM + `PRINT` output); bytecode preview (`web/playground/`).
- **Default package registry** — `moonbasic install demo_extra` without `MOONBASIC_REGISTRY`.
- **`examples/gamejam/`** — zero-asset jam demo.
- **LSP function signatures** — hover + signature help for typed `FUNCTION` headers.
- **Hosted registry fallback** — remote index URL with offline bundled fallback.
- **Playground bytecode preview** — first lines of main chunk disassembly after compile.
- **`EVENT.ON/OFF/ONCE`** — callbacks accept string or `@func` ref.
- **Parser error recovery** — multiple syntax errors reported in one compile pass.
- **`ASSET.PATH` / `ASSET.RESOLVE`** — asset paths relative to the `.mb` file; **`MODEL.LOAD`** uses resolution.
- **`moonbasic pack game.mb`** — zip bundle with `.mbc` + `assets/` + bundled **`moonbasic`** executable (`-no-runtime` to omit).
- **Profiler wall time** — **`--profile`** / **`--profile-out`** include milliseconds per source line.
- **Command browser** — `web/command-browser.html` searchable API reference.
- **Coroutines (initial)** — **`YIELD`**, **`COROUTINE.START/RESUME/WAIT/DONE`**.
- **`IMPORT "package"`** — load packages from configured roots (`pkg/main.mb`, `pkg/index.mb`, or `pkg.mb`).

### Examples and docs

- **`examples/tilemap/`** — runnable Tiled map demo with collision layer.
- **`examples/gamepad/`** — controller axes and buttons with **`GAMEPAD_*`** constants.
- [ROADMAP.md](ROADMAP.md) — forward-looking feature tracker.
- [TILEMAP.md](reference/TILEMAP.md) — corrected **`TILEMAP.DRAW`** arity (handle only).

---

## Version 0.1 (April 2026)

### April 20, 2026 (release hygiene)

-   **Examples / Easy Mode**: Fixed **`examples/mario64/main_easymode.mb`** crashing on first frame when reading star positions. **`s.X()`** / **`s.Y()`** / **`s.Z()`** could parse as a **namespace call** (`S.X`, …) instead of a **handle method** if the parser had not yet registered **`s`** as a variable in that parse context. **Fix:** use **`starEnts(i).X()`** (and **`.Y()` / `.Z()`**) for position reads so the receiver is an **indexed expression**, which always becomes **`HandleCallExpr`** → correct **`ENTITY.*`** dispatch with the entity id. **`s.Hide()`** remains valid after **`s = starEnts(i)`** (statements cannot use **`arr(i).Hide()`** yet).
-   **CI**: Semantic check now includes **`examples/mario64/main_easymode.mb`** alongside **`main_entities.mb`**.
-   **Build**: **`scratch/clean_manifest.go`** and **`scratch/check_tags.go`** marked **`//go:build ignore`** so they do not share **`package main`** with **`verify_array_pt5_test.go`** (fixes duplicate **`main`** when running **`go test ./...`**).
-   **Audit**: Regenerated **`docs/audit/manifest_keys.txt`**, **`docs/audit/runtime_keys.txt`**, and **`docs/MISSING_COMMANDS_AUDIT.md`** via **`python tools/diff_manifest_runtime.py --write`** so **`--check`** matches **`commands.json`** and the runtime scan.

### April 5–6, 2026

-   **Blitz3D-style API**: **`CAMERA.TURN`**, **`ROTATE`**, **`ORBIT`** (alias of **`SETORBIT`**), **`ZOOM`**, **`FOLLOW`**, **`CAMERA.FOLLOWENTITY`**; **`ENTITY.CREATE`**, **`CREATEBOX`**, movement, simple collision/physics, **`DRAWALL`**; input aliases **`KEYHIT`**, **`MOUSEXSPEED`/`MOUSEYSPEED`**, **`JOYX`/`JOYY`/`JOYBUTTON`** (see [BLITZ3D.md](reference/BLITZ3D.md)).
-   **Gameplay / input helpers**: **`LANDBOXES`**, **`PLAYER.MOVERELATIVE`**, **`Input.Orbit`** (alias of **`Input.AxisDeg`**), **`Input.Movement2D`**, plus earlier **`MOVESTEPX`/`MOVESTEPZ`** and **`Input.AxisDeg`**. See [GAMEHELPERS.md](reference/GAMEHELPERS.md), [INPUT.md](reference/INPUT.md), [MATH.md](reference/MATH.md).
-   **Language**: Record types — **`TYPE` … `ENDTYPE`**, **`DIM name AS TypeName(n)`**, **`TypeName(...)`** field initialisers, **`arr(i).field`** access, **`ERASE`** for typed arrays. Documented in [LANGUAGE.md](LANGUAGE.md) and [ARRAY.md](reference/ARRAY.md).
-   **Input**: **`Input.Axis(negKey, posKey)`** returns `{-1, 0, 1}` for two-key axes — [INPUT.md](reference/INPUT.md).
-   **Math / gameplay**: **`MOVEX`** / **`MOVEZ`** (camera-relative XZ from yaw), **`IIF`** (string **`IIF`**) — [MATH.md](reference/MATH.md).
-   **Collision**: **`BOXTOPLAND`** returns a **float** (snap Y or `0.0`), not a boolean — [GAMEHELPERS.md](reference/GAMEHELPERS.md).
-   **Collision / picking**: **`RAY2D.*`** (circle, axis-aligned rect, segment) — pure math, always available; 3D **`RAY.*`** unchanged (Raylib; CGO). Documented in [RAYCAST.md](reference/RAYCAST.md).
-   **Docs**: Regenerate [API_CONSISTENCY.md](API_CONSISTENCY.md) with **`go run ./tools/apidoc`** when builtins change.

### April 4, 2026

-   **Math**: `CLAMP`, `LERP`, and `WRAP` commands now use the formulas from Raylib 5.5 for better consistency and performance. Trigonometric functions still use the standard Go math library.
-   **File I/O**: Clarified the behavior of file writing commands. `FILE.WRITE` and its alias `WRITEFILE` write raw data, while `FILE.WRITELN` and `WRITEFILELN` append a newline character.
-   **Control Flow**: Implemented a full suite of `DO...LOOP` structures (`DO WHILE`, `DO UNTIL`, `DO...LOOP WHILE`, `DO...LOOP UNTIL`). Added `EXIT` and `CONTINUE` statements for all loop types (`FOR`, `WHILE`, `REPEAT`, `DO`) and `EXIT FUNCTION` for early returns from functions.
-   **Parser**: Fixed a bug where `NEXT` in a `FOR` loop could incorrectly consume the first part of the next statement if it was on a new line. The optional variable after `NEXT` is now only considered if it's on the same line.
-   **Codebase**: Refactored loop-related parsing logic into a dedicated `parser_stmts_loop.go` file for better organization.
-   **Internal**: Removed obsolete `strmod.*` command registrations. The core `runtime` now handles all string-related built-in commands.
