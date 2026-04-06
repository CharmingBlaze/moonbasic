# Changelog

This document tracks the recent development history of moonBASIC.

---

## Version 0.1 (April 2026)

### April 5–6, 2026

-   **Blitz3D-style API**: **`CAMERA.TURN`**, **`ROTATE`**, **`ORBIT`** (alias of **`SETORBIT`**), **`ZOOM`**, **`FOLLOW`**, **`CAMERA.FOLLOWENTITY`**; **`ENTITY.CREATE`**, **`CREATEBOX`**, movement, simple collision/physics, **`DRAWALL`**; input aliases **`KEYHIT`**, **`MOUSEXSPEED`/`MOUSEYSPEED`**, **`JOYX`/`JOYY`/`JOYBUTTON`** (see [BLITZ3D.md](reference/BLITZ3D.md)).
-   **Gameplay / input helpers**: **`LANDBOXES`**, **`PLAYER.MOVERELATIVE`**, **`Input.Orbit`** (alias of **`Input.AxisDeg`**), **`Input.Movement2D`**, plus earlier **`MOVESTEPX`/`MOVESTEPZ`** and **`Input.AxisDeg`**. See [GAMEHELPERS.md](reference/GAMEHELPERS.md), [INPUT.md](reference/INPUT.md), [MATH.md](reference/MATH.md).
-   **Language**: Record types — **`TYPE` … `ENDTYPE`**, **`DIM name AS TypeName(n)`**, **`TypeName(...)`** field initialisers, **`arr(i).field`** access, **`ERASE`** for typed arrays. Documented in [LANGUAGE.md](LANGUAGE.md) and [ARRAY.md](reference/ARRAY.md).
-   **Input**: **`Input.Axis(negKey, posKey)`** returns `{-1, 0, 1}` for two-key axes — [INPUT.md](reference/INPUT.md).
-   **Math / gameplay**: **`MOVEX`** / **`MOVEZ`** (camera-relative XZ from yaw), **`IIF$`** (string **`IIF`**) — [MATH.md](reference/MATH.md).
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
