# DBPro — Camera

moonBASIC cameras are **`CAMERA.*`** on a **handle** from **`CAMERA.MAKE`** (see [CAMERA.md](../CAMERA.md)). DBPro integer slot **cam** → store the **handle** your program gets from **`CAMERA.MAKE`**.

| DBPro | moonBASIC | Notes |
|-------|-----------|--------|
| **MAKE CAMERA (cam)** | ✓ **`Camera.Make()`** | Returns handle, not a reserved slot index. |
| **DELETE CAMERA** | ✓ **`Camera.Free()`** | |
| **POSITION CAMERA (x, y, z)** | ✓ **`Camera.SetPos()`** | Sets eye position. |
| **ROTATE CAMERA (pitch, yaw, roll)** | ≈ **`Camera.SetPos()`** + **`Camera.SetTarget()`** | MoonBasic uses eye/target points. |
| **POINT CAMERA (cam, x, y, z)** | ✓ **`Camera.LookAt()`** | Aims at a world point. |
| **MOVE CAMERA (cam, distance)** | ✓ **`Camera.Move()`** | |
| **SET CAMERA RANGE (near, far)** | ✓ **`Camera.SetRange()`** | Near/far clipping planes. |
| **SET CAMERA FOV (angle)** | ✓ **`Camera.SetFOV()`** | |
| **SET CAMERA VIEW (x, y, w, h)** | ≈ **`Render.SetScissor()`** | |
| **SET CAMERA TO OBJECT** | ✓ **`Camera.LookAtEntity()`** | |
| **CAMERA POSITION X/Y/Z** | ✓ **`Camera.GetPos()`** | Returns Vec3 handle. |
