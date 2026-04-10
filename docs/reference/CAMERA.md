# Camera Commands

2D and 3D cameras map to Raylib `Camera2D` / `Camera3D`. In source, use the **`Camera`** and **`Camera2D`** namespaces (calls compile to `CAMERA.*` and `CAMERA2D.*`).

**Quick map (Create, SetMode, FollowEntity, Project, Unproject, …):** [CAMERA_LIGHT_RENDER.md](CAMERA_LIGHT_RENDER.md).

**Threading:** Raylib windowing and GL calls run on the **main thread** (see [ARCHITECTURE.md](../../ARCHITECTURE.md)); do not invoke `CAMERA.*` / `CAMERA2D.*` from background goroutines.

---

Blitz3D-style **`Camera.Turn`**, **`Rotate`**, **`Orbit`**, **`Zoom`**, **`Follow`**, and entity-based **`Camera.FollowEntity`** are documented in **[BLITZ3D.md](BLITZ3D.md)**.

## 3D camera (`Camera.*`)

### `Camera.Make()` / `Cam()` / `CAM()`

Creates a default perspective camera. Returns a **handle** (`Camera3D`). **`Cam`** and **`CAM`** are aliases of **`Camera.Make`** (short Blitz-style names). Dot-syntax on the handle uses normal call syntax, for example **`cam.Pos(x,y,z)`**, **`cam.FOV(55)`**, **`cam.Orbit(tx,ty,tz,yaw,pitch,dist)`**, **`cam.LookAt(x,y,z)`**, **`cam.Zoom(amount)`** — see [BLITZ3D.md](BLITZ3D.md).

### `Camera.SetPos(camera, x#, y#, z#)` / `Camera.SetPosition(...)`

Alias pair; sets the camera **eye** position in world space.

### `Camera.SetTarget(camera, x#, y#, z#)`

Sets the **look-at** point in world space.

### `Camera.LookAt(camera, x#, y#, z#)`

Alias of **`Camera.SetTarget`** (same arguments and behaviour).

### `Camera.SetProjection(camera, mode#)`

Sets the Raylib projection mode: **`0`** = perspective (**`CameraPerspective`**), **`1`** = orthographic (**`CameraOrthographic`**). In orthographic mode, **`fovy`** is interpreted as the **near-plane height** in world units (Raylib convention).

### `Camera.SetRange(camera, near#, far#)`

Calls **`rl.SetClipPlanes`** before **`Camera.Begin`** for that camera (when **near** &lt; **far** and both positive). This is separate from software **`Cull.*`** distance tests.

### `Camera.SetActive(camera)` / `Camera.GetActive()`

**`SetActive`** records the handle for tooling; **`GetActive`** returns the last handle passed to **`Camera.Begin`** (or **`SetActive`**), or **0** if none.

### `Camera.WorldToScreen2D(camera, wx#, wy#, wz#)`

Alias of **`Camera.WorldToScreen`** — returns a **2-element** float array **\[sx, sy\]**.

### `Camera.SetFPSMode(camera, sensitivity#)` / `Camera.ClearFPSMode(camera)` / `Camera.UpdateFPS(camera)`

**`SetFPSMode`** disables the OS cursor; each frame call **`Camera.UpdateFPS`** to run Raylib **`UpdateCamera`** in **first-person** mode. **`ClearFPSMode`** shows the cursor again.

### `Camera.SetFOV(camera, fovy#)`

Vertical field of view in **degrees**.

### `Camera.Begin(camera)`

Starts 3D mode with this camera (sets projection, depth buffer usage, and the active camera for billboards / deferred paths). Pair with `Camera.End`.

### `Camera.End()`

Ends 3D mode (flushes deferred 3D work, then `EndMode3D`).

### `Camera.Move(camera, dx#, dy#, dz#)`

Translates **both** position and target by the delta (pan/strafe/fly without changing orientation).

### `Camera.SetOrbit(camera, tx#, ty#, tz#, yaw#, pitch#, distance#)`

Third-person **spherical** orbit: places the eye on a shell around the target `(tx,ty,tz)` using **yaw** and **pitch** (radians) and **distance** (world units). Yaw follows the usual XZ orbit (sine on X, cosine on Z); pitch raises/lowers the camera.

**Worked example:** `examples/mario64/main_orbit_simple.mb` builds **`yaw` / `pitch` / `dist`** with **`ORBITYAWDELTA`**, **`ORBITPITCHDELTA`**, **`ORBITDISTDELTA`** (see [GAMEHELPERS.md](GAMEHELPERS.md)), then calls **`Camera.SetOrbit`** each frame after **`MOVESTEPX`/`MOVESTEPZ`** with the same **`camYaw`**.

### `Camera.SmoothExp(current#, target#, smoothHz#, dt#)` → **float**

**Registry:** **`CAMERA.SMOOTHEXP`**. One step of **exponential smoothing** toward a target (same idea as critically damped lag on an angle or scalar):

`result = current + (target - current) * (1 - exp(-smoothHz * dt))`

- **`smoothHz`** — larger values follow the target faster (try **20–35** for third-person orbit yaw/pitch at 60 FPS).
- **`dt`** — use **`Time.Delta()`** (clamped) so smoothing stays framerate-independent.

Typical pattern: mouse and keys update **target** yaw/pitch (`camYawT`, `camPitchT`), clamp pitch, then:

`camYaw = Camera.SmoothExp(camYaw, camYawT, 28.0, dt)`  
`camPitch = Camera.SmoothExp(camPitch, camPitchT, 28.0, dt)`

Pass **`camYaw` / `camPitch`** into **`Camera.OrbitEntity`** (and into movement that must match the camera). See **`examples/mario64/main_entities.mb`**.

### `Camera.OrbitAround(camera, tx#, ty#, tz#, yaw#, distance#, cameraY#)` / `Camera.OrbitAroundDeg(...)`

Simpler **third-person** placement: camera stays at fixed world height **`cameraY`**, orbiting the target on the **XZ** plane at **distance** from `(tx,tz)`, with horizontal angle **`yaw`** in **radians** (`OrbitAround`) or **degrees** (`OrbitAroundDeg`). Sets both position and target (target is `(tx,ty,tz)`).

**Keyboard orbit:** store **`yaw`** in radians, then each frame add **`Input.AxisDeg(negKey, posKey, degreesPerSec, dt)`** (same as **`Input.Axis` × `DEGPERSEC`**). Move the player with **`MOVESTEPX`/`MOVESTEPZ`** or **`MOVEX`/`MOVEZ`** × speed × **`dt`** using the same **`yaw`** so walking matches the camera. See **`examples/mario64/main_v2.mb`** and [INPUT.md](INPUT.md).

### `Camera.GetRay(camera, screenX#, screenY#)`

Screen-space to world **ray** for the **current render size**. Returns a **handle** to a **6-float array**: origin `(x,y,z)` then direction `(dx,dy,dz)`. Use with `Draw3D.Ray`, **`RAY.MAKE`** with the same six components, or **`RAY.HITSPHERE_*`** / other **`RAY.HIT*_*`** queries — see **[RAYCAST.md](RAYCAST.md)**.

### `Camera.GetViewRay(screenX#, screenY#, camera, width, height)`

Like `GetRay`, but uses explicit `width` / `height` (positive integers) for the projection instead of `GetRenderWidth` / `GetRenderHeight`.

### `Camera.GetMatrix(camera)`

Returns a **matrix handle** for the camera (view × projection as provided by Raylib). Free with `Matrix.Free` when you are done (see registry: `MATRIX.FREE` on the camera module).

### `Camera.GetPos(camera)` / `Camera.GetTarget(camera)`

Return **Vec3 handles** for the camera position and target (heap objects; use matrix/vector helpers to read components).

### `Camera.SetUp(camera, ux#, uy#, uz#)`

Sets the camera **up** vector (world space).

### `Camera.Free(camera)`

Frees the camera heap object.

### `Camera.WorldToScreen(camera, wx#, wy#, wz#)`

Projects a **world-space** point through the camera to **screen** coordinates (Raylib **`GetWorldToScreen`**). Returns a **handle** to a **2-float array**: `(screenX, screenY)`. Best used while the same camera is active for rendering (inside your 3D pass).

### `Camera.Project(...)` / `CameraProject(...)`

**Alias** of **`Camera.WorldToScreen`** — same arguments and return value. Use for HUD anchors (health bars above entities).

### `Camera.LookAtEntity(camera, entity#)` / `Camera.PointAtEntity(...)`

**Aliases** — sets the camera **target** to the **world position** of **`entity`** (via **`ENTITY.ENTITYX/Y/Z`** + **`SETTARGET`**). For “look at a point” without an entity, use **`Camera.LookAt(x,y,z)`** / **`SetTarget`**.

### `Camera.Pick` / `Camera.GetRay` / picking entities

**`CAMERA.PICK`** is an alias of **`CAMERA.GETRAY`**: it returns a **screen-space ray** (origin + direction), **not** an entity id. To pick objects, use **`RAY.HIT*_*`** or **`ENTITY.PICK`** / collision queries — see [RAYCAST.md](RAYCAST.md), [ENTITY.md](ENTITY.md).

### `Camera.IsOnScreen(camera, wx#, wy#, wz#)` / `Camera.IsOnScreen(camera, wx#, wy#, wz#, margin#)`

Returns **`TRUE`** if the projected point lies inside the current render rectangle, optionally expanded by **`margin`** pixels on each side.

### `Camera.MouseRay(camera)`

Like **`Camera.GetRay`**, but uses **`GetMousePosition`** for **`screenX` / `screenY`** and the current render size. Returns a **6-float ray** handle (origin + direction).

---

## Culling and visibility (`Cull.*`)

Open-world and large 3D scenes should **not** issue a draw call for every object every frame. **CPU-side culling** decides visibility **before** `Draw3D.*` / `Model.Draw` / `Terrain.Draw` runs, so the GPU never receives geometry that is certainly invisible.

moonBASIC extracts the **view frustum** automatically when you call **`Camera.Begin`**: it combines Raylib’s projection and view matrices (`projection * view`), derives **six clip planes** (left, right, bottom, top, near, far), and stores **pitch**, **vertical FOV**, and **camera position** for distance and horizon helpers. **`Camera.End`** clears that state.

**Conservative rule:** tests may mark something visible when it is not (false positive); they must **never** hide something that is actually visible (no false negatives). Distance and frustum tests are tuned for that.

### Frustum math (short)

Each plane is `ax + by + cz + d = 0` with a **normalised** normal `(a,b,c)`. A point is on the **visible** side when `ax + by + cz + d > 0`.

- **Sphere:** For each plane, compute `distance = a*cx + b*cy + c*cz + d`. If `distance < -r` for **any** plane, the sphere is entirely outside — **cull**. Otherwise it may intersect the frustum.
- **AABB (axis-aligned box):** For each plane, pick the **positive vertex**: the corner of the box that is farthest along the plane normal `(a,b,c)` (per axis, use `min` or `max` depending on the sign of `a`, `b`, `c`). If that vertex is behind the plane (`distance < 0`), the whole box is outside — **cull**.
- **Combined matrix:** Planes come from the **same** `projection * view` as rendering, using the current render aspect ratio at **`Camera.Begin`**. Raylib’s `Matrix` stores rows as `(M0,M4,M8,M12)`, …; **column** vectors are `(M0..M3)`, `(M4..M7)`, `(M8..M11)`, `(M12..M15)`. Gribb–Hartmann combinations `c3 ± c0`, `c3 ± c1`, `c3 ± c2` yield the six planes, then each plane is **normalised** so distances are metric.

### Recommended test order (cheapest first)

1. **`Cull.InRange`** (squared distance, no `sqrt`) — drop far objects.
2. **`Cull.BehindHorizon`** — cheap reject for terrain when the camera is high (optional; terrain drawing uses this internally).
3. **`Cull.SphereVisible` / `Cull.AABBVisible` / `Cull.PointVisible`** — frustum tests.
4. **Occlusion** — Phase B (future); APIs exist now and remain stable.

### When the frustum is “active”

Between **`Camera.Begin`** and **`Camera.End`**, frustum tests use the **current** camera. **`Cull.InRange`**, **`Cull.Distance`**, and **`Cull.DistanceSq`** compare against the **last** camera position captured at **`Camera.Begin`**. If you call **`Cull.SphereVisible`** (or AABB/point) **outside** a Begin/End pair, the implementation returns **`TRUE`** (do not cull) so scripts do not accidentally hide objects; you simply get no frustum benefit until you move the test inside Begin/End.

**`Cull.SetMaxDistance` / `Cull.GetMaxDistance`** set a **global default** draw distance (world units). It is used by:

- **`Cull.InRange(cx, cy, cz)`** (three arguments) — compares to that default.
- **`Cull.SphereVisible`**, **`Cull.AABBVisible`**, **`Cull.PointVisible`** — apply distance **before** frustum, using the same default.

### Per-command reference

| Command | Arguments | Returns | Notes |
|--------|-----------|---------|--------|
| **`Cull.SphereVisible`** | `cx, cy, cz, r` | `bool` | Inside Begin/End: distance (default max) then frustum. Outside Begin/End: always `TRUE`. |
| **`Cull.AABBVisible`** | `minX, minY, minZ, maxX, maxY, maxZ` | `bool` | Uses box centre for distance vs default max, then AABB/frustum test. |
| **`Cull.PointVisible`** | `x, y, z` | `bool` | Cheapest frustum test; still uses default distance first. |
| **`Cull.InRange`** | `cx, cy, cz` **or** `cx, cy, cz, maxdist` | `bool` | Three-arg form uses **`Cull.SetMaxDistance`**. If no active Begin, returns **`TRUE`** (no distance reject). |
| **`Cull.Distance`** | `cx, cy, cz` | `float` | Euclidean distance to last Begin camera; **`0`** if no Begin yet. |
| **`Cull.DistanceSq`** | `cx, cy, cz` | `float` | Squared distance; no `sqrt`. |
| **`Cull.BehindHorizon`** | `camera, maxY, cx, cz` | `bool` | **`TRUE`** if terrain/feature top at `maxY` over `(cx,cz)` is fully below the camera’s bottom-of-view angle (uses camera pitch + FOV). Does **not** require Begin/End. |
| **`Cull.BatchSphere`** | `positions, radii, results` | — | **`positions`**: 1D float array `[x0,y0,z0,x1,y1,z1,…]`. **`radii`**: one float per sphere. **`results`**: bool array (same length as radii). Writes **`TRUE`/`FALSE`** per index. Uses default max distance + frustum when Begin is active. |
| **`Cull.OcclusionEnable`** | `enable?` | `bool` | Phase A: stores flag; returns **`TRUE`**. Phase B: depth pyramid. |
| **`Cull.OccluderAdd`** | `model` | — | Phase A: records handle for future use. |
| **`Cull.OccluderClear`** | — | — | Clears occluder list. |
| **`Cull.IsOccluded`** | `cx, cy, cz, r` | `bool` | Phase A: always **`FALSE`**. |
| **`Cull.SetBackfaceCulling`** | `enable?` | — | Maps to Raylib **`EnableBackfaceCulling`** / **`DisableBackfaceCulling`**. |
| **`Cull.SetMaxDistance`** | `maxdist` | — | Default world radius for distance culling. |
| **`Cull.GetMaxDistance`** | — | `float` | Current default. |
| **`Cull.StatsReset`** | — | — | Zeros all counters — call **once per frame** before your cull tests if you want clean numbers. |
| **`Cull.StatsTotal`** | — | `int` | Tests recorded (sphere/AABB/point/batch iteration / horizon). |
| **`Cull.StatsCulled`** | — | `int` | Count that failed a test (sum of culled outcomes). |
| **`Cull.StatsVisible`** | — | `int` | Count that passed. |
| **`Cull.StatsFrustumCulled`** / **`Cull.StatsDistanceCulled`** / **`Cull.StatsHorizonCulled`** / **`Cull.StatsOcclusionCulled`** | — | `int` | Breakdown (occlusion stays **0** in Phase A). |

### Terrain integration

**`Terrain.Draw`** (inside **`Camera.Begin`**) skips chunks that fail, in order:

1. **`BehindHorizonActive`** — uses pitch/FOV from **`Camera.Begin`**.
2. **`WithinDistanceActive`** — chunk centre vs **`Cull.SetMaxDistance`** default.
3. **`AABBVisibleActive`** — world AABB from chunk footprint and cached **min/max height** from the heightfield (stored when the chunk mesh is rebuilt).

Chunk **min/max Y** are cached in the terrain runtime so height is not rescanned every frame for drawing or culling. **`Cull.StatsReset`** is **not** called from **`Terrain.Draw`**; reset stats in your own loop when you want a clean HUD.

### Phase B — occlusion

Future work: software depth pyramid / conservative occlusion against terrain. **`Cull.OcclusionEnable`**, **`Cull.OccluderAdd`**, **`Cull.OccluderClear`**, **`Cull.IsOccluded`** keep stable names and arguments; Phase A behaviour is documented above.

### Troubleshooting

- **0% culled, stats idle:** **`Cull.StatsReset`** not called, or cull calls run **outside** **`Camera.Begin`** / **`Camera.End`**, or the camera never moves relative to objects.
- **Everything culled:** **`Cull.SetMaxDistance`** too small, wrong bounds, or objects placed outside the valid frustum.
- **Terrain pops or missing chunks:** verify **`World.Update`**, **`Chunk.SetRange`**, and that **`Terrain.Draw`** runs **inside** the same **`Camera.Begin`** used for gameplay.

---

## 2D camera (`Camera2D.*`)

### `Camera2D.Make()`

Creates a `Camera2D` with offset initialized to **half the current screen size** (or `800×450` if dimensions are not ready), target `(0,0)`, zoom `1`, rotation `0`.

### `Camera2D.Begin()` / `Camera2D.Begin(camera)`

- **No arguments:** `BeginMode2D` with an **identity** camera (offset and target `(0,0)`, zoom `1`, rotation `0`).
- **One argument:** `BeginMode2D` with the given **camera handle**.

### `Camera2D.End()`

Ends 2D camera mode (`EndMode2D`).

### `Camera2D.SetTarget(camera, x#, y#)` / `Camera2D.SetOffset(camera, x#, y#)` / `Camera2D.SetZoom(camera, zoom#)` / `Camera2D.SetRotation(camera, angle#)`

Update fields on the stored `Camera2D`. **Rotation follows Raylib:** value is in **degrees**. Zoom must be positive (values `<= 0` are clamped to `0.01`).

### `Camera2D.GetMatrix(camera)`

Returns a **matrix handle** for `GetCameraMatrix2D` applied to that camera.

### `Camera2D.WorldToScreen(camera, worldX#, worldY#)` / `Camera2D.ScreenToWorld(camera, screenX#, screenY#)`

Each returns a **handle** to a **2-float array** `[x, y]` for the converted point.

### `Camera2D.Free(camera)`

Frees the **`Camera2D`** heap object (symmetric with **`Camera.Free`** for 3D cameras).

### `Camera2D.Follow(camera, sprite, speed#, dt#)`

Smoothly moves the camera **target** toward the sprite’s world position (requires **`SPRITE.*`**).

### `Camera2D.ZoomToMouse(camera, delta#)` / `Camera2D.ZoomIn` / `Camera2D.ZoomOut`

**`ZoomToMouse`** adjusts zoom while keeping the world point under the cursor fixed; **`ZoomIn`** / **`ZoomOut`** add or subtract zoom amount.

### `Camera2D.Rotation` / `Camera2D.TargetX` / `Camera2D.TargetY`

Read rotation (degrees) and target **X** / **Y**.

---

## Typical layouts

**2D game with scrolling:** `Camera2D.Make` → adjust target/zoom → each frame: `Render.Clear` → `Camera2D.Begin(cam2d)` → `Draw.*` → `Camera2D.End` → `Render.Frame`.

**3D scene:** `Camera.Make` → set position/target → each frame: `Render.Clear` → `Camera.Begin(cam)` → `Draw3D.*` / `Model.Draw` / … → `Camera.End` → (optional 2D/UI pass) → `Render.Frame`.

---

## See also

- [DRAW2D.md](DRAW2D.md), [DRAW3D.md](DRAW3D.md) — what to draw inside each mode.
- [RENDER.md](RENDER.md) — `Render.Clear` / `Render.Frame` and GPU state.
- **Culling** — see **§ Culling and visibility (`Cull.*`)** above; sample: [`testdata/culling_test.mb`](../../testdata/culling_test.mb).
