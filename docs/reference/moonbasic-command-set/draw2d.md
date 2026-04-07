# 2D drawing

No **global pen** — colour is **per call** or stored in your own variables.

| Designed | Implementation | Memory / notes |
|----------|----------------|----------------|
| **Plot (x, y)** | **`DRAW.PLOT`** / **`DRAW.PIXEL`** — **6 args** `(x,y,r,g,b,a)` | No heap. |
| **Line (x1…y2)** | **`DRAW.LINE`** — **8 args** with RGBA | No heap. |
| **Rect (x,y,w,h,filled)** | **`DRAW.RECTANGLE`** (fill) / **`DRAW.RECTLINES`** (outline) | No heap. |
| **Oval (x,y,w,h,filled)** | **`DRAW.OVAL`** / **`DRAW.ELLIPSE`** — centre + radii; convert box → centre & half-extents | No heap. |
| **Text (x, y, text$)** | **`DRAW.TEXT`**, **`DRAW.TEXTEX`** | String in pool; no extra handle for default font. |
| **SetColor / SetAlpha** | *(use locals)* | Pass into each **`DRAW.*`**. |
| **SetOrigin (x, y)** | **`CAMERA2D.SETOFFSET`**, **`SETTARGET`**, **`SETZOOM`**, **`SETROTATION`** | 2D camera handle may be a **heap object** — **`FREE`** when done. |
| **SetViewport (x,y,w,h)** | **`RENDER.SETSCISSOR`**, **`RENDERTARGET.*`** | Render targets are **handles** — **`RENDERTARGET.FREE`**. |

See [DRAW2D.md](../DRAW2D.md), [DRAW_WRAPPERS.md](../DRAW_WRAPPERS.md).
