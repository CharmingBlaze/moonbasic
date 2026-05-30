# Ray2D Commands

2D scalar ray intersection tests against circles, axis-aligned rectangles, and line segments. All arguments are raw floats — no heap handles required.

## Core Workflow

1. Call a `RAY2D.HIT*_HIT` function with the ray and target parameters.
2. If it returns `TRUE`, read `_DISTANCE`, `_POINTX`, `_POINTY` with the same arguments.
3. No cleanup needed — all scalar, no handles.

Ray convention: `(ox, oy)` = origin, `(dx, dy)` = direction (not necessarily normalised).

---

## Circle Tests

### `RAY2D.HITCIRCLE_HIT(ox, oy, dx, dy, cx, cy, radius)` 

Returns `TRUE` if the ray hits the circle at `(cx, cy)` with `radius`.

---

### `RAY2D.HITCIRCLE_DISTANCE(ox, oy, dx, dy, cx, cy, radius)` 

Returns the distance along the ray to the hit point. `0` if no hit.

---

### `RAY2D.HITCIRCLE_POINTX(ox, oy, dx, dy, cx, cy, radius)` / `RAY2D.HITCIRCLE_POINTY(...)` 

Returns the X or Y coordinate of the hit point.

---

## Rectangle Tests

### `RAY2D.HITRECT_HIT(ox, oy, dx, dy, minX, minY, maxX, maxY)` 

Returns `TRUE` if the ray hits the axis-aligned rectangle from `(minX, minY)` to `(maxX, maxY)`.

---

### `RAY2D.HITRECT_DISTANCE(ox, oy, dx, dy, minX, minY, maxX, maxY)` 

Returns the distance to the hit.

---

### `RAY2D.HITRECT_POINTX(...)` / `RAY2D.HITRECT_POINTY(...)` 

Returns the X or Y coordinate of the hit point.

---

## Segment Tests

### `RAY2D.HITSEGMENT_HIT(ox, oy, dx, dy, x1, y1, x2, y2)` 

Returns `TRUE` if the ray intersects the line segment from `(x1, y1)` to `(x2, y2)`.

---

### `RAY2D.HITSEGMENT_DISTANCE(ox, oy, dx, dy, x1, y1, x2, y2)` 

Returns the distance along the ray to the intersection.

---

### `RAY2D.HITSEGMENT_POINTX(...)` / `RAY2D.HITSEGMENT_POINTY(...)` 

Returns the X or Y coordinate of the intersection point.

---

## Full Example

A 2D laser beam hitting walls.

```basic
WINDOW.OPEN(800, 600, "Ray2D Demo")
WINDOW.SETFPS(60)

; wall segment
wx1 = 400  wy1 = 100
wx2 = 400  wy2 = 500

ox = 100.0
oy = 300.0

WHILE NOT WINDOW.SHOULDCLOSE()
    mx = MOUSE.X()
    my = MOUSE.Y()

    dx = mx - ox
    dy = my - oy
    len = SQR(dx * dx + dy * dy)
    IF len > 0 THEN dx = dx / len : dy = dy / len

    hit = RAY2D.HITSEGMENT_HIT(ox, oy, dx, dy, wx1, wy1, wx2, wy2)

    RENDER.CLEAR(10, 10, 20)

    ; wall
    DRAW.LINE(wx1, wy1, wx2, wy2, 120, 120, 180, 255)

    ; laser
    IF hit THEN
        hx = RAY2D.HITSEGMENT_POINTX(ox, oy, dx, dy, wx1, wy1, wx2, wy2)
        hy = RAY2D.HITSEGMENT_POINTY(ox, oy, dx, dy, wx1, wy1, wx2, wy2)
        DRAW.LINE(INT(ox), INT(oy), INT(hx), INT(hy), 255, 80, 80, 255)
        DRAW.CIRCLE(INT(hx), INT(hy), 6, 255, 200, 60, 255)
    ELSE
        DRAW.LINE(INT(ox), INT(oy), mx, my, 80, 80, 200, 200)
    END IF

    DRAW.CIRCLE(INT(ox), INT(oy), 8, 80, 200, 80, 255)
    RENDER.FRAME()
WEND

WINDOW.CLOSE()
```

---

## See also

- [RAY.md](RAY.md) — 3D ray tests
- [COLLISION.md](COLLISION.md) — entity overlap
- [BBOX.md](BBOX.md) — axis-aligned bounding boxes
