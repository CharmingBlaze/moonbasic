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

### `CreateCamera()` / **`CAMERA.CREATE()`** 
Easy Mode **`CreateCamera()`** forwards to **`CAMERA.CREATE`** (same defaults as below). Prefer this or the registry call to avoid deprecated **`CAMERA.MAKE`**.

---

### `CAMERA.MAKE()` (deprecated) 
Same as **`CAMERA.CREATE`**; the compiler warns — prefer **`CreateCamera()`** or **`CAMERA.CREATE()`**. Creates a default perspective camera. Returns a **handle** (`Camera3D`).

---

### `Cam()` / `CAM()` 
Aliases of **`Camera.Make()`** / **`CAMERA.CREATE`**. Short Blitz-style names.

---

### `CAMERA.SETPOS(handle, x, y, z)` · Easy Mode `Camera.SetPos` 
Sets the camera **eye** position in world space.

---

### `CAMERA.SETTARGET(handle, x, y, z)` · Easy Mode `Camera.SetTarget` 
Sets the **look-at** point in world space.

---

### `CAMERA.MOVE(handle, dx, dy, dz)` · Easy Mode `Camera.Move` 
Translates **both** position and target by the delta (pan/strafe/fly without changing orientation).

---

### `RENDER.BEGIN3D(handle)` / `RENDER.END3D()` · `CAMERA.BEGIN(handle)` / `CAMERA.END()` 
Starts / ends 3D mode with this camera. Prefer **`RENDER.BEGIN3D`/`END3D`** in new scripts; **`CAMERA.BEGIN`/`END`** are equivalent. **`CAMERA.END`** takes no arguments.

---

### `CAMERA.ORBIT(handle, entity, distance)` · Easy Mode `Camera.Orbit` 
Third-person orbit-follow around an entity. Each frame, updates internal yaw, pitch, and distance based on mouse/keyboard input.

---

### `CAMERA.YAW(handle)` · Easy Mode `Camera.Yaw` 
Returns the internal orbit yaw in **radians** maintained by **`CAMERA.ORBIT`**. Use this to set player rotation so movement matches the camera.

---

### `CAMERA.PROJECT(handle, wx, wy, wz)` · Easy Mode `Camera.Project` 
Projects a **world-space** point through the camera to **screen** coordinates. Returns a **handle** to a **2-float array**: `(screenX, screenY)`.

---

### `CAMERA.FREE(handle)` · Easy Mode `Camera.Free` 
Frees the camera heap object.

**Orbit defaults:** mouse orbit **on**; **RMB required** to apply mouse delta; keys **Q** / **E**; mouse **0.005**, wheel **1.0**, keyboard **1.5** rad/s; pitch **−1.5…1.5** rad; distance **2…50**; look-at offset **0.5** above entity base.

**Worked examples and recipes:** `examples/mario64/README.md` (orbit configuration section).

#### `CAMERA.USEMOUSEORBIT(camera, useMouse)` · Easy Mode `Camera.UseMouseOrbit`

- **`TRUE`:** mouse (subject to **`UseOrbitRightMouse`**) adjusts yaw/pitch.
- **`FALSE`:** mouse does **not** move the orbit — use when the cursor aims a weapon or UI; combine with **`SetOrbitKeys`** / wheel.

#### `CAMERA.USEORBITRIGHTMOUSE(camera, requireRightMouse)` · Easy Mode `Camera.UseOrbitRightMouse`

- **`TRUE` (default):** mouse delta applies **only while right mouse button is held** (common third-person feel).
- **`FALSE`:** mouse delta applies **without** holding RMB (inspector-style orbit).

#### `CAMERA.SETORBITKEYS(camera, leftKey, rightKey)` · Easy Mode `Camera.SetOrbitKeys`

Raylib keyboard constants (e.g. **`KEY_Q`**, **`KEY_E`**, **`KEY_LEFT`**). **`0`** disables that side; **`(0, 0)`** disables keyboard orbit entirely.

#### `CAMERA.SETORBITLIMITS(camera, minPitch, maxPitch, minDist, maxDist)` · Easy Mode `Camera.SetOrbitLimits`

**Pitch** in **radians** (vertical tilt of the orbit). **Distance** in **world units** (zoom in/out along the orbit radius). Prevents camera flipping and runaway zoom.

#### `CAMERA.SETORBITSPEED(camera, mouseSens, wheelSens)` · Easy Mode `Camera.SetOrbitSpeed`

Scales **mouse drag** (yaw and pitch) and **mouse wheel** zoom. Larger values feel faster.

#### `CAMERA.SETORBITKEYSPEED(camera, keyRadPerSec)` · Easy Mode `Camera.SetOrbitKeySpeed`

Keyboard orbit yaw rate in **radians per second** (framerate-independent).

---

### `CAMERA.SMOOTHEXP(current, target, smoothHz, dt)` → **float** · Easy Mode `Camera.SmoothExp` 

**Registry:** **`CAMERA.SMOOTHEXP`**. One step of **exponential smoothing** toward a target (same idea as critically damped lag on an angle or scalar):

`result = current + (target - current) * (1 - exp(-smoothHz * dt))`

- **`smoothHz`** — larger values follow the target faster (try **20–35** for third-person orbit yaw/pitch at 60 FPS).
- **`dt`** — use **`TIME.DELTA()`** (clamped) so smoothing stays framerate-independent.

Typical pattern: mouse and keys update **target** yaw/pitch (`camYawT`, `camPitchT`), clamp pitch, then:

`camYaw = CAMERA.SMOOTHEXP(camYaw, camYawT, 28.0, dt)`  
`camPitch = CAMERA.SMOOTHEXP(camPitch, camPitchT, 28.0, dt)`

Pass **`camYaw` / `camPitch`** into **`CAMERA.ORBITENTITY`** (and into movement that must match the camera). See **`examples/mario64/main_entities.mb`**.

---

### `CAMERA.ORBITAROUND(...)` / `CAMERA.ORBITAROUNDEG(...)` · Easy Mode `Camera.OrbitAround` / `OrbitAroundDeg` 

Simpler **third-person** placement: camera stays at fixed world height **`cameraY`**, orbiting the target on the **XZ** plane at **distance** from `(tx,tz)`, with horizontal angle **`yaw`** in **radians** (`OrbitAround`) or **degrees** (`OrbitAroundDeg`). Sets both position and target (target is `(tx,ty,tz)`).

**Keyboard orbit:** store **`yaw`** in radians, then each frame add **`INPUT.AXISDEG(negKey, posKey, degreesPerSec, dt)`** (same as **`INPUT.AXIS` × `DEGPERSEC`**). Move the player with **`MOVESTEPX`/`MOVESTEPZ`** or **`MOVEX`/`MOVEZ`** × speed × **`dt`** using the same **`yaw`** so walking matches the camera. See **`examples/mario64/main_v2.mb`** and [INPUT.md](INPUT.md).

---

### `CAMERA.GETRAY(camera, screenX, screenY)` · Easy Mode `Camera.GetRay` 

Screen-space to world **ray** for the **current render size**. Returns a **handle** to a **6-float array**: origin `(x,y,z)` then direction `(dx,dy,dz)`. Use with **`DRAW3D.RAY`**, **`RAY.CREATE`** (canonical) or deprecated **`RAY.MAKE`** with the same six components, or **`RAY.HITSPHERE_*`** / other **`RAY.HIT*_*`** queries — see **[RAYCAST.md](RAYCAST.md)**.

---

### `CAMERA.GETVIEWRAY(screenX, screenY, camera, width, height)` · Easy Mode `Camera.GetViewRay` 

Like **`CAMERA.GETRAY`**, but uses explicit `width` / `height` (positive integers) for the projection instead of `GetRenderWidth` / `GetRenderHeight`.

---

### `CAMERA.GETMATRIX(camera)` · Easy Mode `Camera.GetMatrix` 

Returns a **matrix handle** for the camera (view × projection as provided by Raylib). Free with **`MATRIX.FREE`** when you are done.

---

### `CAMERA.GETPOS(camera)` / `CAMERA.GETTARGET(camera)` · Easy Mode `Camera.GetPos` / `GetTarget` 

Return **Vec3 handles** for the camera position and target (heap objects; use matrix/vector helpers to read components).

---

### `CAMERA.SETUP(camera, ux, uy, uz)` · Easy Mode `Camera.SetUp` 

Sets the camera **up** vector (world space).

---

### `CAMERA.WORLDTOSCREEN(camera, wx, wy, wz)` · Easy Mode `Camera.WorldToScreen` 

Projects a **world-space** point through the camera to **screen** coordinates (Raylib **`GetWorldToScreen`**). Returns a **handle** to a **2-float array**: `(screenX, screenY)`. Best used while the same camera is active for rendering (inside your 3D pass).

---

### `CAMERA.PROJECT(...)` / `CameraProject(...)` · Easy Mode 

**Alias** of **`CAMERA.WORLDTOSCREEN`** — same arguments and return value. Use for HUD anchors (health bars above entities).

---

### `CAMERA.LOOKATENTITY` / `CAMERA.POINTATENTITY` · Easy Mode `Camera.LookAtEntity` / `PointAtEntity` 

**Aliases** — sets the camera **target** to the **world position** of **`entity`** (via **`ENTITY.ENTITYX/Y/Z`** + **`SETTARGET`**). For “look at a point” without an entity, use **`CAMERA.LOOKAT(x,y,z)`** / **`CAMERA.SETTARGET`**.

---

### `CAMERA.PICK` / `CAMERA.GETRAY` / picking entities 

**`CAMERA.PICK`** is an alias of **`CAMERA.GETRAY`**: it returns a **screen-space ray** (origin + direction), **not** an entity id. To pick objects, use **`RAY.HIT*_*`** or **`ENTITY.PICK`** / collision queries — see [RAYCAST.md](RAYCAST.md), [ENTITY.md](ENTITY.md).

---

### `CAMERA.ISONSCREEN(...)` · Easy Mode `Camera.IsOnScreen` 

Returns **`TRUE`** if the projected point lies inside the current render rectangle, optionally expanded by **`margin`** pixels on each side.

---

### `CAMERA.MOUSERAY(camera)` · Easy Mode `Camera.MouseRay` 

Like **`CAMERA.GETRAY`**, but uses **`GetMousePosition`** for **`screenX` / `screenY`** and the current render size. Returns a **6-float ray** handle (origin + direction).

---

## Culling and visibility (`CULL.*`)

Open-world and large 3D scenes should **not** issue a draw call for every object every frame. **CPU-side culling** decides visibility **before** **`DRAW3D.*`** / **`MODEL.DRAW`** / **`TERRAIN.DRAW`** runs, so the GPU never receives geometry that is certainly invisible.

moonBASIC extracts the **view frustum** automatically when you call **`CAMERA.BEGIN`** (or **`RENDER.BEGIN3D`**): it combines Raylib’s projection and view matrices (`projection * view`), derives **six clip planes** (left, right, bottom, top, near, far), and stores **pitch**, **vertical FOV**, and **camera position** for distance and horizon helpers. **`CAMERA.END`** (or **`RENDER.END3D`**) clears that state.

**Conservative rule:** tests may mark something visible when it is not (false positive); they must **never** hide something that is actually visible (no false negatives). Distance and frustum tests are tuned for that.

### Frustum math (short) 

Each plane is `ax + by + cz + d = 0` with a **normalised** normal `(a,b,c)`. A point is on the **visible** side when `ax + by + cz + d > 0`.

- **Sphere:** For each plane, compute `distance = a*cx + b*cy + c*cz + d`. If `distance < -r` for **any** plane, the sphere is entirely outside — **cull**. Otherwise it may intersect the frustum.
- **AABB (axis-aligned box):** For each plane, pick the **positive vertex**: the corner of the box that is farthest along the plane normal `(a,b,c)` (per axis, use `min` or `max` depending on the sign of `a`, `b`, `c`). If that vertex is behind the plane (`distance < 0`), the whole box is outside — **cull**.
- **Combined matrix:** Planes come from the **same** `projection * view` as rendering, using the current render aspect ratio at **`CAMERA.BEGIN`** (or **`RENDER.BEGIN3D`**). Raylib’s `Matrix` stores rows as `(M0,M4,M8,M12)`, …; **column** vectors are `(M0..M3)`, `(M4..M7)`, `(M8..M11)`, `(M12..M15)`. Gribb–Hartmann combinations `c3 ± c0`, `c3 ± c1`, `c3 ± c2` yield the six planes, then each plane is **normalised** so distances are metric.

---

### Recommended test order (cheapest first) 

1. **`CULL.INRANGE`** (squared distance, no `sqrt`) — drop far objects.
2. **`CULL.BEHINDHORIZON`** — cheap reject for terrain when the camera is high (optional; terrain drawing uses this internally).
3. **`CULL.SPHEREVISIBLE`** / **`CULL.AABBVISIBLE`** / **`CULL.POINTVISIBLE`** — frustum tests.
4. **Occlusion** — Phase B (future); APIs exist now and remain stable.

---

### When the frustum is “active” 

Between **`CAMERA.BEGIN`** and **`CAMERA.END`** (or **`RENDER.BEGIN3D`** / **`RENDER.END3D`**), frustum tests use the **current** camera. **`CULL.INRANGE`**, **`CULL.DISTANCE`**, and **`CULL.DISTANCESQ`** compare against the **last** camera position captured at **`CAMERA.BEGIN`**. If you call **`CULL.SPHEREVISIBLE`** (or AABB/point) **outside** a Begin/End pair, the implementation returns **`TRUE`** (do not cull) so scripts do not accidentally hide objects; you simply get no frustum benefit until you move the test inside Begin/End.

**`CULL.SETMAXDISTANCE`** / **`CULL.GETMAXDISTANCE`** set a **global default** draw distance (world units). It is used by:

- **`CULL.INRANGE(cx, cy, cz)`** (three arguments) — compares to that default.
- **`CULL.SPHEREVISIBLE`**, **`CULL.AABBVISIBLE`**, **`CULL.POINTVISIBLE`** — apply distance **before** frustum, using the same default.

---

### Per-command reference 

| Command | Arguments | Returns | Notes |
|--------|-----------|---------|--------|
| **`CULL.SPHEREVISIBLE`** | `cx, cy, cz, r` | `bool` | Inside Begin/End: distance (default max) then frustum. Outside Begin/End: always `TRUE`. |
| **`CULL.AABBVISIBLE`** | `minX, minY, minZ, maxX, maxY, maxZ` | `bool` | Uses box centre for distance vs default max, then AABB/frustum test. |
| **`CULL.POINTVISIBLE`** | `x, y, z` | `bool` | Cheapest frustum test; still uses default distance first. |
| **`CULL.INRANGE`** | `cx, cy, cz` **or** `cx, cy, cz, maxdist` | `bool` | Three-arg form uses **`CULL.SETMAXDISTANCE`**. If no active Begin, returns **`TRUE`** (no distance reject). |
| **`CULL.DISTANCE`** | `cx, cy, cz` | `float` | Euclidean distance to last Begin camera; **`0`** if no Begin yet. |
| **`CULL.DISTANCESQ`** | `cx, cy, cz` | `float` | Squared distance; no `sqrt`. |
| **`CULL.BEHINDHORIZON`** | `camera, maxY, cx, cz` | `bool` | **`TRUE`** if terrain/feature top at `maxY` over `(cx,cz)` is fully below the camera’s bottom-of-view angle (uses camera pitch + FOV). Does **not** require Begin/End. |
| **`CULL.BATCHSPHERE`** | `positions, radii, results` | — | **`positions`**: 1D float array `[x0,y0,z0,x1,y1,z1,…]`. **`radii`**: one float per sphere. **`results`**: bool array (same length as radii). Writes **`TRUE`/`FALSE`** per index. Uses default max distance + frustum when Begin is active. |
| **`CULL.OCCLUSIONENABLE`** | `enable` | `bool` | Phase A: stores flag; returns **`TRUE`**. Phase B: depth pyramid. |
| **`CULL.OCCLUDERADD`** | `model` | — | Phase A: records handle for future use. |
| **`CULL.OCCLUDERCLEAR`** | — | — | Clears occluder list. |
| **`CULL.ISOCCLUDED`** | `cx, cy, cz, r` | `bool` | Phase A: always **`FALSE`**. |
| **`CULL.SETBACKFACECULLING`** | `enable` | — | Maps to Raylib **`EnableBackfaceCulling`** / **`DisableBackfaceCulling`**. |
| **`CULL.SETMAXDISTANCE`** | `maxdist` | — | Default world radius for distance culling. |
| **`CULL.GETMAXDISTANCE`** | — | `float` | Current default. |
| **`CULL.STATSRESET`** | — | — | Zeros all counters — call **once per frame** before your cull tests if you want clean numbers. |
| **`CULL.STATSTOTAL`** | — | `int` | Tests recorded (sphere/AABB/point/batch iteration / horizon). |
| **`CULL.STATSCULLED`** | — | `int` | Count that failed a test (sum of culled outcomes). |
| **`CULL.STATSVISIBLE`** | — | `int` | Count that passed. |
| **`CULL.STATSFRUSTUMCULLED`** / **`CULL.STATSDISTANCECULLED`** / **`CULL.STATSHORIZONCULLED`** / **`CULL.STATSOCCLUSIONCULLED`** | — | `int` | Breakdown (occlusion stays **0** in Phase A). |

---

### Terrain integration 

**`TERRAIN.DRAW`** (inside **`CAMERA.BEGIN`**) skips chunks that fail, in order:

1. **`BehindHorizonActive`** — uses pitch/FOV from **`CAMERA.BEGIN`**.
2. **`WithinDistanceActive`** — chunk centre vs **`CULL.SETMAXDISTANCE`** default.
3. **`AABBVisibleActive`** — world AABB from chunk footprint and cached **min/max height** from the heightfield (stored when the chunk mesh is rebuilt).

Chunk **min/max Y** are cached in the terrain runtime so height is not rescanned every frame for drawing or culling. **`CULL.STATSRESET`** is **not** called from **`TERRAIN.DRAW`**; reset stats in your own loop when you want a clean HUD.

---

### Phase B — occlusion 

Future work: software depth pyramid / conservative occlusion against terrain. **`CULL.OCCLUSIONENABLE`**, **`CULL.OCCLUDERADD`**, **`CULL.OCCLUDERCLEAR`**, **`CULL.ISOCCLUDED`** keep stable names and arguments; Phase A behaviour is documented above.

---

### Troubleshooting 

- **0% culled, stats idle:** **`CULL.STATSRESET`** not called, or cull calls run **outside** **`CAMERA.BEGIN`** / **`CAMERA.END`**, or the camera never moves relative to objects.
- **Everything culled:** **`CULL.SETMAXDISTANCE`** too small, wrong bounds, or objects placed outside the valid frustum.
- **Terrain pops or missing chunks:** verify **`WORLD.UPDATE`**, **`CHUNK.SETRANGE`**, and that **`TERRAIN.DRAW`** runs **inside** the same **`CAMERA.BEGIN`** used for gameplay.

---

## 2D camera (`CAMERA2D.*`)

### `CAMERA2D.CREATE()` (canonical; deprecated `CAMERA2D.MAKE`) · Easy Mode `Camera2D.Make()` 

Creates a `Camera2D` with offset initialized to **half the current screen size** (or `800×450` if dimensions are not ready), target `(0,0)`, zoom `1`, rotation `0`.

---

### `CAMERA2D.BEGIN()` / `CAMERA2D.BEGIN(camera)` · Easy Mode `Camera2D.Begin` 

- **No arguments:** `BeginMode2D` with an **identity** camera (offset and target `(0,0)`, zoom `1`, rotation `0`).
- **One argument:** `BeginMode2D` with the given **camera handle**.

---

### `CAMERA2D.END()` · Easy Mode `Camera2D.End` 

Ends 2D camera mode (`EndMode2D`).

---

### `CAMERA2D.SETTARGET` / `CAMERA2D.SETOFFSET` / `CAMERA2D.SETZOOM` / `CAMERA2D.SETROTATION` · Easy Mode `Camera2D.Set*` 

Update fields on the stored `Camera2D`. **Rotation follows Raylib:** value is in **degrees**. Zoom must be positive (values `<= 0` are clamped to `0.01`).

---

### `CAMERA2D.GETMATRIX(camera)` · Easy Mode `Camera2D.GetMatrix` 

Returns a **matrix handle** for `GetCameraMatrix2D` applied to that camera.

---

### `CAMERA2D.WORLDTOSCREEN` / `CAMERA2D.SCREENTOWORLD` · Easy Mode `Camera2D.WorldToScreen` / `ScreenToWorld` 

Each returns a **handle** to a **2-float array** `[x, y]` for the converted point.

---

### `CAMERA2D.FREE(camera)` · Easy Mode `Camera2D.Free` 

Frees the **`Camera2D`** heap object (symmetric with **`CAMERA.FREE`** for 3D cameras).

---

### `CAMERA2D.FOLLOW(camera, sprite, speed, dt)` · Easy Mode `Camera2D.Follow` 

Smoothly moves the camera **target** toward the sprite’s world position (requires **`SPRITE.*`**).

---

### `CAMERA2D.ZOOMTOMOUSE` / `CAMERA2D.ZOOMIN` / `CAMERA2D.ZOOMOUT` · Easy Mode `Camera2D.ZoomToMouse` / … 

**`CAMERA2D.ZOOMTOMOUSE`** adjusts zoom while keeping the world point under the cursor fixed; **`ZOOMIN`** / **`ZOOMOUT`** add or subtract zoom amount.

---

### `CAMERA2D.ROTATION` / `CAMERA2D.TARGETX` / `CAMERA2D.TARGETY` · Easy Mode `Camera2D.Rotation` / … 

Read rotation (degrees) and target **X** / **Y**.

---

## Typical layouts

**2D game with scrolling:** **`CAMERA2D.CREATE`** → adjust target/zoom → each frame: **`RENDER.CLEAR`** → **`CAMERA2D.BEGIN(cam2d)`** → **`DRAW.*`** → **`CAMERA2D.END`** → **`RENDER.FRAME`**.

**3D scene:** **`CAMERA.CREATE`** → set position/target → each frame: **`RENDER.CLEAR`** → **`RENDER.BEGIN3D(cam)`** → **`DRAW3D.*`** / **`MODEL.DRAW`** / … → **`RENDER.END3D()`** → (optional 2D/UI pass) → **`RENDER.FRAME`**.

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
- **Culling** — see **§ Culling and visibility (`CULL.*`)** above; sample: [`testdata/culling_test.mb`](../../testdata/culling_test.mb).
