# Graphics (window + frame)

| Designed | Implementation | Memory / notes |
|----------|----------------|----------------|
| **Graphics (w, h, depth)** | **`WINDOW.OPEN`** — depth via **3D mode** / MSAA / Z-buffer, not a third `Graphics` arg | Returns nothing; window owns GL context. |
| **Graphics3D (w, h)** | **`Graphics3D`** (Blitz alias) — **resizes** the window after **`WINDOW.OPEN`**; optional 4-arg form sets reserved depth + HighDPI bit | See [WINDOW.md](../WINDOW.md), [RENDER.md](../RENDER.md), [MODERN_BLITZ_COMMANDS.md](../MODERN_BLITZ_COMMANDS.md). |

**Example (after `Window.Open`):**

```basic
Window.Open(1920, 1080, "Hi-Fi")
Window.SetFPS(60)
AppTitle("Project")
Graphics3D(1920, 1080)   ; any width/height; optional if Open already matched
SetMSAA(4)
```
| **SetVSync (mode)** | **`WINDOW.SETFLAG`** / vsync-related **window** flags if exposed | Check manifest for exact key. |
| **SetClearColor (r, g, b)** | **`RENDER.CLEAR`** args or clear colour on **`RENDER.CLEAR`** | Typically **`r,g,b,a`** 0–255. |
| **Clear ()** | **`RENDER.CLEAR`** | No handle. |
| **Flip ()** | **`RENDER.FRAME`** | Presents the back buffer (Raylib frame). |
