# Camera Commands

2D and 3D cameras map to Raylib `Camera2D` / `Camera3D`. In source, use the **`Camera`** and **`Camera2D`** namespaces (calls compile to **`CAMERA.*`** and **`CAMERA2D.*`**). For 3D frames, **`RENDER.BEGIN3D(cam)`** / **`RENDER.END3D()`** are the usual pair (they delegate to **`CAMERA.BEGIN`/`CAMERA.END`**).

**Quick map (Create, SetMode, FollowEntity, Project, Unproject, …):** [CAMERA_LIGHT_RENDER.md](CAMERA_LIGHT_RENDER.md).

**Threading:** Raylib windowing and GL calls run on the **main thread** (see [ARCHITECTURE.md](../../ARCHITECTURE.md)); do not invoke **`CAMERA.*`** / **`CAMERA2D.*`** from background goroutines.

**Page shape:** [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) — large multi-topic reference (**3D**, **culling**, **2D**); headings use registry **`CAMERA.*`**. For a compact single-namespace layout see [WAVE.md](WAVE.md).

## Core Workflow (3D)

**`CAMERA.CREATE()`** → configure **`CAMERA.SETPOS`** / **`SETTARGET`** / orbit helpers → bracket scene draws with **`RENDER.BEGIN3D(cam)`** … **`RENDER.END3D()`** (equivalent to **`CAMERA.BEGIN`/`CAMERA.END`**) → **`CAMERA.FREE`** when done.

---

Blitz3D-style **`Camera.Turn`**, **`Rotate`**, **`Orbit`**, **`Zoom`**, **`Follow`**, and entity-based **`Camera.FollowEntity`** are documented in **[BLITZ3D.md](BLITZ3D.md)**.

## 3D camera (`CAMERA.*`)

### `CAMERA.CREATE()`
Creates a new 3D perspective camera.

- **Returns**: (Handle) The new camera handle.
- **Example**:
    ```basic
    cam = CAMERA.CREATE()
    ```

---

### `CAMERA.SETPOS(handle, x, y, z)` / `SETTARGET`
Sets the camera eye position or look-at point.

- **Arguments**:
    - `handle`: (Handle) The camera to modify.
    - `x, y, z`: (Float) World coordinates.
- **Returns**: (Handle) The camera handle (for chaining).

---

### `CAMERA.MOVE(handle, dx, dy, dz)`
Translates **both** position and target by the delta.

- **Returns**: (Handle) The camera handle (for chaining).

---

### `RENDER.BEGIN3D(handle)` / `RENDER.END3D()`
Starts and ends 3D rendering mode.

- **Arguments**:
    - `handle`: (Handle) The camera to use for the pass.
- **Returns**: (None)

---

### `CAMERA.ORBIT(handle, entity, distance)`
Activates a third-person orbit follow behavior.

- **Arguments**:
    - `handle`: (Handle) The camera to move.
    - `entity`: (Handle) The entity to orbit.
    - `distance`: (Float) Preferred orbit radius.
- **Returns**: (Handle) The camera handle (for chaining).

---

### `CAMERA.YAW(handle)`
Returns the internal orbit yaw in **radians** maintained by **`CAMERA.ORBIT`**.

- **Returns**: (Float)

---

### `CAMERA.PROJECT(handle, wx, wy, wz)`
Projects a **world-space** point through the camera to **screen** coordinates.

- **Returns**: (Handle) A 2-float array handle `[screenX, screenY]`.

---

### `CAMERA.FREE(handle)`
Frees the camera heap object.

---

## Culling and visibility (`CULL.*`)

Open-world and large 3D scenes should **not** issue a draw call for every object every frame. **CPU-side culling** decides visibility **before** rendering commands run.

### `CULL.SPHEREVISIBLE(cx, cy, cz, r)`
Returns `TRUE` if a sphere is within the active camera frustum.

- **Arguments**:
    - `cx, cy, cz`: (Float) Sphere center.
    - `r`: (Float) Sphere radius.
- **Returns**: (Boolean)

---

### `CULL.AABBVISIBLE(minX, minY, minZ, maxX, maxY, maxZ)`
Returns `TRUE` if an axis-aligned box is within the frustum.

- **Returns**: (Boolean)

---

### `CULL.INRANGE(cx, cy, cz [, maxDist])`
Returns `TRUE` if a point is within range of the active camera.

- **Arguments**:
    - `cx, cy, cz`: (Float) World position.
    - `maxDist`: (Float, Optional) Override default max distance.
- **Returns**: (Boolean)

---

### `CULL.SETMAXDISTANCE(dist)`
Sets the default world radius for distance culling.

- **Returns**: (None)

---

## 2D camera (`CAMERA2D.*`)

### `CAMERA2D.CREATE()`
Creates a new `Camera2D` handle.

- **Returns**: (Handle) The new camera.

---

### `CAMERA2D.BEGIN([camera])` / `CAMERA2D.END()`
Starts / ends 2D rendering mode.

---

### `CAMERA2D.SETTARGET` / `SETOFFSET` / `SETZOOM` / `SETROTATION`
Updates the 2D camera fields.

- **Returns**: (Handle) The camera handle (for chaining).

---

## Full Example

```basic
WINDOW.OPEN(1280, 720, "Camera Demo")
WINDOW.SETFPS(60)

cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 5, -10)
CAMERA.SETTARGET(cam, 0, 0, 0)

WHILE NOT WINDOW.SHOULDCLOSE()
    RENDER.CLEAR(40, 40, 60)
    RENDER.BEGIN3D(cam)
        DRAW3D.GRID(10, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND

CAMERA.FREE(cam)
WINDOW.CLOSE()
```

---

## See also

- [DRAW2D.md](DRAW2D.md), [DRAW3D.md](DRAW3D.md) — what to draw inside each mode.
- [RENDER.md](RENDER.md) — **`RENDER.CLEAR`** / **`RENDER.FRAME`** and GPU state.
- **Culling** — see **§ Culling and visibility (`CULL.*`)** above.
