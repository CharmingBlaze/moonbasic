# Camera2D Commands

2D scrolling camera for sprites, tilemaps, and side-scrollers. Controls world-to-screen mapping, zoom, rotation, and coordinate conversion.

## Core Workflow

1. `CAMERA2D.CREATE()` — allocate a 2D camera handle.
2. `CAMERA2D.SETTARGET(cam, x, y)` — world position the camera centres on.
3. `CAMERA2D.SETOFFSET(cam, ox, oy)` — screen-space offset (usually half screen size to centre).
4. `CAMERA2D.BEGIN(cam)` / `CAMERA2D.END()` — wrap all 2D draw calls that should be scrolled.
5. `CAMERA2D.FREE(cam)` when done.

---

## Creation & Lifetime

### `CAMERA2D.CREATE()` 

Allocates a new 2D camera handle with defaults: target `(0, 0)`, offset `(0, 0)`, zoom `1.0`, rotation `0`.

- *Handle shortcut*: use `cam = CAMERA2D.CREATE()`

---

### `CAMERA2D.FREE(cam)` 

Releases the camera handle.

- *Handle shortcut*: `cam.free()`

---

## Target & Offset

### `CAMERA2D.SETTARGET(cam, x, y)` 

Sets the world-space point the camera is centred on. Typically the player's world position.

- *Handle shortcut*: `cam.setTarget(x, y)`

---

### `CAMERA2D.SETOFFSET(cam, ox, oy)` 

Sets the screen-space offset where the target maps to. Use `WINDOW.WIDTH() / 2, WINDOW.HEIGHT() / 2` to centre the view on the target.

- *Handle shortcut*: `cam.setOffset(ox, oy)`

---

### `CAMERA2D.TARGETX(cam)` 

Returns the current target X as a float.

---

### `CAMERA2D.TARGETY(cam)` 

Returns the current target Y as a float.

---

### `CAMERA2D.GETPOS(cam)` 

Returns the current target position as a 2-element array handle `[x, y]`.

---

## Zoom

### `CAMERA2D.SETZOOM(cam, zoom)` 

Sets the zoom level. `1.0` = no zoom, `2.0` = 2× magnified, `0.5` = zoomed out.

- *Handle shortcut*: `cam.setZoom(zoom)`

---

### `CAMERA2D.ZOOMIN(cam, amount)` 

Adds `amount` to the current zoom level.

---

### `CAMERA2D.ZOOMOUT(cam, amount)` 

Subtracts `amount` from the current zoom level.

---

### `CAMERA2D.ZOOMTOMOUSE(cam, amount)` 

Adjusts zoom and offsets the target so the world point under the cursor stays fixed. Ideal for editor-style pan/zoom.

---

## Rotation

### `CAMERA2D.SETROTATION(cam, degrees)` 

Sets the camera rotation in **degrees**.

- *Handle shortcut*: `cam.setRotation(degrees)`

---

### `CAMERA2D.GETROTATION(cam)` / `CAMERA2D.ROTATION(cam)` 

Returns the current rotation in degrees.

---

## Follow

### `CAMERA2D.FOLLOW(cam, target, lerpX, lerpY)` 

Smoothly moves the camera target toward a `target` entity handle each frame. `lerpX` / `lerpY` are 0–1 smoothing factors per axis (`1.0` = instant snap).

- *Handle shortcut*: `cam.follow(target, lerpX, lerpY)`

---

## Coordinate Conversion

### `CAMERA2D.WORLDTOSCREEN(cam, worldX, worldY)` 

Converts a world-space coordinate to screen-space pixels. Returns a 2-element array handle `[sx, sy]`.

- *Handle shortcut*: `cam.worldToScreen(worldX, worldY)`

---

### `CAMERA2D.SCREENTOWORLD(cam, screenX, screenY)` 

Converts screen-space pixels to world-space coordinates. Returns a 2-element array handle `[wx, wy]`. Use for mouse picking in a scrolled scene.

- *Handle shortcut*: `cam.screenToWorld(screenX, screenY)`

---

### `CAMERA2D.GETMATRIX(cam)` 

Returns the 4×4 camera matrix handle used internally. Useful for custom batch rendering.

---

## Render Pair

### `CAMERA2D.BEGIN()` / `CAMERA2D.BEGIN(cam)` 

Starts 2D camera-space rendering. All `DRAW.*` and `SPRITE.*` calls between `BEGIN` and `END` are transformed by the camera. Call with no argument to use the most recently set active camera, or pass `cam` to activate a specific one.

---

### `CAMERA2D.END()` 

Ends 2D camera-space rendering. Always pair with `CAMERA2D.BEGIN`.

---

## Full Example

A side-scroller camera that follows a sprite with smooth lerp.

```basic
WINDOW.OPEN(800, 450, "Camera2D Demo")
WINDOW.SETFPS(60)

cam = CAMERA2D.CREATE()
CAMERA2D.SETOFFSET(cam, 400, 225)
CAMERA2D.SETZOOM(cam, 1.0)

px = 400.0
py = 225.0

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()

    IF INPUT.KEYDOWN(KEY_RIGHT) THEN px = px + 200 * dt
    IF INPUT.KEYDOWN(KEY_LEFT)  THEN px = px - 200 * dt
    IF INPUT.KEYDOWN(KEY_UP)    THEN py = py - 200 * dt
    IF INPUT.KEYDOWN(KEY_DOWN)  THEN py = py + 200 * dt

    ; smooth follow
    tx = CAMERA2D.TARGETX(cam)
    ty = CAMERA2D.TARGETY(cam)
    CAMERA2D.SETTARGET(cam, tx + (px - tx) * 0.1, ty + (py - ty) * 0.1)

    RENDER.CLEAR(30, 30, 50)
    CAMERA2D.BEGIN(cam)
        ; draw world grid
        FOR i = 0 TO 20
            DRAW.LINE(i * 100, 0, i * 100, 2000, 50, 50, 70, 255)
            DRAW.LINE(0, i * 100, 2000, i * 100, 50, 50, 70, 255)
        NEXT i
        ; draw player
        DRAW.RECT(INT(px) - 16, INT(py) - 16, 32, 32, 80, 160, 255, 255)
    CAMERA2D.END()
    RENDER.FRAME()
WEND

CAMERA2D.FREE(cam)
WINDOW.CLOSE()
```

---

## Extended Command Reference

| Command | Description |
|--------|-------------|
| `CAMERA2D.MAKE(...)` | Deprecated alias of `CAMERA2D.CREATE`. |

---

## See also

- [CAMERA.md](CAMERA.md) — 3D camera
- [SPRITE.md](SPRITE.md) — sprite drawing inside camera space
- [TILEMAP.md](TILEMAP.md) — tile maps with 2D camera
- [INPUT.md](INPUT.md) — mouse world-position via `CAMERA2D.SCREENTOWORLD`
