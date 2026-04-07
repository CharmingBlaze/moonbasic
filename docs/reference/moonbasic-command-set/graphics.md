# Graphics (window + frame)

| Designed | Implementation | Memory / notes |
|----------|----------------|----------------|
| **Graphics (w, h, depth)** | **`WINDOW.OPEN`** — depth via **3D mode** / MSAA / Z-buffer, not a third `Graphics` arg | Returns nothing; window owns GL context. |
| **Graphics3D (w, h, depth)** | Same as above + **`CAMERA.MAKE`** / **`RENDER.*`** pipeline | See [WINDOW.md](../WINDOW.md), [RENDER.md](../RENDER.md). |
| **SetVSync (mode)** | **`WINDOW.SETFLAG`** / vsync-related **window** flags if exposed | Check manifest for exact key. |
| **SetClearColor (r, g, b)** | **`RENDER.CLEAR`** args or clear colour on **`RENDER.CLEAR`** | Typically **`r,g,b,a`** 0–255. |
| **Clear ()** | **`RENDER.CLEAR`** | No handle. |
| **Flip ()** | **`RENDER.FRAME`** | Presents the back buffer (Raylib frame). |
