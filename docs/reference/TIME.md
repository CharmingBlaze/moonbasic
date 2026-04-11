# Time Commands

Commands for getting the current time, date, and measuring elapsed time.

## Core Concepts

-   **Delta Time**: `Time.Delta()` is the most important command for game development. It gives you the time elapsed since the last frame, allowing you to create movement and logic that is independent of the frame rate.
-   **Program Time**: `Time.Get()` or `TIMER()` gives you the total elapsed time since the program started.
-   **Wall-Clock Time**: `DATE()`, `TIME()`, `YEAR()`, etc., give you the real-world time from the system clock.

---

### `Time.Delta()`
Returns seconds since last frame.

### `Time.Get()`
Returns total elapsed seconds since start.

### `Time.TickCount()`
Returns total elapsed milliseconds since start.

---

### `Date()` / `Time()` / `DateTime()`
Returns formatted wall-clock strings.

### `Timestamp()`
Returns Unix epoch seconds.

---

## Full Example: A Simple Stopwatch

```basic
Window.Open(800, 600, "Stopwatch Example")
Window.SetFPS(60)

start_time = 0.0
stop_time = 0.0
running = FALSE

WHILE NOT Window.ShouldClose()
    ; --- LOGIC ---
    IF Input.KeyPressed(KEY_SPACE) THEN
        IF running THEN
            ; Stop the timer
            stop_time = Time.Get()
            running = FALSE
        ELSE
            ; Start or reset the timer
            start_time = Time.Get()
            stop_time = 0.0
            running = TRUE
        ENDIF
    ENDIF

    ; --- DRAWING ---
    elapsed_time = 0.0
    IF running THEN
        elapsed_time = Time.Get() - start_time
    ELSE
        elapsed_time = stop_time - start_time
    ENDIF

    Render.Clear(0,0,0)
    Camera2D.Begin()
        Draw.Text("Press SPACE to start/stop", 210, 150, 20, 150, 150, 150, 255)
        Draw.Text(FORMAT(elapsed_time, "%.2f"), 300, 250, 60, 100, 200, 255, 255)
    Camera2D.End()
    Render.Frame()
WEND

Window.Close()
```
