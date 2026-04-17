# Water Commands

A **horizontal water plane** with simple wave motion, color gradients, and queries for camera/gameplay (**depth**, **underwater**). **CGO + Raylib** required.

Page shape: [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

Create with **`WATER.CREATE(...)`** (arity matches your overload — see **`commands.json`**), place with **`WATER.SETPOS`**, advance waves with **`WATER.UPDATE(dt)`** (single **`dt`**; updates **all** water instances), and draw with **`WATER.DRAW`** inside **`RENDER.BEGIN3D(cam)`** / **`RENDER.END3D()`** (or **`CAMERA.BEGIN`** / **`CAMERA.END`**). **Draw order:** after opaque terrain and props, before transparent weather/particles when possible.

### Method chaining 

Mutating builtins take the water handle as the first argument and return that handle on success, so you can chain setters and **`DRAW`**, e.g. **`water.setPos(x, y, z).draw()`** (same as **`WATER.SETPOS`** then **`WATER.DRAW`**). **`water.update(dt)`** is also chainable per-handle. The global **`WATER.UPDATE(dt)`** still advances all instances. See [DRAW3D.md](DRAW3D.md).

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

### `WATER.UPDATE(dt)` / `water.update(dt)` 

Advances wave simulation. Global form advances **all** active instances; handle form advances only that instance.

---

### `WATER.AUTOPHYSICS(handle)` 

Enables automatic physics integration for this water body (buoyancy/wave coupling). Handle shortcut: `water.autoPhysics()`.

---

### `WATER.DRAW(handle)` 

Renders the water surface (must be inside an active 3D camera block).

---

### `WATER.SETWAVE(handle, speed, height)` / `WATER.SETWAVEHEIGHT(handle, height)` 

Wave frequency/amplitude-style parameters — see **`commands.json`** / runtime for exact semantics.

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
    WATER.UPDATE(dt)
    RENDER.CLEAR(15, 20, 35)
    RENDER.BEGIN3D(cam)
        ; ... opaque terrain ...
        WATER.DRAW(water)
        ; or: water.draw()
        ; or chain: water.setPos(0, 0, 0).draw()
    RENDER.END3D()
    RENDER.FRAME()
WEND

WATER.FREE(water)
```

**Common mistake:** Using **`TERRAIN.GETHEIGHT`** alone for water level — water has its own **Y** from **`WATER.SETPOS`**; use **`WATER.GETWAVEY`** or **`WATER.ISUNDER`** for gameplay consistency.

---

## See also

- [DRAW3D.md](DRAW3D.md) — 3D pass + handle chaining
- [TERRAIN.md](TERRAIN.md)
- [SKY.md](SKY.md) — horizon tint
