# Game — English / Blitz-style helpers

Blitz-like names for **third-person** movement on **XZ** and **orbit yaw** deltas. Full detail: [GAMEHELPERS.md](../../GAMEHELPERS.md) (section *Blitz-style English helpers*).

## 2D mover handle (`PLAYER2D`)

| Designed | Implementation | Memory |
|----------|----------------|--------|
| **Create mover** | **`PLAYER2D.MAKE`** | Heap handle — **`PLAYER2D.FREE`** or **`ERASE ALL`**. |
| **Set position** | **`PLAYER2D.SETPOS`** `(p, x#, z#)` | |
| **MoveEntity2D / MovePlayer** | **`MOVEENTITY2D`**, **`MOVEPLAYER`**, **`PLAYER2D.MOVE`** `(p, camYaw#, f#, s#, speed#, dt#)` | Same math as **`MOVESTEPX`/`MOVESTEPZ`** in place. |
| **ClampEntity2D** | **`CLAMPENTITY2D`**, **`PLAYER2D.CLAMP`** `(p, minX, maxX, minZ, maxZ)` | Stores bounds and clamps **now**. |
| **KeepPlayerInBounds** | **`KEEPPLAYERINBOUNDS`**, **`PLAYER2D.KEEPINBOUNDS`** `(p)` | Re-clamps to **last** **`CLAMPENTITY2D`** bounds. |
| **Read X/Z** | **`PLAYER2D.GETX`**, **`PLAYER2D.GETZ`** | Use with **`py#`** for **`BOXTOPLAND`**, drawing, etc. |

## Camera yaw helpers (radians)

The **camera** argument validates the handle; **yaw** still lives in your **`camYaw#`** variable.

| Designed | Implementation | Returns |
|----------|----------------|---------|
| **TurnCameraLeft** | **`CAMERA.TURNLEFT`** `(cam, amount#)` | **float** — add to **`camYaw#`** (negative **`abs(amount)`**). |
| **TurnCameraRight** | **`CAMERA.TURNRIGHT`** `(cam, amount#)` | **float** — add to **`camYaw#`** (**`+abs(amount)`**). |
| **OrbitCamera** | **`CAMERA.ORBITCAMERA`** `(cam, mouseSens#, keyDegPerSec#, dt#)` | **float** — mouse **X** delta × sens **+** **`Input.Orbit(KEY_Q, KEY_E, …)`**-style Q/E yaw. |

Example:

```basic
camYaw# = camYaw# + CAMERA.ORBITCAMERA(cam, MOUSE_ORBIT_SENS#, 77.0, dt#)
```
