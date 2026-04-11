# Camera

| Designed | moonBASIC | Notes |
|----------|------------|-------|
| **CreateCamera(parent)** | **`Camera.Make()`** | Returns **camera handle**. Parenting via `Entity.Parent()`. |
| **PositionCamera(x, y, z)** | **`Camera.SetPos()`** | Sets eye position. |
| **PointCamera(x, y, z)** | **`Camera.SetTarget()`** | Sets look-at point. |
| **RotateCamera(p, y, r)** | **`Camera.SetRot()`** | |
| **MoveCamera(d)** | **`Camera.Move()`** | |
| **CameraRange(n, f)** | **`Camera.SetRange()`** | Near/Far clipping planes. |
| **CameraZoom(f)** | **`Camera.SetFOV()`** | moonBASIC uses degrees FOV. |
| **CameraProject(x, y, z)** | **`Camera.Project()`** | World to Screen. |
| **CameraPick(sx, sy)** | **`Camera.GetRay()`** | Screen to World ray. |
| **Viewport(x, y, w, h)** | **`Render.SetScissor()`** | |

**Cleanup:** **`CAMERA.FREE(cam)`** when the camera handle is no longer needed.
