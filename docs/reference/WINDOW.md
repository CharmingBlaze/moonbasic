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

### `Window.Open(width, height, title$)`

Initializes the window and sets its dimensions and title. This must be the first window command called. **`Window.Open` does not return a value.** If the window cannot be created, the runtime prints a short message to **stderr** and **exits the process** (so you do not need `IF NOT Window.Open ...` boilerplate in every program).

-   `width`, `height`: The dimensions of the window's client area in pixels.
-   `title$`: The text to display in the window's title bar.

### `Window.CanOpen(width, height, title$)`

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

### `Window.SetTitle(title$)`

Updates the window's title text while the program is running.

-   `title$`: The new title to display.

```basic
score = 100
Window.SetTitle("My Game | Score: " + STR$(score))
```

---

### `Window.SetPosition(x, y)`

Sets the position of the top-left corner of the window on the screen.

---

### `Window.SetIcon(filePath$)`

Sets the window's icon from an image file. Best results with a square `.png` file (e.g., 64x64).

-   `filePath$`: Path to the image file.

---

### `Window.SetOpacity(alpha#)`

Sets the window's transparency.

-   `alpha#`: A value from `0.0` (fully transparent) to `1.0` (fully opaque).

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
