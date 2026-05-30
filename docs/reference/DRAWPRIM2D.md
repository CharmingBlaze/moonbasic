# DrawPrim2D Commands

Retained-mode 2D primitive draw objects. Configure once, update properties each frame, draw with a single call. Useful for shapes that need per-frame position/color updates without rebuilding draw calls from scratch.

## Core Workflow

1. Create a primitive (via `DRAW.*` constructors that return a handle, or configure after creation).
2. Set properties: `DRAWPRIM2D.POS`, `DRAWPRIM2D.SIZE`, `DRAWPRIM2D.COLOR`, etc.
3. `DRAWPRIM2D.DRAW(prim)` each frame.
4. `DRAWPRIM2D.FREE(prim)` when done.

---

## Properties

### `DRAWPRIM2D.POS(prim, x, y)` 

Sets the draw position in screen pixels.

---

### `DRAWPRIM2D.SIZE(prim, size)` 

Sets the primary size (radius for circles/rings, half-extent for rects).

---

### `DRAWPRIM2D.COLOR(prim, r, g, b, a)` / `DRAWPRIM2D.COL(prim, r, g, b, a)` 

Sets the fill color (0–255 per channel).

---

### `DRAWPRIM2D.OUTLINE(prim, enabled)` 

Sets whether to draw filled (`FALSE`) or outline-only (`TRUE`).

---

### `DRAWPRIM2D.P2(prim, x, y)` 

Sets a second point for line or triangle primitives.

---

### `DRAWPRIM2D.P3(prim, x, y)` 

Sets a third point for triangle primitives.

---

### `DRAWPRIM2D.RING(prim, innerRadius, outerRadius, startAngle, endAngle)` 

Configures a ring/arc sector.

---

### `DRAWPRIM2D.SEGS(prim, segments)` / `DRAWPRIM2D.SIDES(prim, sides)` 

Sets the number of segments (circle smoothness) or sides (polygon).

---

### `DRAWPRIM2D.ROT(prim, degrees)` 

Sets the rotation angle in degrees.

---

### `DRAWPRIM2D.THICK(prim, thickness)` 

Sets line thickness in pixels.

---

## Draw & Free

### `DRAWPRIM2D.DRAW(prim)` 

Draws the configured primitive this frame.

---

### `DRAWPRIM2D.FREE(prim)` 

Destroys the primitive handle.

---

## Full Example

A rotating polygon with a color pulse.

```basic
WINDOW.OPEN(800, 600, "DrawPrim2D Demo")
WINDOW.SETFPS(60)

poly = DRAW.POLYGON(400, 300, 80, 6, 255, 100, 50, 255)   ; returns handle

t = 0.0
WHILE NOT WINDOW.SHOULDCLOSE()
    t = t + TIME.DELTA()
    g = INT(100 + SIN(t * 2) * 100)
    DRAWPRIM2D.POS(poly, 400, 300)
    DRAWPRIM2D.ROT(poly, t * 30)
    DRAWPRIM2D.COLOR(poly, 255, g, 50, 255)
    DRAWPRIM2D.SIDES(poly, 6)

    RENDER.CLEAR(15, 15, 30)
    DRAWPRIM2D.DRAW(poly)
    RENDER.FRAME()
WEND

DRAWPRIM2D.FREE(poly)
WINDOW.CLOSE()
```

---

## See also

- [DRAW2D.md](DRAW2D.md) — immediate-mode 2D drawing
- [DRAWPRIM3D.md](DRAWPRIM3D.md) — 3D retained primitives
