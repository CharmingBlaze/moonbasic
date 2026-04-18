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

## Extended Command Reference

### Creation aliases

| Command | Description |
|--------|-------------|
| `CAMERA.MAKE(...)` | Deprecated alias of `CAMERA.CREATE`. |

### Queries

| Command | Description |
|--------|-------------|
| `CAMERA.GETPOS(cam)` | Returns `[x,y,z]` world position. |
| `CAMERA.GETROT(cam)` | Returns `[pitch,yaw,roll]` in radians. |
| `CAMERA.GETTARGET(cam)` | Returns `[x,y,z]` look-at target. |
| `CAMERA.GETUP(cam)` | Returns `[x,y,z]` up vector. |
| `CAMERA.GETFOV(cam)` | Returns field of view in degrees. |
| `CAMERA.FOV(cam, fov)` / `CAMERA.SETFOV(cam, fov)` | Set field of view in degrees. |
| `CAMERA.GETMATRIX(cam)` | Returns the view matrix as a `MAT4` handle. |
| `CAMERA.GETPROJECTION(cam)` | Returns the projection matrix as a `MAT4` handle. |
| `CAMERA.SETPROJECTION(cam, mat4)` | Override the projection matrix. |
| `CAMERA.PROJECTION(cam, fov, aspect, near, far)` | Reconfigure perspective projection. |
| `CAMERA.SETRANGE(cam, near, far)` | Set near/far clip planes. |
| `CAMERA.GETACTIVE()` | Returns the currently active camera handle. |
| `CAMERA.SETACTIVE(cam)` | Set active camera for 3D rendering. |
| `CAMERA.GETYAW(cam)` | Returns current yaw angle in radians. |

### Targeting & follow

| Command | Description |
|--------|-------------|
| `CAMERA.LOOKAT(cam, x, y, z)` | Set look-at target world position. |
| `CAMERA.LOOKATENTITY(cam, entity)` | Continuously track an entity. |
| `CAMERA.POINTATENTITY(cam, entity)` | Alias of `LOOKATENTITY`. |
| `CAMERA.SETTARGETENTITY(cam, entity)` | Set entity target for follow cam. |
| `CAMERA.SETUP(cam, ux, uy, uz)` | Set custom up vector. |
| `CAMERA.SETPOSITION(cam, x, y, z)` | Alias of `CAMERA.SETPOS`. |
| `CAMERA.CAMERAFOLLOW(cam, entity, dist, height, speed)` | Smooth follow with lag. |
| `CAMERA.LERPTO(cam, x, y, z, tx, ty, tz, t)` | Lerp camera position and target simultaneously. |
| `CAMERA.SMOOTHEXP(cam, targetX, targetY, targetZ, factor)` | Exponential smooth move toward target. |
| `CAMERA.SHAKE(cam, intensity, duration)` | Trauma-based camera shake. |

### FPS & orbit

| Command | Description |
|--------|-------------|
| `CAMERA.SETFPSMODE(cam, bool)` | Enable first-person camera mode. |
| `CAMERA.CLEARFPSMODE(cam)` | Disable FPS mode. |
| `CAMERA.UPDATEFPS(cam, dt, sensitivity)` | Update FPS mouselook. |
| `CAMERA.TURNLEFT(cam, speed)` | Yaw left at `speed` degrees/s. |
| `CAMERA.TURNRIGHT(cam, speed)` | Yaw right at `speed` degrees/s. |
| `CAMERA.SETORBIT(cam, pivotX, pivotY, pivotZ, dist)` | Set orbit pivot and distance. |
| `CAMERA.SETORBITSPEED(cam, speed)` | Set orbit drag speed. |
| `CAMERA.SETORBITLIMITS(cam, minPitch, maxPitch, minDist, maxDist)` | Clamp orbit pitch and zoom. |
| `CAMERA.SETORBITKEYS(cam, bool)` | Enable keyboard orbit control. |
| `CAMERA.SETORBITKEYSPEED(cam, speed)` | Keyboard orbit speed. |
| `CAMERA.ORBITCAMERA(cam, dx, dy, dz)` | Orbit around current target by delta. |
| `CAMERA.ORBITAROUND(cam, x, y, z, dist, yaw, pitch)` | Orbit around a world point. |
| `CAMERA.ORBITAROUNDEG(cam, entity, dist, yaw, pitch)` | Orbit around entity with Euler control. |
| `CAMERA.ORBITENTITY(cam, entity, dist)` | Auto-orbit around entity. |
| `CAMERA.USEMOUSEORBIT(cam, btn)` | Enable mouse-drag orbit when `btn` held. |
| `CAMERA.USEORBITRIGHTMOUSE(cam, bool)` | Enable right-mouse-button orbit. |

### Raycasting & picking

| Command | Description |
|--------|-------------|
| `CAMERA.GETRAY(cam, sx, sy)` | Returns ray `[ox,oy,oz, dx,dy,dz]` from screen point. |
| `CAMERA.GETVIEWRAY(cam)` | Returns forward view ray from camera center. |
| `CAMERA.MOUSERAY(cam)` | Returns ray from current mouse position. |
| `CAMERA.RAYCASTMOUSE(cam, maxDist)` | Cast ray from mouse; returns entity id hit. |
| `CAMERA.PICK(cam, sx, sy, maxDist)` | Pick entity at screen coords. |

### Screen-space

| Command | Description |
|--------|-------------|
| `CAMERA.WORLDTOSCREEN(cam, x, y, z)` | Returns `[sx, sy]` screen coords of world point. |
| `CAMERA.WORLDTOSCREEN2D(cam, x, y, z)` | Returns `[sx, sy]` in 2D overlay space. |
| `CAMERA.ISONSCREEN(cam, x, y, z)` | Returns `TRUE` if world point is in the view frustum. |
| `CAMERA.XZBASIS(cam)` | Returns the camera's XZ forward/right basis vectors for WASD mapping. |

---

## See also

- [DRAW2D.md](DRAW2D.md), [DRAW3D.md](DRAW3D.md) — what to draw inside each mode.
- [RENDER.md](RENDER.md) — **`RENDER.CLEAR`** / **`RENDER.FRAME`** and GPU state.
- **Culling** — see **§ Culling and visibility (`CULL.*`)** above.
