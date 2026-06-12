# Camera Commands

Commands for creating and controlling 3D and 2D cameras. The camera determines the viewpoint for all rendered geometry. moonBASIC supports multiple cameras, orbit modes, FPS modes, shake effects, and smooth interpolation.

## Core Concepts

- **Camera3D** — A perspective or orthographic 3D camera with position, target, up vector, and FOV.
- **Camera2D** — A 2D camera with offset, target, rotation, and zoom. Used for scrolling tile maps and side-scrollers.
- **Orbit mode** — Camera rotates around a target entity at a fixed distance. Controlled by mouse and/or keys.
- **FPS mode** — First-person camera controlled by mouse look.
- Every camera is a **handle** that must be freed when no longer needed.

---

## 3D Camera

### `Camera.Create()`

Creates a new 3D camera with default settings: position (0, 10, 10), looking at origin, FOV 45, perspective projection.

**Returns:** `handle`

```basic
cam = Camera.Create()
cam.pos(0, 10, 20)
cam.look(0, 0, 0)
cam.fov(60)
```

---

### `Camera.Free(cameraHandle)`

Frees a camera handle.

- `cameraHandle` (handle) — Camera to free.

---

### `Camera.SetPos(cameraHandle, x, y, z)` / `cam.pos(x, y, z)`

Sets the camera's world position.

- `x`, `y`, `z` (float) — World coordinates.

```basic
Camera.SetPos(cam, 0, 10, 20)
; or
cam.pos(0, 10, 20)
```

---

### `Camera.SetTarget(cameraHandle, x, y, z)` / `Camera.LookAt(cameraHandle, x, y, z)` / `cam.look(x, y, z)`

Sets the point the camera is looking at.

- `x`, `y`, `z` (float) — World coordinates of the target.

```basic
cam.look(0, 0, 0)  ; Look at the origin
```

---

### `Camera.SetFOV(cameraHandle, degrees)` / `cam.fov(degrees)`

Sets the field of view in degrees. Wider FOV gives a fish-eye effect; narrower FOV gives a telephoto/zoom effect.

- `degrees` (float) — FOV in degrees (typically 45–90).

```basic
cam.fov(60)
```

---

### `Camera.SetProjection(cameraHandle, mode)`

Sets the projection mode.

- `mode` (int) — 0 = Perspective (default), 1 = Orthographic.

---

### `Camera.SetRange(cameraHandle, near, far)`

Sets the near and far clipping planes. Objects closer than `near` or farther than `far` are not rendered.

- `near` (float) — Near clip distance.
- `far` (float) — Far clip distance.

```basic
Camera.SetRange(cam, 0.1, 1000)
```

---

### `Camera.Setup(cameraHandle, upX, upY, upZ)`

Sets the camera's up vector. Default is (0, 1, 0).

- `upX`, `upY`, `upZ` (float) — Up direction.

---

### `Camera.GetPos(cameraHandle)` / `Camera.GetRot(cameraHandle)` / `Camera.GetTarget(cameraHandle)`

Returns the camera's current position, rotation, or target. Returns values via the stack.

---

### `Camera.GetMatrix(cameraHandle)`

Returns the camera's view-projection matrix as a Mat4 handle.

**Returns:** `handle`

---

### `Camera.GetViewRay(cameraHandle)`

Returns the camera's forward ray as a Ray handle.

**Returns:** `handle`

---

## 3D Camera Rendering

### `Camera.Begin(cameraHandle)`

Begins a 3D rendering pass with this camera. Equivalent to `Render.Begin3D(cam)`. All 3D draw calls must be between `Begin` and `End`.

- `cameraHandle` (handle) — Camera to render with.

```basic
Camera.Begin(cam)
    Draw.Grid(20, 1.0)
    Entity.DrawAll()
Camera.End(cam)
```

---

### `Camera.End(cameraHandle)`

Ends the 3D rendering pass.

---

## Camera Movement

### `Camera.Move(cameraHandle, forward, right, up)`

Moves the camera relative to its orientation. Positive `forward` moves toward the target; positive `right` strafes right.

- `forward` (float) — Forward/backward movement.
- `right` (float) — Left/right strafe.
- `up` (float) — Up/down movement.

---

### `Camera.Turn(cameraHandle, pitch, yaw, roll)`

Rotates the camera by the given Euler angles (relative).

---

### `Camera.Rotate(cameraHandle, pitch, yaw, roll)`

Sets the camera's absolute rotation.

---

### `Camera.Zoom(cameraHandle, delta)`

Adjusts the camera's zoom level or distance from target.

- `delta` (float) — Zoom change amount.

---

## Orbit Camera

### `Camera.Orbit(cameraHandle, entityHandle, distance)`

Puts the camera into orbit mode around an entity. The camera revolves around the entity at the specified distance.

- `cameraHandle` (handle) — Camera.
- `entityHandle` (handle) — Target entity to orbit.
- `distance` (float) — Orbit radius.

**How it works:** Updates the camera's yaw, pitch, and distance each frame to follow the entity while maintaining orbit constraints. Use `Camera.SetOrbitLimits` and `Camera.SetOrbitSpeed` to fine-tune.

```basic
cam = Camera.Create()
Camera.Orbit(cam, player, 10)
Camera.UseMouseOrbit(cam, TRUE)
Camera.SetOrbitLimits(cam, -80, 80, 3, 50)  ; pitch min/max, dist min/max
```

---

### `Camera.UseMouseOrbit(cameraHandle, enabled)`

Enables mouse-controlled orbit rotation. Mouse movement rotates the camera around the target.

- `enabled` (bool) — `TRUE` to enable.

---

### `Camera.UseOrbitRightMouse(cameraHandle, enabled)`

When `TRUE`, orbit mouse control only activates while the right mouse button is held.

---

### `Camera.SetOrbitKeys(cameraHandle, leftKey, rightKey)`

Sets keyboard keys for orbit rotation.

- `leftKey` (int) — Key to rotate left (e.g., `KEY_Q`).
- `rightKey` (int) — Key to rotate right (e.g., `KEY_E`).

---

### `Camera.SetOrbitLimits(cameraHandle, minPitch, maxPitch, minDist, maxDist)`

Sets pitch angle and distance limits for orbit mode.

- `minPitch`, `maxPitch` (float) — Vertical angle limits in degrees.
- `minDist`, `maxDist` (float) — Minimum and maximum distance from target.

---

### `Camera.SetOrbitSpeed(cameraHandle, mouseSensitivity, wheelSensitivity)`

Sets orbit rotation speed from mouse and mouse wheel.

---

### `Camera.SetOrbitKeySpeed(cameraHandle, radiansPerSecond)`

Sets orbit rotation speed from keyboard.

---

## FPS Camera

### `Camera.SetFPSMode(cameraHandle, enabled)`

Enables first-person camera mode. The camera locks the cursor and uses mouse look for rotation.

- `enabled` (bool) — `TRUE` to enable FPS mode.

---

### `Camera.ClearFPSMode(cameraHandle)`

Disables FPS mode and unlocks the cursor.

---

### `Camera.UpdateFPS(cameraHandle)`

Updates the FPS camera with mouse look input. Call every frame.

---

## Camera Effects

### `Camera.Shake(cameraHandle, intensity, duration)` / `World.Shake(intensity, duration)`

Applies a camera shake effect for impact or explosion feedback.

- `intensity` (float) — Shake magnitude.
- `duration` (float) — Duration in seconds.

```basic
Camera.Shake(cam, 5.0, 0.3)
```

---

### `Camera.Follow(cameraHandle, entityHandle)`

Makes the camera smoothly follow an entity.

- `entityHandle` (handle) — Entity to follow.

---

### `Camera.SmoothExp(cameraHandle, targetX, targetY, targetZ, smoothFactor)`

Moves the camera toward a target position with exponential smoothing. Higher smoothFactor = faster response.

---

### `Camera.LookAtEntity(cameraHandle, entityHandle)` / `Camera.PointAtEntity(cameraHandle, entityHandle)`

Points the camera at an entity's current position.

---

## World-Space Conversion

### `Camera.WorldToScreen(cameraHandle, x, y, z)` / `World.ToScreen(x, y, z)`

Converts a 3D world position to 2D screen coordinates.

**Returns:** Screen X and Y (via stack or Vec2).

---

### `World.ToWorld(screenX, screenY)` / `Camera.Unproject(cameraHandle, screenX, screenY)`

Converts 2D screen coordinates to a 3D world position.

---

### `Camera.Pick(cameraHandle, screenX, screenY)`

Casts a ray from the camera through screen coordinates for entity picking.

**Returns:** `handle` — Ray handle.

---

### `Camera.MouseRay(cameraHandle)`

Returns a ray from the camera through the current mouse position.

**Returns:** `handle`

---

### `World.MouseFloor(y)`

Returns the world position where the mouse ray intersects a horizontal plane at height `y`.

---

### `World.MousePick()`

Performs entity picking under the mouse cursor.

---

### `Camera.IsOnScreen(cameraHandle, x, y, z)`

Returns `TRUE` if a world point is visible on screen.

**Returns:** `bool`

---

## 2D Camera

### `Camera2D.Create()`

Creates a 2D camera with default settings (no offset, no zoom, no rotation).

**Returns:** `handle`

```basic
cam2d = Camera2D.Create()
Camera2D.SetTarget(cam2d, playerX, playerY)
Camera2D.SetOffset(cam2d, 640, 360)
Camera2D.SetZoom(cam2d, 2.0)
```

---

### `Camera2D.Begin(camera2DHandle)` / `Camera2D.End(camera2DHandle)`

Begins and ends a 2D camera rendering pass. All 2D draw calls between these will be transformed by the camera.

```basic
Camera2D.Begin(cam2d)
    Draw.Rectangle(playerX - 16, playerY - 16, 32, 32, 255, 100, 100, 255)
Camera2D.End(cam2d)
```

---

### `Camera2D.SetTarget(handle, x, y)`

Sets the point the camera is centered on.

---

### `Camera2D.SetOffset(handle, x, y)`

Sets the camera offset (typically half the screen size to center the target).

---

### `Camera2D.SetZoom(handle, zoom)`

Sets the zoom level. 1.0 = normal, 2.0 = 2x zoom in.

---

### `Camera2D.SetRotation(handle, angle)`

Sets the camera rotation in degrees.

---

### `Camera2D.GetPos(handle)` / `Camera2D.GetRotation(handle)`

Returns camera position or rotation.

---

### `Camera2D.Follow(handle, entityOrX, entityOrY)`

Makes the 2D camera smoothly follow a position or entity.

---

### `Camera2D.ZoomToMouse(handle, delta)`

Zooms in/out centered on the mouse position.

---

### `Camera2D.WorldToScreen(handle, x, y)` / `Camera2D.ScreenToWorld(handle, x, y)`

Coordinate conversion between world and screen space.

---

### `Camera2D.Free(handle)`

Frees a 2D camera handle.

---

## Culling

Culling commands test whether objects are visible from the camera's perspective, to avoid drawing off-screen geometry.

### `Cull.SphereVisible(x, y, z, radius)`

Returns `TRUE` if a sphere is visible in the camera frustum.

---

### `Cull.AABBVisible(minX, minY, minZ, maxX, maxY, maxZ)`

Returns `TRUE` if an axis-aligned bounding box is visible.

---

### `Cull.PointVisible(x, y, z)`

Returns `TRUE` if a point is visible.

---

### `Cull.InRange(x, y, z, maxDist)`

Returns `TRUE` if a point is within the maximum render distance.

---

### `Cull.Distance(x, y, z)` / `Cull.DistanceSQ(x, y, z)`

Returns the distance (or squared distance) from the camera to a point. Squared distance is faster for comparisons.

---

### `Cull.SetMaxDistance(distance)` / `Cull.GetMaxDistance()`

Sets/gets the maximum culling distance.

---

## Easy Mode Shortcuts

| Shortcut | Maps To |
|----------|---------|
| `CreateCamera()` | `Camera.Create()` |
| `CameraFOV(cam, deg)` | `Camera.SetFOV(cam, deg)` |
| `CameraShake(cam, i, d)` | `Camera.Shake(cam, i, d)` |
| `CameraLookAt(cam, x, y, z)` | `Camera.LookAt(cam, x, y, z)` |
| `CameraRange(cam, n, f)` | `Camera.SetRange(cam, n, f)` |
| `CameraZoom(cam, z)` | `Camera.Zoom(cam, z)` |
| `PositionCamera(cam, x, y, z)` | `Camera.SetPos(cam, x, y, z)` |
| `CAM(...)` | `Camera.Create()` |

---

## Full Example

A 3D scene with an orbit camera that can be rotated with the mouse.

```basic
Window.Open(1280, 720, "Camera Demo")
Window.SetFPS(60)

; Create orbit camera
cam = Camera.Create()
cam.pos(0, 10, 20)
cam.look(0, 0, 0)
cam.fov(60)

; Enable mouse orbit
Camera.UseMouseOrbit(cam, TRUE)
Camera.UseOrbitRightMouse(cam, TRUE)
Camera.SetOrbitLimits(cam, -85, 85, 5, 50)

angle = 0

WHILE NOT Window.ShouldClose()
    dt = Time.Delta()
    angle = angle + 30 * dt

    Render.Clear(30, 30, 50)

    Camera.Begin(cam)
        Draw.Grid(20, 1.0)
        Draw.Cube(0, 1, 0, 2, 2, 2, 200, 50, 50, 255)
        Draw.Sphere(SIN(angle * 0.0174) * 5, 1, COS(angle * 0.0174) * 5, 0.5, 50, 200, 50, 255)
    Camera.End(cam)

    Draw.Text("Right-click + drag to orbit", 10, 10, 18, 255, 255, 255, 255)
    Draw.Text("Scroll to zoom", 10, 32, 18, 200, 200, 200, 255)
    Render.Frame()
WEND

Camera.Free(cam)
Window.Close()
```

---

## See Also

- [RENDER](RENDER.md) — Frame lifecycle and 3D pass
- [ENTITY](ENTITY.md) — Entities to look at and orbit
- [INPUT](INPUT.md) — Mouse input for camera control
