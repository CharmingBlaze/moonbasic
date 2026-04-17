# Raylib Extras

Index of Raylib-mapped namespaces: window, input, render, draw, time, textures, and more.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

This page is a quick-map to dedicated reference pages. Use registry-first keys (`WINDOW.OPEN`, `INPUT.KEYDOWN`, `RENDER.CLEAR`, etc.). Requires **CGO** for the full Raylib stack.

---

## Namespaces (quick map)

| Area | Registry prefix | Reference |
|------|------------------|-----------|
| Window | `WINDOW.*` | This page, [PROGRAMMING.md](../PROGRAMMING.md) |
| Input | `INPUT.*`, `KEY_*` | [INPUT.md](INPUT.md) |
| Gestures | `GESTURE.*` | [INPUT.md](INPUT.md) |
| Render state | `RENDER.*` | [RENDER.md](RENDER.md) |
| 2D draw | `DRAW.*` | [DRAW2D.md](DRAW2D.md) |
| 3D draw | `DRAW3D.*` | [DRAW3D.md](DRAW3D.md) |
| Time | `TIME.*` | [TIME.md](TIME.md) if present, else built-in docs |
| GPU textures | `TEXTURE.*`, `RENDERTARGET.*` | [TEXTURE.md](TEXTURE.md) |
| Clipboard | `SYSTEM.GETCLIPBOARD`, `SYSTEM.SETCLIPBOARD` | [SYSTEM.md](SYSTEM.md) if present |

---

## `WINDOW.*` — open window, placement, state

Core lifecycle: **`WINDOW.OPEN`**, **`WINDOW.SETFPS`** / **`WINDOW.SETTARGETFPS`** (alias; both call `SetTargetFPS`), **`WINDOW.CLOSE`**, **`WINDOW.SHOULDCLOSE`**.

**Placement and size (desktop):**

| Command | Role |
|---------|------|
| **`WINDOW.SETPOS`** `(x, y)` (canonical) / deprecated **`WINDOW.SETPOSITION`** | `SetWindowPosition` |
| **`WINDOW.SETSIZE`** `(w, h)` | `SetWindowSize` |
| **`WINDOW.GETPOSITIONX`** / **`WINDOW.GETPOSITIONY`** | Current window position |

**Window chrome:**

| Command | Role |
|---------|------|
| **`WINDOW.MINIMIZE`** | `MinimizeWindow` |
| **`WINDOW.MAXIMIZE`** | `MaximizeWindow` (when resizable) |
| **`WINDOW.RESTORE`** | `RestoreWindow` |
| **`WINDOW.TOGGLEFULLSCREEN`** | `ToggleFullscreen` |
| **`WINDOW.SETTITLE`** | `SetWindowTitle` |

**Flags and monitors:** **`WINDOW.SETFLAG`**, **`WINDOW.CLEARFLAG`**, **`WINDOW.CHECKFLAG`**, min/max size, monitor queries, DPI, icon, opacity — see manifest and `runtime/window/window_state_cgo.go`.

---

## `RENDER.*` — clear, frame, GL-ish state

Handled mainly in **`runtime/window`** (clear/frame, blend, depth, scissor, wireframe, screenshot, MSAA, ambient, shadow map size, IBL, etc.) and **`runtime/mbmodel3d`** for PBR-related hooks. See **[RENDER.md](RENDER.md)** for the full loop and **[LIGHT.md](LIGHT.md)** for lighting-related **`RENDER.SET*`** calls.

---

## `INPUT.*` — keyboard, mouse, gamepad

See **[INPUT.md](INPUT.md)**. Key codes use **`KEY_*`** globals. **`INPUT.GETKEYNAME`** resolves names where supported.

---

## `DRAW.*` / `DRAW3D.*`

Primitives, text, textures, billboards — see **[DRAW2D.md](DRAW2D.md)** and **[DRAW3D.md](DRAW3D.md)**. **`TEXTURE.*`** and **[TEXTURE.md](TEXTURE.md)** cover GPU image handles used by **`Draw.Texture`**.

---

## Clipboard

**`SYSTEM.GETCLIPBOARD`** (returns string) and **`SYSTEM.SETCLIPBOARD`** `(text)` wrap Raylib clipboard when CGO is enabled; see **`runtime/system/clipboard_cgo.go`**.

---

## See also

- [PROGRAMMING.md](../PROGRAMMING.md) — main loop, shutdown
- [ARCHITECTURE.md](../../ARCHITECTURE.md) — threading and main thread
- [TEXTURE.md](TEXTURE.md) — render targets and textures
