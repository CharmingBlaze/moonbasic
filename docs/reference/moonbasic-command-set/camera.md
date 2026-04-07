# Camera

| Designed | Implementation | Memory / notes |
|----------|----------------|----------------|
| **CreateCamera (parent)** | **`CAMERA.MAKE`** — parenting via **`ENTITY`** / your own transform graph | Returns **VM heap handle** — **`CAMERA.FREE`**. |
| **PositionCamera** | **`CAMERA.SETPOS`**, **`SETPOSITION`** | |
| **RotateCamera** | **`CAMERA.ROTATE`**, **`TURN`**, **`SETORBIT`** | Radians vs degrees: see [CAMERA.md](../CAMERA.md). |
| **TurnCameraLeft / TurnCameraRight** (yaw delta) | **`CAMERA.TURNLEFT`**, **`CAMERA.TURNRIGHT`** `(cam, amount#)` → **float** radians | Add return value to your **`camYaw#`**; validates **`cam`** handle. See [game-english.md](game-english.md). |
| **OrbitCamera** (mouse + Q/E yaw delta) | **`CAMERA.ORBITCAMERA`** `(cam, mouseSens#, keyDegPerSec#, dt#)` → **float** | Same idea as mouse delta × sens **`+`** **`Input.Orbit`** for **Q/E**. |
| **MoveCamera** | **`CAMERA.MOVE`** | |
| **CameraRange** | **`CAMERA.SETRANGE`** | Near/far. |
| **CameraZoom** | **`CAMERA.ZOOM`** | Adjusts FOV. |
| **CameraViewport** | **`RENDER.SETSCISSOR`**, **`RENDERTARGET.*`** | |
| **CameraProject** | **`CAMERA.WORLDTOSCREEN`**, **`WORLDTOSCREEN*`** | |
| **CameraPick** | **`CAMERA.GETRAY`**, **`CAMERA.PICK`** alias | |

**Cleanup:** **`CAMERA.FREE(cam)`** when the camera handle is no longer needed.
