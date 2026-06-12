# Time, Timer & Stopwatch Commands

Commands for measuring elapsed time, delta time, frame timing, countdown timers, and stopwatches. Time is essential for frame-rate-independent movement, animations, cooldowns, and gameplay scheduling.

## Core Concepts

- **Delta time (`dt`)** — The time in seconds since the last frame. Multiply all movement and animation by `dt` for consistent behavior at any frame rate.
- **Time scale** — A global multiplier on `dt`. Set to 0.5 for slow-motion, 0 to pause, 2.0 for fast-forward.
- **Timers** — Countdown objects that track whether a duration has elapsed. Two flavors: **wall-clock timers** (real time) and **simulation timers** (affected by time scale and manual `dt`).
- **Stopwatches** — Measure elapsed real time from a start point.

---

## Delta Time

### `Time.Delta()`

Returns the time in seconds since the last frame. This is the most important time function — use it for all frame-rate-independent logic.

**Returns:** `float`

**How it works:** Returns `Raylib.GetFrameTime()` multiplied by `rt.TimeScale`. If `rt.GamePaused` is `TRUE` or a hit-stop is active (`rt.HitStopEndAt > currentTime`), returns 0.

```basic
WHILE NOT Window.ShouldClose()
    dt = Time.Delta()

    ; Frame-rate independent movement
    playerX = playerX + speed * dt
    angle = angle + rotSpeed * dt

    ; ...
WEND
```

---

### `DELTATIME` / `DT()`

Easy Mode aliases for `Time.Delta()`.

---

## Wall Clock

### `MILLISECS` / `Time.Millisecs()`

Returns the number of milliseconds since the program started.

**Returns:** `int`

**How it works:** Calls `Raylib.GetTime()` and converts to milliseconds. Not affected by time scale.

```basic
startMs = MILLISECS
; ... do work ...
elapsed = MILLISECS - startMs
PRINT "Took " + STR(elapsed) + "ms"
```

---

### `FPS` / `Window.GetFPS()`

Returns the current frames per second.

**Returns:** `int`

---

## Countdown Timers (Wall-Clock)

Wall-clock timers count down in real time (not affected by time scale). Best for UI cooldowns, network timeouts, and animation delays.

### `Timer.New(durationSeconds)`

Creates a wall-clock countdown timer. The timer starts immediately.

- `durationSeconds` (float) — Duration in seconds.

**Returns:** `handle`

**How it works:** Records `time.Now() + duration` as the end time. `Timer.Finished` checks if `time.Now()` has passed that point.

```basic
respawnTimer = Timer.New(3.0)   ; 3-second respawn
```

---

### `Timer.Reset(timerHandle, durationSeconds)`

Resets and restarts the timer with a new duration.

- `durationSeconds` (float) — New duration.

```basic
Timer.Reset(respawnTimer, 3.0)
```

---

### `Timer.Finished(timerHandle)`

Returns `TRUE` if the timer's duration has elapsed.

**Returns:** `bool`

```basic
IF Timer.Finished(respawnTimer) THEN
    SpawnPlayer()
    Timer.Reset(respawnTimer, 3.0)
ENDIF
```

---

### `Timer.Remaining(timerHandle)`

Returns the remaining time in seconds. Works with both wall-clock and simulation timers.

**Returns:** `float`

```basic
remaining = Timer.Remaining(respawnTimer)
Draw.Text("Respawn in: " + STR(INT(remaining)) + "s", 10, 10, 20, 255, 255, 255, 255)
```

---

### `Timer.Free(timerHandle)`

Frees a timer handle.

---

## Simulation Timers

Simulation timers advance by manual `dt` updates, so they respect time scale, pausing, and slow-motion. Best for gameplay mechanics, ability cooldowns, and timed events.

### `Timer.Create(durationSeconds)` / `Timer.Make(durationSeconds)`

Creates a simulation timer.

- `durationSeconds` (float) — Duration in seconds.

**Returns:** `handle`

```basic
attackCooldown = Timer.Create(0.5)
```

---

### `Timer.Start(timerHandle)`

Starts (or restarts) the timer.

---

### `Timer.Stop(timerHandle)`

Pauses the timer without resetting it.

---

### `Timer.Rewind(timerHandle)`

Resets the timer's elapsed time to 0 and restarts it.

---

### `Timer.SetLoop(timerHandle, looping)`

Sets whether the timer automatically restarts when finished.

- `looping` (bool) — `TRUE` to loop.

```basic
Timer.SetLoop(attackCooldown, TRUE)   ; Auto-restart
```

---

### `Timer.Update(timerHandle, dt)`

Advances the timer by `dt` seconds. Call every frame.

- `dt` (float) — Delta time.

**How it works:** Adds `dt` to the timer's elapsed time. If elapsed exceeds the duration and looping is enabled, wraps around. Fires the "done" edge once per cycle.

```basic
WHILE NOT Window.ShouldClose()
    dt = Time.Delta()
    Timer.Update(attackCooldown, dt)

    IF Timer.Done(attackCooldown) THEN
        FireWeapon()
    ENDIF
WEND
```

---

### `Timer.Done(timerHandle)`

Returns `TRUE` on the **first frame** the timer's duration has been reached (edge-triggered). Does not repeat unless the timer loops or is rewound.

**Returns:** `bool`

---

### `Timer.Fraction(timerHandle)`

Returns the timer's progress from 0.0 (just started) to 1.0 (done).

**Returns:** `float`

```basic
; Smooth animation driven by timer progress
t = Timer.Fraction(fadeTimer)
alpha = INT(255 * t)
```

---

## Stopwatches

Stopwatches measure real elapsed time from a start point. Not affected by time scale.

### `Stopwatch.New()`

Creates a new stopwatch, started immediately.

**Returns:** `handle`

```basic
sw = Stopwatch.New()
```

---

### `Stopwatch.Reset(stopwatchHandle)`

Resets the stopwatch to 0 and restarts it.

---

### `Stopwatch.Elapsed(stopwatchHandle)`

Returns the elapsed time in seconds since the stopwatch was created or last reset.

**Returns:** `float`

```basic
elapsed = Stopwatch.Elapsed(sw)
Draw.Text("Time: " + STR(elapsed, 2) + "s", 10, 10, 24, 255, 255, 255, 255)
```

---

### `Stopwatch.Free(stopwatchHandle)`

Frees a stopwatch handle.

---

## Time Scale

### `Game.SetTimeScale(scale)` / `World.SetTimeScale(scale)`

Sets the global time multiplier. Affects `Time.Delta()` and all simulation timers.

- `scale` (float) — 1.0 = normal, 0.5 = half speed, 0.0 = paused.

```basic
; Slow motion on hit
Game.SetTimeScale(0.25)
```

---

### `Game.GetTimeScale()`

Returns the current time scale.

**Returns:** `float`

---

## Easy Mode Shortcuts

| Shortcut | Maps To |
|----------|---------|
| `DELTATIME` | `Time.Delta()` |
| `DT()` | `Time.Delta()` |
| `MILLISECS` | `Time.Millisecs()` |
| `FPS` | `Window.GetFPS()` |
| `TIMEMS` | `Time.Millisecs()` |

---

## Full Example

A game with cooldown timers, a stopwatch, and time-scale effects.

```basic
Window.Open(800, 600, "Timer Demo")
Window.SetFPS(60)

; Create a simulation cooldown timer
fireCooldown = Timer.Create(0.3)
Timer.Start(fireCooldown)
Timer.SetLoop(fireCooldown, TRUE)

; Create a wall-clock respawn timer
respawnTimer = Timer.New(5.0)

; Create a stopwatch for total game time
gameTimer = Stopwatch.New()

score = 0
canFire = TRUE

WHILE NOT Window.ShouldClose()
    dt = Time.Delta()

    ; Update simulation timer
    Timer.Update(fireCooldown, dt)

    ; Fire on click when cooldown is ready
    IF Input.MousePressed(0) AND Timer.Done(fireCooldown) THEN
        score = score + 1
    ENDIF

    ; Check respawn timer
    IF Timer.Finished(respawnTimer) THEN
        PRINT "Respawn available!"
        Timer.Reset(respawnTimer, 5.0)
    ENDIF

    ; Slow motion on TAB
    IF Input.KeyPressed(KEY_TAB) THEN
        IF Game.GetTimeScale() < 1.0 THEN
            Game.SetTimeScale(1.0)
        ELSE
            Game.SetTimeScale(0.25)
        ENDIF
    ENDIF

    ; Render
    Render.Clear(25, 25, 40)

    ; Cooldown bar
    t = Timer.Fraction(fireCooldown)
    Draw.Rectangle(10, 10, INT(200 * t), 20, 50, 200, 50, 255)
    Draw.Text("Cooldown: " + STR(INT(t * 100)) + "%", 220, 12, 16, 200, 200, 200, 255)

    ; Respawn countdown
    rem = Timer.Remaining(respawnTimer)
    Draw.Text("Respawn in: " + STR(rem, 1) + "s", 10, 40, 18, 255, 200, 100, 255)

    ; Game time
    elapsed = Stopwatch.Elapsed(gameTimer)
    Draw.Text("Game Time: " + STR(elapsed, 1) + "s", 10, 65, 18, 150, 150, 255, 255)

    ; Score & controls
    Draw.Text("Score: " + STR(score), 10, 95, 24, 255, 255, 255, 255)
    Draw.Text("Click = Score | TAB = Slow-mo", 10, 130, 16, 180, 180, 180, 255)
    Draw.Text("Time Scale: " + STR(Game.GetTimeScale()), 10, 150, 16, 180, 180, 180, 255)

    Render.Frame()
WEND

Timer.Free(fireCooldown)
Timer.Free(respawnTimer)
Stopwatch.Free(gameTimer)
Window.Close()
```

---

## See Also

- [WORLD](WORLD.md) — Time scale, hit-stop
- [MATH](MATH.md) — Easing functions for timer-driven animation
