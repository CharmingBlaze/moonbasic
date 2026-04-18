# Light2D Commands

2D point light handles for sprite lighting, dynamic shadows, and atmosphere in top-down and side-scrolling scenes.

## Core Workflow

1. `LIGHT2D.CREATE()` — allocate a light handle.
2. `LIGHT2D.SETPOS(light, x, y)` — position in world space.
3. `LIGHT2D.SETCOLOR(light, r, g, b, a)` and `LIGHT2D.SETRADIUS(light, r)` — configure appearance.
4. Draw your scene between `CAMERA2D.BEGIN` / `CAMERA2D.END` — the 2D lighting shader uses active `LIGHT2D` handles automatically.
5. `LIGHT2D.FREE(light)` when done.

---

## Creation

### `LIGHT2D.CREATE()` 

Creates a 2D point light with defaults: white color, radius 200, intensity 1.0. Returns a **light handle**.

---

## Position

### `LIGHT2D.SETPOS(light, x, y)` 

Sets the world-space position of the light.

- *Handle shortcut*: `light.setPos(x, y)`

---

### `LIGHT2D.GETPOS(light)` 

Returns `[x, y]` position array.

---

## Appearance

### `LIGHT2D.SETCOLOR(light, r, g, b, a)` 

Sets the light color (0–255 per channel). Alpha scales intensity.

- *Handle shortcut*: `light.setColor(r, g, b, a)`

---

### `LIGHT2D.GETCOLOR(light)` 

Returns the color as a color handle.

---

### `LIGHT2D.SETRADIUS(light, radius)` 

Sets the light falloff radius in pixels (world units).

---

### `LIGHT2D.GETRADIUS(light)` 

Returns the radius.

---

### `LIGHT2D.SETINTENSITY(light, intensity)` 

Sets the brightness multiplier (0.0–2.0+).

---

### `LIGHT2D.GETINTENSITY(light)` 

Returns the intensity.

---

## Lifetime

### `LIGHT2D.FREE(light)` 

Destroys the light handle.

---

## Full Example

A torch that follows the player in a dark dungeon.

```basic
WINDOW.OPEN(800, 600, "Light2D Demo")
WINDOW.SETFPS(60)

cam = CAMERA2D.CREATE()
CAMERA2D.SETOFFSET(cam, 400, 300)

px = 400.0
py = 300.0

torch = LIGHT2D.CREATE()
LIGHT2D.SETCOLOR(torch, 255, 180, 80, 255)
LIGHT2D.SETRADIUS(torch, 250)
LIGHT2D.SETINTENSITY(torch, 1.2)

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    IF INPUT.KEYDOWN(KEY_RIGHT) THEN px = px + 150 * dt
    IF INPUT.KEYDOWN(KEY_LEFT)  THEN px = px - 150 * dt
    IF INPUT.KEYDOWN(KEY_DOWN)  THEN py = py + 150 * dt
    IF INPUT.KEYDOWN(KEY_UP)    THEN py = py - 150 * dt

    CAMERA2D.SETTARGET(cam, px, py)
    LIGHT2D.SETPOS(torch, px, py)

    RENDER.CLEAR(5, 5, 10)
    CAMERA2D.BEGIN(cam)
        ; draw world tiles here
        DRAW.RECT(INT(px) - 8, INT(py) - 8, 16, 16, 200, 200, 255, 255)
    CAMERA2D.END()
    RENDER.FRAME()
WEND

LIGHT2D.FREE(torch)
CAMERA2D.FREE(cam)
WINDOW.CLOSE()
```

---

## Extended Command Reference

| Command | Description |
|--------|-------------|
| `LIGHT2D.MAKE(...)` | Deprecated alias of `LIGHT2D.CREATE`. |
| `LIGHT2D.SETPOSITION(l, x, y)` | Alias of `LIGHT2D.SETPOS`. |

---

## See also

- [LIGHT.md](LIGHT.md) — 3D lighting
- [SPRITE.md](SPRITE.md) — sprite rendering with 2D lighting
- [CAMERA2D.md](CAMERA2D.md) — 2D camera
