# Time Commands

Commands for getting the current time, date, and measuring elapsed time.

## Core Concepts

-   **Delta Time**: `Time.Delta()` is the most important command for game development. It gives you the time elapsed since the last frame, allowing you to create movement and logic that is independent of the frame rate.
-   **Program Time**: `Time.Get()` or `TIMER()` gives you the total elapsed time since the program started.
-   **Wall-Clock Time**: `DATE$()`, `TIME$()`, `YEAR()`, etc., give you the real-world time from the system clock.

---

## Game Time

### `Time.Delta()`

Returns the time in seconds that has passed since the last frame.

By default the value is **clamped** to at most **0.05** seconds (so a hitch longer than ~20 FPS does not apply a huge delta to physics). **`GAME.DT()`** / **`DT()`** use the same clamped value.

- **`Time.SetMaxDelta(seconds#)`** — set a different cap; use **`<= 0`** to **disable** clamping if you manage timing yourself.

```basic
; Move the player at 200 pixels per second, regardless of FPS
speed# = 200.0
player_x# = player_x# + speed# * Time.Delta()
```

### `Time.Get()` / `TIMER()`

Returns the total time in seconds since the program started.

### `TICKCOUNT()`

Returns the number of milliseconds since the program started.

---

## System Time

### `DATE$()` / `TIME$()`

Returns the current system date or time as a formatted string.

### `YEAR()` / `MONTH()` / `DAY()`

Returns the individual components of the current date.

### `HOUR()` / `MINUTE()` / `SECOND()`

Returns the individual components of the current time.

---

## Full Example: A Simple Stopwatch

```basic
Window.Open(800, 600, "Stopwatch Example")
Window.SetFPS(60)

start_time# = 0.0
stop_time# = 0.0
running? = FALSE

WHILE NOT Window.ShouldClose()
    ; --- LOGIC ---
    IF Input.KeyPressed(KEY_SPACE) THEN
        IF running? THEN
            ; Stop the timer
            stop_time# = Time.Get()
            running? = FALSE
        ELSE
            ; Start or reset the timer
            start_time# = Time.Get()
            stop_time# = 0.0
            running? = TRUE
        ENDIF
    ENDIF

    ; --- DRAWING ---
    elapsed_time# = 0.0
    IF running? THEN
        elapsed_time# = Time.Get() - start_time#
    ELSE
        elapsed_time# = stop_time# - start_time#
    ENDIF

    Render.Clear(0,0,0)
    Camera2D.Begin()
        Draw.Text("Press SPACE to start/stop", 210, 150, 20, 150, 150, 150, 255)
        Draw.Text(FORMAT$(elapsed_time#, "%.2f"), 300, 250, 60, 100, 200, 255, 255)
    Camera2D.End()
    Render.Frame()
WEND

Window.Close()
```
