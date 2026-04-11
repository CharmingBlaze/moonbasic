# Window Commands

Commands for managing the application window.

## Core Workflow

A typical moonBASIC application follows this structure:

1.  **Open a window**: Call `Window.Open(width, height, title)` at the very beginning.
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

Opens the application window with specified width, height, and title string.

-   `width`, `height`: The dimensions of the window's client area in pixels.
-   `title`: The text to display in the window's title bar.

### `Window.Close()` 

Closes the window and terminates the application.

### `Window.ShouldClose()` 

Returns `TRUE` if the user has requested to close the window.

### `Window.SetFPS(fps)` 

Sets the target frames per second.

-   `fps`: The target FPS value (e.g., 30, 60, 144).

### `Window.SetTargetFPS(fps)` 

Alias for `Window.SetFPS(fps)`.

## Appearance & Position

### `Window.SetTitle(title)` 

Updates the window title at runtime.

-   `title`: The new title to display.

```basic
score = 100
Window.SetTitle("My Game | Score: " + STR$(score))
```

### `Window.SetPosition(x, y)` 

Sets the position of the top-left corner of the window on the screen.

### `Window.SetIcon(filePath)` 

Sets the window's icon from an image file. Best results with a square `.png` file (e.g., 64x64).

-   `filePath`: Path to the image file.

### `Window.SetOpacity(alpha)` 

Sets the window's transparency.

### `Window.SetSize(w, h)` 

Resizes the window to the specified dimensions in pixels.

### `Window.GetPositionX()` / `Window.GetPositionY()` 

Returns the current screen position of the top-left corner of the window.

### `Window.DPISCALE()` 

Returns the global DPI scale factor (float64) for high-DPI displays.

---

### `Window.GetMonitorCount()` 

Returns the number of connected monitors.

### `Window.GetMonitorWidth(monitor)` / `Window.GetMonitorHeight(monitor)` 

Returns the width or height of the specified monitor in pixels.

-   `monitor`: The monitor index (0 for the primary monitor).

