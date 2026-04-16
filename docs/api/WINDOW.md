# Window Commands

Commands for creating, configuring, and managing the application window. The window is the foundation of every graphical moonBASIC program — it must be opened before any rendering, input, or audio commands will function.

## Core Workflow

Every graphical program follows the same lifecycle: **open** a window, run a **game loop** with `RENDER.CLEAR` / `RENDER.FRAME`, and **close** on exit. The window owns the OpenGL context and the frame timing, so nothing visual works until `WINDOW.OPEN` succeeds.

---

## Creation & Lifecycle

### `Window.Open(width, height, title)`

Opens a new window with the given pixel dimensions and title bar text. Internally initializes the Raylib OpenGL context, sets default FPS to 60, polls initial input events, and primes the GPU driver for drawing.

- `width` (int) — Window width in pixels.
- `height` (int) — Window height in pixels.
- `title` (string) — Text shown in the title bar.

**How it works:** Calls `InitWindow` through the HAL driver, sets target FPS to 60, resizes the framebuffer to match the requested client area (important for HiDPI), then runs a short warm-up sequence of blank frames to prevent visual glitches on some GPU drivers (Intel Iris, etc). If a window is already open, it is closed first.

**Errors:**
- `"WINDOW.OPEN expects 3 arguments (width, height, title)"`
- `"WINDOW.OPEN: width and height must be numeric"`
- `"WINDOW.OPEN: title must be a string"`

```basic
Window.Open(1280, 720, "My Game")
```

---

### `Window.CanOpen(width, height, title)`

Tests whether a window could be opened with these parameters, without actually opening it. Returns `TRUE` if width > 0, height > 0, and title is non-empty.

- `width` (int) — Desired width.
- `height` (int) — Desired height.
- `title` (string) — Desired title.

**Returns:** `bool`

```basic
IF Window.CanOpen(1920, 1080, "HD Game") THEN
    Window.Open(1920, 1080, "HD Game")
ELSE
    Window.Open(1280, 720, "HD Game")
ENDIF
```

---

### `Window.Close()`

Closes the window, destroys the OpenGL context, frees all heap handles, and shuts down audio if it was initialized. This should be the last call in your program.

**How it works:** Calls `CloseWindow` through the HAL driver, then calls `Heap.FreeAll()` to release every allocated handle (textures, models, sounds, etc). The window cannot be used after this call.

```basic
Window.Close()
```

---

### `Window.ShouldClose()`

Returns `TRUE` when the user has requested the window to close (clicked the X button, pressed Alt+F4, etc). Use this as your main loop condition.

**Returns:** `bool`

**How it works:** Delegates to Raylib's `WindowShouldClose()`. Returns `FALSE` if no window is open yet (so you won't spin forever before calling `WINDOW.OPEN`).

```basic
WHILE NOT Window.ShouldClose()
    Render.Clear(40, 40, 40)
    Draw.Text("Hello moonBASIC!", 10, 10, 20, 255, 255, 255, 255)
    Render.Frame()
WEND
```

---

## Frame Rate & Timing

### `Window.SetFPS(fps)`

Sets the target frames per second. The engine will sleep between frames to hit this cap. Set to 0 to uncap (limited only by VSync or GPU).

- `fps` (int) — Target frame rate (e.g., 60, 120, 144).

**How it works:** Calls `SetTargetFPS` through the HAL. The actual frame rate may be lower if your game logic or GPU can't keep up, or different if VSync is active.

**Errors:**
- `"WINDOW.SETFPS: window is not open (call WINDOW.OPEN first)"`

```basic
Window.Open(1280, 720, "Smooth Game")
Window.SetFPS(144)
```

---

### `Window.SetTargetFPS(fps)`

Alias for `Window.SetFPS`. Identical behavior.

---

### `Window.GetFPS()`

Returns the current actual frames per second as reported by the GPU driver.

**Returns:** `int`

```basic
fps = Window.GetFPS()
Draw.Text("FPS: " + STR(fps), 10, 10, 20, 255, 255, 255, 255)
```

---

## Window Properties

### `Window.Width()`

Returns the current window width in pixels.

**Returns:** `int`

---

### `Window.Height()`

Returns the current window height in pixels.

**Returns:** `int`

```basic
w = Window.Width()
h = Window.Height()
Draw.Text("Resolution: " + STR(w) + "x" + STR(h), 10, 40, 16, 200, 200, 200, 255)
```

---

### `Window.DPIScale()`

Returns the DPI scaling factor for High-DPI / Retina displays. On a standard display this returns `1.0`.

**Returns:** `float`

---

### `Window.SetTitle(title)`

Changes the window title bar text at runtime.

- `title` (string) — New title text.

```basic
Window.SetTitle("My Game - Score: " + STR(score))
```

---

### `Window.SetSize(width, height)`

Resizes the window to new pixel dimensions.

- `width` (int) — New width.
- `height` (int) — New height.

---

### `Window.SetPos(x, y)`

Moves the window to a screen position.

- `x` (int) — X position in screen pixels.
- `y` (int) — Y position in screen pixels.

---

### `Window.SetPosition(x, y)`

Alias for `Window.SetPos`. Deprecated — use `SetPos`.

---

### `Window.GetPositionX()`

Returns the window's current X position on screen.

**Returns:** `int`

---

### `Window.GetPositionY()`

Returns the window's current Y position on screen.

**Returns:** `int`

---

## Window State

### `Window.SetMSAA(samples)`

Sets the MSAA (Multi-Sample Anti-Aliasing) sample count. Must be called **before** `Window.Open` to take effect, as MSAA is a context-creation flag.

- `samples` (int) — Number of samples (typically 2, 4, or 8).

```basic
Window.SetMSAA(4)
Window.Open(1280, 720, "Anti-Aliased")
```

---

### `Window.IsFullscreen()`

Returns `TRUE` if the window is currently in fullscreen mode.

**Returns:** `bool`

---

### `Window.ToggleFullscreen()`

Toggles between fullscreen and windowed mode.

```basic
IF Input.KeyPressed(KEY_F11) THEN
    Window.ToggleFullscreen()
ENDIF
```

---

### `Window.IsResized()`

Returns `TRUE` if the window was resized since the last frame. Useful for recalculating UI layouts.

**Returns:** `bool`

---

### `Window.Minimize()`

Minimizes the window to the taskbar.

---

### `Window.Maximize()`

Maximizes the window to fill the screen.

---

### `Window.Restore()`

Restores the window from minimized or maximized state.

---

### `Window.SetFlag(flag)`

Sets a Raylib window configuration flag. Flags are integer constants (see Raylib documentation for flag values like `FLAG_VSYNC_HINT`, `FLAG_WINDOW_RESIZABLE`, etc).

- `flag` (int) — Raylib window flag constant.

---

### `Window.ClearFlag(flag)`

Removes a previously set window flag.

- `flag` (int) — Raylib window flag constant.

---

### `Window.CheckFlag(flag)`

Tests if a window flag is currently active.

- `flag` (int) — Raylib window flag constant.

**Returns:** `bool`

---

### `Window.SetOpacity(alpha)`

Sets the overall window opacity (transparency). Value from 0.0 (invisible) to 1.0 (fully opaque).

- `alpha` (float) — Opacity level.

---

### `Window.SetIcon(imagePath)`

Sets the window icon from an image file.

- `imagePath` (string) — Path to icon image file.

---

## Monitor Information

### `Window.GetMonitorCount()`

Returns the number of connected monitors.

**Returns:** `int`

---

### `Window.SetMonitor(index)`

Moves the window to a specific monitor.

- `index` (int) — Monitor index (0-based).

---

### `Window.GetMonitorName(index)`

Returns the name of a monitor.

- `index` (int) — Monitor index.

**Returns:** `string`

---

### `Window.GetMonitorWidth(index)`

Returns the width of a monitor in pixels.

- `index` (int) — Monitor index.

**Returns:** `int`

---

### `Window.GetMonitorHeight(index)`

Returns the height of a monitor in pixels.

- `index` (int) — Monitor index.

**Returns:** `int`

---

### `Window.GetMonitorRefreshRate(index)`

Returns the refresh rate of a monitor in Hz.

- `index` (int) — Monitor index.

**Returns:** `int`

---

### `Window.GetScaleDPIX()`

Returns the horizontal DPI scale factor.

**Returns:** `float`

---

### `Window.GetScaleDPIY()`

Returns the vertical DPI scale factor.

**Returns:** `float`

---

## Loading Mode

### `Window.SetLoadingMode(enabled)`

Enables or disables loading mode. When active, the engine can render loading screens during asset loading without blocking the main thread.

- `enabled` (bool) — `TRUE` to enable.

---

### `Window.LoadingMode()`

Returns whether loading mode is currently active.

**Returns:** `bool`

---

## Easy Mode Shortcuts

| Shortcut | Maps To |
|----------|---------|
| `SCREENWIDTH` | `Window.Width()` |
| `SCREENHEIGHT` | `Window.Height()` |
| `FPS` | `Window.GetFPS()` |
| `WindowWidth()` | `Window.Width()` |
| `WindowHeight()` | `Window.Height()` |
| `ScreenWidth()` | `Window.Width()` |
| `ScreenHeight()` | `Window.Height()` |
| `AppTitle(w, h, title)` | `Window.Open(w, h, title)` |
| `Graphics3D(w, h, depth)` | `Window.Open(w, h, "")` |

---

## Full Example

A complete program showing the window lifecycle, FPS control, fullscreen toggle, and basic drawing.

```basic
; Configure anti-aliasing before opening
Window.SetMSAA(4)

; Open a 720p window
Window.Open(1280, 720, "moonBASIC Window Demo")
Window.SetFPS(60)

score = 0
frames = 0

WHILE NOT Window.ShouldClose()
    frames = frames + 1

    ; Toggle fullscreen with F11
    IF Input.KeyPressed(KEY_F11) THEN
        Window.ToggleFullscreen()
    ENDIF

    ; Update title every 60 frames
    IF frames MOD 60 = 0 THEN
        Window.SetTitle("Demo - FPS: " + STR(Window.GetFPS()))
    ENDIF

    ; Render
    Render.Clear(25, 25, 40)

    Draw.Text("Window: " + STR(Window.Width()) + "x" + STR(Window.Height()), 10, 10, 20, 255, 255, 255, 255)
    Draw.Text("FPS: " + STR(Window.GetFPS()), 10, 35, 20, 100, 255, 100, 255)
    Draw.Text("DPI Scale: " + STR(Window.DPIScale()), 10, 60, 20, 100, 200, 255, 255)
    Draw.Text("Press F11 for fullscreen", 10, 90, 16, 180, 180, 180, 255)

    Render.Frame()
WEND

Window.Close()
```

---

## See Also

- [RENDER](RENDER.md) — Frame control, clear, and advanced rendering
- [INPUT](INPUT.md) — Keyboard, mouse, and gamepad input
- [DRAW](DRAW.md) — 2D and 3D drawing commands
