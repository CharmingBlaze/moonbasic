# DrawPrim3D Commands

Retained-mode 3D primitive draw objects. Configure once, update each frame, draw with a single call. Useful for persistent debug overlays, gizmos, and visualisers.

## Core Workflow

1. Obtain a primitive handle (from a `DRAW3D.*` constructor or configure after allocation).
2. Set properties: `DRAWPRIM3D.POS`, `DRAWPRIM3D.SIZE`, `DRAWPRIM3D.COLOR`, etc.
3. `DRAWPRIM3D.DRAW(prim)` each frame inside `RENDER.BEGIN3D` / `RENDER.END3D`.
4. `DRAWPRIM3D.FREE(prim)` when done.

---

## Properties

### `DRAWPRIM3D.POS(prim, x, y, z)` 

Sets the world-space position.

---

### `DRAWPRIM3D.SIZE(prim, size)` 

Sets the primary dimension (radius, half-extent, or scale).

---

### `DRAWPRIM3D.COLOR(prim, r, g, b, a)` / `DRAWPRIM3D.COL(prim, r, g, b, a)` 

Sets the draw color (0–255).

---

### `DRAWPRIM3D.WIRE(prim, enabled)` 

When `TRUE` draws wireframe; `FALSE` draws solid.

---

### `DRAWPRIM3D.RADIUS(prim, radius)` 

Sets the radius for sphere and capsule primitives.

---

### `DRAWPRIM3D.ENDPOINT(prim, x, y, z)` 

Sets the end point for line primitives.

---

### `DRAWPRIM3D.CYL(prim, topRadius, bottomRadius, height)` 

Configures cylinder/frustum dimensions.

---

### `DRAWPRIM3D.BBOX(prim, minX, minY, minZ, maxX, maxY, maxZ)` 

Configures a bounding-box wireframe.

---

### `DRAWPRIM3D.SLICES(prim, slices)` / `DRAWPRIM3D.RINGS(prim, rings)` 

Sets the tessellation resolution for sphere/cylinder primitives.

---

### `DRAWPRIM3D.GRID(prim, cells, cellSize)` 

Configures a grid overlay with `cells` divisions and `cellSize` spacing.

---

### `DRAWPRIM3D.SETRAY(prim, rayHandle)` 

Binds a `RAY.*` handle to the primitive for ray-visualisation.

---

### `DRAWPRIM3D.SETTEXTURE(prim, texHandle)` 

Assigns a texture to the primitive.

---

### `DRAWPRIM3D.SRCTEX(prim, u, v, w, h)` 

Sets the UV source rectangle for a textured primitive.

---

## Draw & Free

### `DRAWPRIM3D.DRAW(prim)` 

Draws the primitive this frame.

---

### `DRAWPRIM3D.FREE(prim)` 

Destroys the primitive handle.

---

## Full Example

A pulsing sphere debug overlay tracking an entity.

```basic
WINDOW.OPEN(960, 540, "DrawPrim3D Demo")
WINDOW.SETFPS(60)

cam  = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 5, -10)
CAMERA.SETTARGET(cam, 0, 0, 0)

cube = ENTITY.CREATECUBE(1.0)

sphere = DRAW3D.SPHERE(0, 0, 0, 0.8, 80, 200, 255, 180)

t = 0.0
WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    t  = t + dt
    ENTITY.TURN(cube, 0, 45 * dt, 0)
    ENTITY.UPDATE(dt)

    pulse = 0.7 + SIN(t * 3) * 0.2
    DRAWPRIM3D.SIZE(sphere, pulse)
    DRAWPRIM3D.COLOR(sphere, 80, INT(150 + SIN(t) * 80), 255, 180)

    RENDER.CLEAR(20, 25, 35)
    RENDER.BEGIN3D(cam)
        ENTITY.DRAWALL()
        DRAWPRIM3D.DRAW(sphere)
        DRAW3D.GRID(20, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND

DRAWPRIM3D.FREE(sphere)
ENTITY.FREE(cube)
CAMERA.FREE(cam)
WINDOW.CLOSE()
```

---

## See also

- [DRAW3D.md](DRAW3D.md) — immediate-mode 3D drawing
- [DRAWPRIM2D.md](DRAWPRIM2D.md) — 2D retained primitives
- [DEBUG.md](DEBUG.md) — debug overlays
