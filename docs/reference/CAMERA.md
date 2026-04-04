# Camera Commands

2D and 3D cameras map to Raylib `Camera2D` / `Camera3D`. In source, use the **`Camera`** and **`Camera2D`** namespaces (calls compile to `CAMERA.*` and `CAMERA2D.*`).

---

## 3D camera (`Camera.*`)

### `Camera.Make()`

Creates a default perspective camera. Returns a **handle** (`Camera3D`).

### `Camera.SetPos(camera, x#, y#, z#)` / `Camera.SetPosition(...)`

Alias pair; sets the camera **eye** position in world space.

### `Camera.SetTarget(camera, x#, y#, z#)`

Sets the **look-at** point in world space.

### `Camera.SetFOV(camera, fovy#)`

Vertical field of view in **degrees**.

### `Camera.Begin(camera)`

Starts 3D mode with this camera (sets projection, depth buffer usage, and the active camera for billboards / deferred paths). Pair with `Camera.End`.

### `Camera.End()`

Ends 3D mode (flushes deferred 3D work, then `EndMode3D`).

### `Camera.Move(camera, dx#, dy#, dz#)`

Translates **both** position and target by the delta (pan/strafe/fly without changing orientation).

### `Camera.GetRay(camera, screenX#, screenY#)`

Screen-space to world **ray** for the **current render size**. Returns a **handle** to a **6-float array**: origin `(x,y,z)` then direction `(dx,dy,dz)`. Use with `Draw3D.Ray` or your own picking.

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

Each returns a **handle** to a **2-float array** `[x#, y#]` for the converted point.

---

## Typical layouts

**2D game with scrolling:** `Camera2D.Make` → adjust target/zoom → each frame: `Render.Clear` → `Camera2D.Begin(cam2d)` → `Draw.*` → `Camera2D.End` → `Render.Frame`.

**3D scene:** `Camera.Make` → set position/target → each frame: `Render.Clear` → `Camera.Begin(cam)` → `Draw3D.*` / `Model.Draw` / … → `Camera.End` → (optional 2D/UI pass) → `Render.Frame`.

---

## See also

- [DRAW2D.md](DRAW2D.md), [DRAW3D.md](DRAW3D.md) — what to draw inside each mode.
- [RENDER.md](RENDER.md) — `Render.Clear` / `Render.Frame` and GPU state.
