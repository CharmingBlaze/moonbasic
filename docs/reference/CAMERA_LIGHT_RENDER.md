# Camera, Light, and Rendering

A high-level overview of atmospheric controls, camera management, and light setup in the MoonBASIC engine.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. **Camera:** Create a viewpoint with `CAMERA.CREATE` and link it to the scene.
2. **Lights:** Populate the world with `LIGHT.CREATEPOINT` or `LIGHT.CREATEDIRECTIONAL`.
3. **Atmosphere:** Configure global effects like `RENDER.SETFOG` and `RENDER.SETAMBIENT`.
4. **Post-Process:** Apply final polish using `RENDER.SETBLOOM` or `CAMERA.SHAKE`.

---

## 1. Camera Module

| Command | Arguments | Returns | Notes |
|---------|-----------|---------|-------|
| `CAMERA.CREATE()` | None | Handle | Creates a new 3D camera. |
| `CAMERA.SETPOS(c, x, y, z)`| Handle, Float...| Handle | Positions the camera. |
| `CAMERA.SETTARGET(c, x, y, z)`| Handle, Float...| None | Sets the look-at target. |
| `CAMERA.SETFOV(c, fov)` | Handle, Float | None | Sets field of view in degrees. |
| `CAMERA.FOLLOWENTITY(c, e, d, h, s)`| Handle, Int, Float...| None | Smooth third-person follow. |
| `CAMERA.SHAKE(c, a, d)` | Handle, Float...| None | Screen shake effect. |

---

## 2. Light Module

| Command | Arguments | Returns | Notes |
|---------|-----------|---------|-------|
| `LIGHT.CREATEPOINT(x, y, z, r, g, b, e)`| Float... | Handle | Omnidirectional point light. |
| `LIGHT.CREATEDIRECTIONAL(dx, dy, dz, r, g, b, e)`| Float... | Handle | Sun-like directional light. |
| `LIGHT.CREATESPOT(x, y, z, tx, ty, tz, r, g, b, cone, e)`| Float... | Handle | Focused spotlight. |
| `LIGHT.SETCOLOR(lt, r, g, b)`| Handle, Float...| None | Sets light tint (0.0-1.0 or 0-255). |

---

## 3. Atmosphere and Post-Processing

| Command | Arguments | Returns | Notes |
|---------|-----------|---------|-------|
| `RENDER.SETAMBIENT(r, g, b, i)`| Float... | None | Sets base scene lighting. |
| `RENDER.SETFOG(r, g, b, start, end, density)`| Float... | None | Configures distance fog. |
| `RENDER.SETBLOOM(threshold)`| Float | None | Enables post-process bloom. |
| `RENDER.SCREENSHOT(path)` | String | None | Saves frame to disk. |

---

## Full Example

A complete scene setup with a camera follow, directional sun, and atmospheric fog.

```basic
WINDOW.OPEN(1280, 720, "Render Demo")
cam = CAMERA.CREATE()
sun = LIGHT.CREATEDIRECTIONAL(0, -1, 0, 255, 255, 200, 2.0)
player = ENTITY.CREATECUBE(1.0).SETPOS(0, 0.5, 0)

RENDER.SETFOG(20, 20, 30, 10.0, 100.0, 0.01)
RENDER.SETAMBIENT(50, 50, 60, 1.0)

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    
    ; 1. Logic
    CAMERA.FOLLOWENTITY(cam, player, 10.0, 3.0, 5.0)
    
    ; 2. Rendering
    RENDER.CLEAR(20, 20, 30)
    RENDER.BEGIN3D(cam)
        ENTITY.DRAWALL()
        DRAW3D.GRID(100, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND
```

## See also

- [CAMERA.md](CAMERA.md) — dedicated camera reference
- [LIGHT.md](LIGHT.md) — dedicated lighting reference
- [RENDER.md](RENDER.md) — render pipeline and shaders
