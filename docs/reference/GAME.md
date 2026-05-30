# Game Commands

Convenience namespace aggregating common game-loop helpers: screen dimensions, delta time, mouse/keyboard input shortcuts, gamepad, time scale, screen flash, debug rect, and miscellaneous game utilities.

`GAME.*` commands are thin wrappers over core subsystems — prefer the specific namespace (`WINDOW.*`, `MOUSE.*`, `INPUT.*`, `TIME.*`) for new code. Use `GAME.*` for concise game-jam style scripts.

## Core Workflow

Most `GAME.*` commands are read-only queries called each frame inside the game loop.

```basic
WHILE NOT WINDOW.SHOULDCLOSE()
    dt = GAME.DT()
    IF GAME.KEYDOWN(KEY_SPACE) THEN ...
    IF GAME.MLEFT() THEN ...
    RENDER.CLEAR(20, 25, 35)
    RENDER.FRAME()
WEND
```

---

## Screen

### `GAME.SCREENW()` / `GAME.SCREENH()` 

Returns the window width or height in pixels.

---

### `GAME.SCREENCX()` / `GAME.SCREENCY()` 

Returns the screen centre X or Y as a float.

---

## Time

### `GAME.DT()` 

Returns delta time in seconds (alias of `TIME.DELTA()`).

---

### `GAME.FPS()` 

Returns the current measured frames-per-second.

---

### `GAME.SETTIMESCALE(scale)` 

Scales delta time globally. `1.0` = normal, `0.5` = slow-mo, `0.0` = paused.

---

### `GAME.GETTIMESCALE()` 

Returns the current time scale.

---

### `GAME.SLOWMOTION(scale, duration)` 

Temporarily sets time scale to `scale` for `duration` seconds, then restores.

---

### `GAME.SETPAUSE(enabled)` 

Pauses (`1`) or resumes (`0`) the game time.

---

## Mouse

### `GAME.MX()` / `GAME.MY()` 

Mouse X/Y position (aliases of `MOUSE.X`, `MOUSE.Y`).

---

### `GAME.MOUSEX()` / `GAME.MOUSEY()` 

Aliases of `GAME.MX` / `GAME.MY`.

---

### `GAME.MDX()` / `GAME.MDY()` 

Mouse delta X/Y this frame.

---

### `GAME.MOUSEXSPEED()` / `GAME.MOUSEYSPEED()` 

Returns the mouse movement speed on X or Y (pixels/second equivalent). Use for speed-sensitive camera sensitivity scaling rather than raw delta.

---

### `GAME.MWHEEL()` 

Mouse scroll wheel delta.

---

### `GAME.MLEFT()` / `GAME.MRIGHT()` / `GAME.MMIDDLE()` 

Returns `TRUE` while the left/right/middle button is held.

---

### `GAME.MLEFTPRESSED()` / `GAME.MRIGHTPRESSED()` 

Returns `TRUE` on the first frame the button is pressed.

---

### `GAME.ISCURSORONSCREEN()` 

Returns `TRUE` if the cursor is within the window.

---

## Keyboard

### `GAME.KEYDOWN(keyCode)` 

Returns `TRUE` while `keyCode` is held.

---

### `GAME.KEYPRESSED(keyCode)` 

Returns `TRUE` on the first frame `keyCode` is pressed.

---

### `GAME.KEYRELEASED(keyCode)` 

Returns `TRUE` on the first frame `keyCode` is released.

---

### `GAME.KEYHIT(key)` 

Returns `TRUE` if the key was hit (alias of `GAME.KEYPRESSED`).

---

### `GAME.KEYCHAR()` 

Returns the character code of the last key pressed.

---

### `GAME.ANYKEY()` 

Returns `TRUE` if any key was pressed this frame.

---

## Gamepad

### `GAME.JOYX()` / `GAME.JOYY()` 

Returns the left stick X/Y axis values.

---

### `GAME.JOYBUTTON(button)` 

Returns `TRUE` if gamepad `button` is pressed.

---

### `GAME.ISGAMEPADAVAILABLE(id)` 

Returns `TRUE` if gamepad `id` is connected.

---

### `GAME.GETGAMEPADNAME(id)` 

Returns the gamepad name string.

---

## Camera Orbit Helpers

### `GAME.ORBITDISTDELTA(sensitivity)` 

Returns a scroll-wheel-based orbit distance delta. Use to zoom a `CAMERA.SETORBIT` distance each frame.

---

### `GAME.ORBITPITCHDELTA(sensitivity)` 

Returns a mouse-Y-based orbit pitch delta.

---

### `GAME.ORBITYAWDELTA(sensitivity, threshold, btn, btnAlt, speed)` 

Returns a mouse-X-based orbit yaw delta when button is held.

---

## Screen Effects

### `GAME.SCREENFLASH(r, g, b, a)` 

Triggers a one-frame screen flash overlay with the given color.

---

### `GAME.DRAWSCREENFLASH()` 

Draws the current screen flash overlay. Call once per frame in the render loop.

---

## Debug

### `GAME.DEBUGRECT(x, y, w, h, r, g, b, a)` 

Draws a debug rectangle overlay (visible in debug builds).

---

## Volume

### `GAME.SETMASTERVOLUME(volume)` 

Sets master audio volume 0.0–1.0.

---

### `GAME.GETMASTERVOLUME()` 

Returns current master volume.

---

## Misc

### `GAME.ENDGAME()` 

Signals the game loop to exit cleanly.

---

### `GAME.SETAESTHETIC(mode)` 

Sets a global rendering aesthetic preset (see render docs).

---

### `GAME.BURSTSPAWN(template, count, x, y, z)` 

Spawns `count` copies of `template` entity around position.

---

### `GAME.MAKEFLOATARRAY(size)` 

Returns a new handle to a zero-initialised float array of `size` elements.

---

### `GAME.SPRITETILEBRIDGE(...)` 

Internal bridge for sprite tile mapping.

---

## Full Example

Minimal game-jam style loop using only `GAME.*`.

```basic
WINDOW.OPEN(800, 600, "Game Demo")
WINDOW.SETFPS(60)

px = GAME.SCREENCX()
py = GAME.SCREENCY()

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = GAME.DT()

    IF GAME.KEYDOWN(KEY_RIGHT) THEN px = px + 200 * dt
    IF GAME.KEYDOWN(KEY_LEFT)  THEN px = px - 200 * dt
    IF GAME.KEYDOWN(KEY_DOWN)  THEN py = py + 200 * dt
    IF GAME.KEYDOWN(KEY_UP)    THEN py = py - 200 * dt

    IF GAME.MLEFTPRESSED() THEN
        GAME.SCREENFLASH(255, 255, 255, 200)
    END IF

    RENDER.CLEAR(20, 20, 40)
    DRAW.CIRCLE(INT(px), INT(py), 16, 80, 160, 255, 255)
    DRAW.TEXT("FPS: " + STR(GAME.FPS()), 10, 10, 18, 200, 200, 200, 255)
    GAME.DRAWSCREENFLASH()
    RENDER.FRAME()
WEND

WINDOW.CLOSE()
```

---

## See also

- [INPUT.md](INPUT.md) — full unified input (keys, mouse, gamepad, touch)
- [MOUSE.md](MOUSE.md) — mouse-only commands
- [TIME.md](TIME.md) — `TIME.DELTA`, time utilities
- [WINDOW.md](WINDOW.md) — window size, title, close
