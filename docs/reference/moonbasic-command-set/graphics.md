# Graphics (window + frame)

| Designed | moonBASIC | Notes |
|----------|------------|-------|
| **Graphics(w, h)** | **`Window.Open()`** | Opens the application window. |
| **Graphics3D(w, h)** | **`Window.Open()`** | Same as Graphics. |
| **AppTitle(title)** | **`Window.SetTitle()`** | Updates the window title. |
| **SetFPS(fps)** | **`Window.SetFPS()`** | Sets the target frame rate. |
| **Flip()** | **`Render.Frame()`** | Presents the rendered frame (swap buffers). |
| **Cls()** | **`Render.Clear()`** | Clears the screen buffers. |
| **SetClearColor(r, g, b)** | **`Render.SetClearColor()`** | |
| **SetVSync(on)** | **`Window.SetFlag()`** | Use `FLAG_VSYNC_HINT`. |
