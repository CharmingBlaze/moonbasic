# Window (`WINDOW.*`)

**Conventions:** [STYLE_GUIDE.md](../../STYLE_GUIDE.md), [API_CONVENTIONS.md](API_CONVENTIONS.md) — reference pages use uppercase **`NAMESPACE.ACTION`**; Easy Mode (`Window.Open`, …) is [compatibility only](../../STYLE_GUIDE.md#easy-mode-compatibility-layer).

**Page shape:** [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) — see [WAVE.md](WAVE.md) (registry-first headings, **Full Example** at the end).

## Core Workflow

1. **Open:** **`WINDOW.OPEN(width, height, title)`** at startup.
2. **Configure:** **`WINDOW.SETFPS`**, **`WINDOW.SETTITLE`**, …
3. **Loop:** **`WHILE NOT WINDOW.SHOULDCLOSE()`**
4. **Shutdown:** **`WINDOW.CLOSE()`**

## Full Example

```basic
WINDOW.OPEN(1280, 720, "My Game")
WINDOW.SETFPS(60)

WHILE NOT WINDOW.SHOULDCLOSE()
    RENDER.CLEAR(0, 0, 0)
    RENDER.FRAME()
WEND

WINDOW.CLOSE()
```

---

## Window management

### `WINDOW.OPEN(width, height, title)`
Opens the application window and initializes the OpenGL context.

- **Arguments**:
    - `width, height`: (Integer) Client area dimensions.
    - `title`: (String) Window title bar text.
- **Returns**: (None)

---

### `WINDOW.SHOULDCLOSE()`
Returns `TRUE` if the user has requested to close the window.

- **Returns**: (Boolean)

---

### `WINDOW.SETFPS(fps)`
Sets the target frame rate (e.g., 60).

- **Returns**: (None)

---

### `WINDOW.CLOSE()`
Closes the window and terminates the engine host.

---

### `WINDOW.SETTITLE(title)`
Updates the window title bar at runtime.

---

### `WINDOW.SETPOSITION(x, y)`
Sets the screen position of the window's top-left corner.

---

### `WINDOW.SETSIZE(w, h)`
Resizes the client area in pixels.

---

### `WINDOW.GETMONITORCOUNT()`
Returns the number of connected physical displays.

- **Returns**: (Integer)

---

## Easy Mode name map (compatibility only)

| Facade | Registry |
|--------|----------|
| `Window.Open` | `WINDOW.OPEN` |
| `Window.Close` | `WINDOW.CLOSE` |
| `Window.ShouldClose` | `WINDOW.SHOULDCLOSE` |
| `Window.SetFPS` | `WINDOW.SETFPS` |
| `Window.SetTitle` | `WINDOW.SETTITLE` |
| `Window.SetPosition` | `WINDOW.SETPOSITION` |
| `Window.SetIcon` | `WINDOW.SETICON` |
| `Window.SetOpacity` | `WINDOW.SETOPACITY` |
| `Window.SetSize` | `WINDOW.SETSIZE` |
| `Window.GetPositionX` / `GetPositionY` | `WINDOW.GETPOSITIONX` / `WINDOW.GETPOSITIONY` |
| `Window.DPIScale` | `WINDOW.DPISCALE` |
