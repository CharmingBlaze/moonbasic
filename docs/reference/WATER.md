# Water Commands

A **horizontal water plane** with simple wave motion, color gradients, and queries for camera/gameplay (**depth**, **underwater**). **CGO + Raylib** required.

Page shape: [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

Create with **`WATER.CREATE(...)`** (several arities — see manifest), place with **`WATER.SETPOS`**, animate with **`WATER.UPDATE(handle, dt)`**, and draw with **`WATER.DRAW`** inside **`RENDER.BEGIN3D(cam)`** / **`RENDER.END3D()`** (or **`CAMERA.BEGIN`** / **`CAMERA.END`**). **Draw order:** after opaque terrain and props, before transparent weather/particles when possible.

---

### `WATER.CREATE(...)`

Creates a subdivided water plane; returns a **handle**. **`WATER.MAKE`** is a deprecated alias. Match the arity you need from **`commands.json`**.

---

### `WATER.FREE(handle)`

Frees the water resource.

---

### `WATER.SETPOS(handle, x, y, z)`

Sets the water surface world position. **`WATER.SETPOSITION`** is a deprecated alias.

---

### `WATER.UPDATE(handle, dt)`

Advances wave simulation.

---

### `WATER.DRAW(handle)`

Renders the water surface (must be inside an active 3D camera block).

---

### `WATER.SETHEIGHT(handle, ...)` / `WATER.SETWAVE(handle, ...)` / `WATER.SETWAVEHEIGHT(handle, ...)`

Plane height and wave parameters — see manifest / runtime for overloads.

---

### `WATER.SHOW(handle, visible)`

Toggles visibility where supported.

---

### `WATER.GETWAVEY(handle, x, z)`

Returns surface **Y** including wave offset at **XZ**.

---

### `WATER.GETDEPTH(handle, x, z)`

Returns depth from surface to bed at **XZ**.

---

### `WATER.ISUNDER(handle, x, y, z)`

Returns **`TRUE`** if the point is below the animated surface.

---

### `WATER.SETSHALLOWCOLOR(handle, ...)` / `WATER.SETDEEPCOLOR(handle, ...)` / `WATER.SETCOLOR(handle, ...)`

Shallow vs deep tint and combined color — see runtime for blending.

---

## Full Example

Sketch only (camera and terrain omitted):

```basic
; x, z, width, depth, water level Y — see manifest for other WATER.CREATE overloads
water = WATER.CREATE(0, 0, 80.0, 80.0, 0.0)

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    WATER.UPDATE(water, dt)
    RENDER.CLEAR(15, 20, 35)
    RENDER.BEGIN3D(cam)
        ; ... opaque terrain ...
        WATER.DRAW(water)
    RENDER.END3D()
    RENDER.FRAME()
WEND

WATER.FREE(water)
```

**Common mistake:** Using **`TERRAIN.GETHEIGHT`** alone for water level — water has its own **Y** from **`WATER.SETPOS`**; use **`WATER.GETWAVEY`** or **`WATER.ISUNDER`** for gameplay consistency.

---

## See also

- [TERRAIN.md](TERRAIN.md)
- [SKY.md](SKY.md) — horizon tint
