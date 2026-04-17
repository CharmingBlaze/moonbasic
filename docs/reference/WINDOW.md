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
Opens the window (**client** width/height in pixels) and title bar text.

---

### `WINDOW.CLOSE()` 
Closes the window and tears down the host.

---

### `WINDOW.SHOULDCLOSE()` 
Returns **`TRUE`** when the user asked to close (title bar / Alt+F4 / etc.).

---

### `WINDOW.SETFPS(fps)` / `WINDOW.SETTARGETFPS(fps)` 
Target frame rate (**`SETTARGETFPS`** is the paired name in the manifest where both exist).

---

## Appearance and position

### `WINDOW.SETTITLE(title)` 
Runtime title change.

```basic
score = 100
WINDOW.SETTITLE("My Game | Score: " + STR(score))
```

---

### `WINDOW.SETPOSITION(x, y)` 
Screen position of the window’s top-left corner.

---

### `WINDOW.SETICON(filePath)` 
Loads a window icon (square **`.png`** recommended, e.g. 64×64).

---

### `WINDOW.SETOPACITY(alpha)` 
Window transparency.

---

### `WINDOW.SETSIZE(w, h)` 
Resizes the client area (pixels).

---

### `WINDOW.GETPOSITIONX()` / `WINDOW.GETPOSITIONY()` 
Current screen position of the top-left corner.

---

### `WINDOW.DPISCALE()` 
Global DPI scale factor for high-DPI displays.

---

## Monitors

### `WINDOW.GETMONITORCOUNT()` 
Number of connected monitors.

---

### `WINDOW.GETMONITORWIDTH(monitor)` / `WINDOW.GETMONITORHEIGHT(monitor)` 
Pixel size of the given monitor index (**0** = primary).

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
