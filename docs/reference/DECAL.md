# Decal Commands

World-space decal projectors: bullet holes, blood splatters, scorch marks, and other surface overlays with optional fade lifetime.

## Core Workflow

1. `DECAL.CREATE(texHandle)` — create a decal handle with a texture.
2. `DECAL.SETPOS(decal, x, y, z)` and `DECAL.SETSIZE(decal, w, h)` — place it.
3. `DECAL.SETLIFETIME(decal, seconds)` — set auto-fade time (0 = permanent).
4. `DECAL.DRAW(decal)` each frame (or let the entity system draw it).
5. `DECAL.FREE(decal)` when done.

---

## Creation

### `DECAL.CREATE(texHandle)` 

Creates a decal handle using a texture loaded with `TEXTURE.LOAD`. Returns a **decal handle**.

---

## Position & Size

### `DECAL.SETPOS(decal, x, y, z)` 

Sets the world-space position.

- *Handle shortcut*: `decal.setPos(x, y, z)`

---

### `DECAL.GETPOS(decal)` 

Returns `[x, y, z]` position array.

---

### `DECAL.SETSIZE(decal, width, height)` 

Sets the decal projection size in world units.

---

### `DECAL.GETSIZE(decal)` 

Returns `[width, height]` as a VEC2-compatible handle.

---

## Rotation

### `DECAL.SETROT(decal, yaw)` / `DECAL.SETROT(decal, pitch, yaw, roll)` 

Sets Y-axis rotation in degrees (single float) or full Euler rotation (three floats).

---

### `DECAL.GETROT(decal)` 

Returns `[pitch, yaw, roll]` rotation array.

---

## Color & Transparency

### `DECAL.SETCOLOR(decal, r, g, b)` / `DECAL.SETCOLOR(decal, r, g, b, a)` 

Sets the decal tint color (0–255).

---

### `DECAL.GETCOLOR(decal)` 

Returns `[r, g, b, a]` tint array.

---

### `DECAL.SETALPHA(decal, alpha)` 

Sets transparency 0.0 (invisible) to 1.0 (opaque).

---

### `DECAL.GETALPHA(decal)` 

Returns the current alpha value.

---

## Lifetime

### `DECAL.SETLIFETIME(decal, seconds)` 

Sets auto-fade duration. After `seconds` the decal fades out. `0` = no fade.

---

### `DECAL.GETLIFETIME(decal)` 

Returns the last lifetime value set.

---

## Draw & Free

### `DECAL.DRAW(decal)` 

Renders the decal this frame.

---

### `DECAL.FREE(decal)` 

Destroys the decal handle.

---

## Full Example

Bullet hole decals spawning at the mouse click position.

```basic
WINDOW.OPEN(960, 540, "Decal Demo")
WINDOW.SETFPS(60)

cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 5, -10)
CAMERA.SETTARGET(cam, 0, 0, 0)

tex      = TEXTURE.LOAD("assets/bullethole.png")
decals   = ARRAY.MAKE(0)

WHILE NOT WINDOW.SHOULDCLOSE()
    IF MOUSE.PRESSED(0)
        d = DECAL.CREATE(tex)
        wx = RNDF(-4, 4)
        wy = 0.01
        wz = RNDF(-4, 4)
        DECAL.SETPOS(d, wx, wy, wz)
        DECAL.SETSIZE(d, 0.4, 0.4)
        DECAL.SETLIFETIME(d, 5.0)
        ARRAY.PUSH(decals, d)
    END IF

    RENDER.CLEAR(20, 25, 35)
    RENDER.BEGIN3D(cam)
        DRAW3D.GRID(10, 1.0)
        FOR i = 0 TO ARRAY.LEN(decals) - 1
            DECAL.DRAW(ARRAY.GET(decals, i))
        NEXT i
    RENDER.END3D()
    RENDER.FRAME()
WEND

FOR i = 0 TO ARRAY.LEN(decals) - 1
    DECAL.FREE(ARRAY.GET(decals, i))
NEXT i
TEXTURE.UNLOAD(tex)
CAMERA.FREE(cam)
WINDOW.CLOSE()
```

---

## Extended Command Reference

| Command | Description |
|--------|-------------|
| `DECAL.SETPOSITION(d, x,y,z)` | Alias of `DECAL.SETPOS`. |

---

## See also

- [TEXTURE.md](TEXTURE.md) — texture loading
- [ENTITY.md](ENTITY.md) — entity-based scene objects
- [PARTICLE3D.md](PARTICLE3D.md) — particle effects
