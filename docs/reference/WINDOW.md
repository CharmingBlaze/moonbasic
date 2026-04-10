# Window Commands

Commands for managing the application window.

## Core Workflow

A typical moonBASIC application follows this structure:

1.  **Open a window**: Call `Window.Open()` at the very beginning.
2.  **Set properties**: Configure FPS, title, etc.
3.  **Main loop**: Use a `WHILE NOT Window.ShouldClose()` loop to keep the application running.
4.  **Close window**: Call `Window.Close()` after the loop exits.

```basic
; 1. Open window
Window.Open(1280, 720, "My Game")

; 2. Set properties
Window.SetFPS(60)

; 3. Main loop
WHILE NOT Window.ShouldClose()
    ; Update and draw your game here
    Render.Clear(0,0,0)
    Render.Frame()
WEND

; 4. Close window
Window.Close()
```

---

## Window Management

### `Window.Open(width, height, title)`

Initializes the window and sets its dimensions and title. This must be the first window command called. **`Window.Open` does not return a value.** If the window cannot be created, the runtime prints a short message to **stderr** and **exits the process** (so you do not need `IF NOT Window.Open ...` boilerplate in every program).

-   `width`, `height`: The dimensions of the window's client area in pixels.
-   `title`: The text to display in the window's title bar.

After a successful open, the **fullruntime (CGO)** host calls **`SetWindowSize`** again with the requested **`width`/`height`** (helps when the framebuffer size disagrees with the client area on scaled displays). The VM passes dimensions as **64-bit integers**; the host converts with **`int(width), int(height)`** for **raylib-go**, whose **`SetWindowSize`** takes Go **`int`** (not **`int32`**). That keeps the nudge compatible across 32-bit and 64-bit builds and avoids mismatched framebuffer vs client size on some Windows + Intel setups. The host then runs a short **presentation guard** before your script continues: one **`BeginDrawing` → cleared frame → `EndDrawing`** swap, a no-op toggle of **`FLAG_WINDOW_RESIZABLE`** (off again immediately) to nudge a Windows repaint on some Intel GPUs, **`SetWindowFocused`** (unless skipped; see env below), then **`PollInputEvents()`**. It then **drains the message queue** with additional **`PollInputEvents()`** calls, then presents **N** extra blank black frames (**Windows default N=2** when **`MOONBASIC_OPEN_WARMUP_FRAMES`** is unset; set to **`0`** to disable). None of this sets the internal “in frame” flag used by **`RENDER.CLEAR`** / **`RENDER.FRAME`**.

**Portability defaults (Easy Mode):** **`FLAG_WINDOW_HIGHDPI`** is **off** unless you opt in (see below). **`FLAG_MSAA_4X_HINT`** is only set if the script called **`SetMSAA`** with **2+** samples **before** **`Window.Open`**. **`SetMSAA(0)`** (as in the stock examples) keeps MSAA off.

Environment (optional):

- **`MOONBASIC_ENABLE_HIGHDPI=1`** — set **`FLAG_WINDOW_HIGHDPI`** when opening (Retina / high-DPI aware windows).
- **`MOONBASIC_SKIP_OPEN_PRESENT_KICK=1`** — skip the swap/toggle/focus step and only **`PollInputEvents()`** once inside the guard (if the kick causes trouble on a specific driver).
- **`MOONBASIC_SKIP_WINDOW_FOCUS=1`** — skip **`SetWindowFocused`** in the open guard.
- **`MOONBASIC_SAFE_WINDOW=1`** — use a **longer** post-open **`PollInputEvents`** drain (integrated / laptop “white window” experiments).
- **`MOONBASIC_OPEN_WARMUP_FRAMES=N`** — after open, present **N** blank black frames (**0–120**) before your script runs. **Unset on Windows defaults to 2**; set **`0`** to turn off. Other platforms default **0** when unset.
- **`MOONBASIC_MINIMAL_OPEN_HANDSHAKE=1`** — skip the open **presentation guard**, the long **`PollInputEvents`** drain, and the **blank-frame warmup** (only one poll runs). Use to compare behavior with older builds or strict drivers. Roughly equivalent to **`MOONBASIC_SKIP_OPEN_PRESENT_KICK=1`** + **`MOONBASIC_OPEN_WARMUP_FRAMES=0`** + no extra drain (there is no separate env for the drain count alone).

### `Window.CanOpen(width, height, title)`

Returns **`TRUE`** if the given size and non-empty title are acceptable **before** calling `Window.Open`. Use this when you need to choose a fallback resolution or show a custom error without letting `Window.Open` terminate the process.

---

### `Window.Close()`

Closes the application window and terminates the program. This is typically the last command in your script.

---

### `Window.ShouldClose()`

Returns `TRUE` if the user has attempted to close the window (e.g., by clicking the 'X' button or pressing Alt+F4). This is the standard way to control the main game loop.

---

### `Window.SetFPS(fps)`

Sets the desired frames per second (FPS) for the application. The `Time.Delta()` command will be calculated based on this target.

-   `fps`: The target FPS value (e.g., 30, 60, 144).

---

## Appearance & Position

### `Window.SetTitle(title)`

Updates the window's title text while the program is running.

-   `title`: The new title to display.

```basic
score = 100
Window.SetTitle("My Game | Score: " + STR$(score))
```

---

### `Window.SetPosition(x, y)`

Sets the position of the top-left corner of the window on the screen.

---

### `Window.SetIcon(filePath)`

Sets the window's icon from an image file. Best results with a square `.png` file (e.g., 64x64).

-   `filePath`: Path to the image file.

---

### `Window.SetOpacity(alpha)`

Sets the window's transparency.

-   `alpha`: A value from `0.0` (fully transparent) to `1.0` (fully opaque).

---

## Sizing and Monitors

### `Window.SetSize(width, height)`

Sets the dimensions of the window.

### `Window.SetMinSize(width, height)`

Sets the minimum allowed dimensions for a resizable window.

### `Window.SetMaxSize(width, height)`

Sets the maximum allowed dimensions for a resizable window.

---

### `Window.GetMonitorCount()`

Returns the number of connected monitors.

### `Window.GetMonitorWidth(monitor)` / `Window.GetMonitorHeight(monitor)`

Returns the width or height of the specified monitor in pixels.

-   `monitor`: The monitor index (0 for the primary monitor).
