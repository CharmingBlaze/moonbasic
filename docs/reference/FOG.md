# Fog Commands

Global scene fog: enable depth-based fog with color, near and far distances.

## Core Workflow

1. `FOG.ENABLE(TRUE)` — switch on fog.
2. `FOG.SETNEAR(dist)` / `FOG.SETFAR(dist)` or `FOG.SETRANGE(near, far)` — set distances.
3. `FOG.SETCOLOR(r, g, b, a)` — match the clear color for seamless blending.
4. `FOG.ENABLE(FALSE)` to disable.

---

## Enable

### `FOG.ENABLE(enabled)` 

Enables (`TRUE`) or disables (`FALSE`) scene fog globally.

---

## Distances

### `FOG.SETNEAR(distance)` 

Sets the distance at which fog begins. Objects closer than this are unaffected.

---

### `FOG.SETFAR(distance)` 

Sets the distance at which objects are completely obscured by fog.

---

### `FOG.SETRANGE(near, far)` 

Convenience — sets both near and far in one call.

---

## Color

### `FOG.SETCOLOR(r, g, b, a)` 

Sets the fog color (0–255 per channel). Match to `RENDER.CLEAR` color for seamless horizon blending.

---

## Full Example

Atmospheric fog matching the sky clear color.

```basic
WINDOW.OPEN(960, 540, "Fog Demo")
WINDOW.SETFPS(60)

FOG.ENABLE(TRUE)
FOG.SETRANGE(8.0, 40.0)
FOG.SETCOLOR(50, 70, 90, 255)

cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 2, 0)
CAMERA.SETTARGET(cam, 0, 2, -1)

; scatter some cubes in the distance
FOR i = 1 TO 12
    e = ENTITY.CREATECUBE(1.0)
    ENTITY.SETPOS(e, RNDF(-15, 15), 0.5, RNDF(-5, -35))
NEXT i

WHILE NOT WINDOW.SHOULDCLOSE()
    ENTITY.UPDATE(TIME.DELTA())
    RENDER.CLEAR(50, 70, 90)
    RENDER.BEGIN3D(cam)
        ENTITY.DRAWALL()
        DRAW3D.GRID(40, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND

FOG.ENABLE(FALSE)
WINDOW.CLOSE()
```

---

## See also

- [SKY.md](SKY.md) — skybox and sky color
- [EFFECT.md](EFFECT.md) — depth-of-field and post-process
- [RENDER.md](RENDER.md) — clear color and ambient
