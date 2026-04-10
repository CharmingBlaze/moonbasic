# DBPro — Camera

moonBASIC cameras are **`CAMERA.*`** on a **handle** from **`CAMERA.MAKE`** (see [CAMERA.md](../CAMERA.md)). DBPro integer slot **cam** → store the **handle** your program gets from **`CAMERA.MAKE`**.

| DBPro | moonBASIC | Notes |
|-------|-----------|--------|
| **MAKE CAMERA (cam)** | ✓ **`CAMERA.MAKE`** | Returns handle, not a reserved slot index. |
| **DELETE CAMERA** | ✓ **`CAMERA.FREE`** | |
| **POSITION CAMERA** | ✓ **`CAMERA.SETPOS`** / **`SETPOSITION`** | |
| **ROTATE CAMERA** | ✓ **`CAMERA.ROTATE`**, **`TURN`**, **`SETORBIT`** | Many orbit/follow helpers — [CAMERA.md](../CAMERA.md). |
| **MOVE CAMERA (cam, distance)** | ≈ **`CAMERA.MOVE`** | |
| **POINT CAMERA (cam, x, y, z)** | ✓ **`CAMERA.LOOKAT`**, **`SETTARGET`** | |
| **SET CAMERA RANGE (near, far)** | ✓ **`CAMERA.SETRANGE`** | |
| **SET CAMERA FOV** | ✓ **`CAMERA.SETFOV`** | |
| **SET CAMERA ASPECT** | ≈ projection / window aspect | Often implicit from **`WINDOW` / `RENDER`**. |
| **SET CAMERA VIEW (x, y, w, h)** | ≈ **`RENDER.SETSCISSOR`**, **`RENDERTARGET.*`** | |
| **SET CAMERA TO OBJECT** | ✓ **`CAMERA.SETTARGETENTITY`**, **`FOLLOWENTITY`**, **`ORBITENTITY`** | |
| **CAMERA POSITION X/Y/Z** | ✓ **`CAMERA.GETPOS`**, split getters | See manifest. |
