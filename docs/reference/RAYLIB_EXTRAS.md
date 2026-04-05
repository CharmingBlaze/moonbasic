# Raylib extras — window, input, render, draw

moonBASIC exposes Raylib through focused namespaces. This page maps **where** common behaviors live; see each topic file for parameters.

| Namespace | Package | Topics |
|-----------|---------|--------|
| `WINDOW.*` | `runtime/window` | Open/close, FPS, VSYNC, flags |
| `INPUT.*` / `KEY_*` | `runtime/input` | Keyboard, mouse, gamepad |
| `GESTURE.*` | `runtime/input` | Touch gestures |
| `RENDER.*` | `runtime/window` + `runtime/draw` | Clear, frame, 2D/3D mode |
| `DRAW.*` / `DRAW3D.*` | `runtime/draw` | Primitives, text, billboards |
| `TIME.*` | `runtime/time` | Delta, wall clock |

**Requires CGO** for the full Raylib path; stub builds return errors with a **`CGO_ENABLED=1`** hint.

---

## See also

- [PROGRAMMING.md](../PROGRAMMING.md) — game loop and shutdown
- [DRAW2D.md](DRAW2D.md) — 2D drawing
- [INPUT.md](INPUT.md) — input reference
