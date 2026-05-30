# Game — English / Blitz-style helpers

Blitz-like names for **third-person** movement on **XZ** and **orbit yaw** deltas. Full detail: [GAMEHELPERS.md](../../GAMEHELPERS.md) (section *Blitz-style English helpers*).

## 2D mover handle (`PLAYER2D`)

| Designed | moonBASIC | Memory / notes |
|----------|------------|----------------|
| **Create mover** | **`PLAYER2D.Create()`** (deprecated **`PLAYER2D.Make()`**) | Heap handle — **`PLAYER2D.Free()`** or **`ERASE ALL`**; registry **`PLAYER2D.CREATE`**. |
| **Set position** | **`PLAYER2D.SetPos(p, x, z)`** | Registry **`PLAYER2D.SETPOS`** (deprecated **`PLAYER2D.SETPOSITION`**). |
| **MoveEntity2D / MovePlayer** | **`PLAYER2D.Move(p, camYaw, f, s, speed, dt)`** | Same math as **`MOVESTEPX`/`MOVESTEPZ`** in place. |
| **ClampEntity2D** | **`PLAYER2D.Clamp(p, minX, maxX, minZ, maxZ)`** | Stores bounds and clamps **now**. |
| **KeepPlayerInBounds** | **`PLAYER2D.KeepInBounds(p)`** | Re-clamps to **last** **`CLAMPENTITY2D`** bounds. |

## Camera yaw helpers (radians)

The **camera** argument validates the handle; **yaw** still lives in your **`camYaw`** variable.

| Designed | moonBASIC | Returns |
|----------|------------|---------|
| **TurnCameraLeft** | **`Camera.TurnLeft(cam, n)`** | **float** — add to `camYaw`. |
| **TurnCameraRight** | **`Camera.TurnRight(cam, n)`** | **float** — add to `camYaw`. |
| **OrbitCamera** | **`Camera.Orbit(cam, s, d, dt)`** | **float** — mouse + keys delta. |

Example:

```basic
camYaw = camYaw + Camera.Orbit(cam, MOUSE_ORBIT_SENS, 77.0, dt)
