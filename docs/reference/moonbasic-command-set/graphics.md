# Graphics (window + frame)

**Conventions:** [STYLE_GUIDE.md](../../../STYLE_GUIDE.md), [API_CONVENTIONS.md](../API_CONVENTIONS.md). Use **`WINDOW.*`** / **`RENDER.*`** registry keys in new examples.

| Designed | Registry (use in new examples) | Notes |
|----------|-------------------------------|-------|
| **Graphics(w, h)** | **`WINDOW.OPEN(w, h, title)`** | Third arg is title string. |
| **Graphics3D(w, h)** | **`WINDOW.OPEN`** | Same window path; 3D is camera/render state. |
| **AppTitle(title)** | **`WINDOW.SETTITLE(title)`** | |
| **SetFPS(fps)** | **`WINDOW.SETFPS(fps)`** | |
| **Flip()** | **`RENDER.FRAME()`** | Presents the frame (swap buffers). |
| **Cls()** | **`RENDER.CLEAR(r, g, b)`** | Background clear each frame. |
| **SetClearColor(r, g, b)** | Use **`RENDER.CLEAR`** each frame, or Raylib **`RAYLIB.CLEARBACKGROUND`** if exposed | Per manifest / host. |
| **SetVSync(on)** | **`WINDOW.SETFLAG`** with vsync hint | See **`WINDOW.*`** / Raylib flags. |
