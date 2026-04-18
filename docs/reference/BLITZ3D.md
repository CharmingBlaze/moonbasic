# Blitz3D-Style Helpers

A high-level command surface inspired by the classic Blitz3D engine, providing familiar globals and ergonomic wrappers for entities, cameras, and input.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. **Initialization:** Use `WINDOW.OPEN` followed by `CAMERA.CREATE` (or the short `CAM()` alias).
2. **Entity Creation:** Use `ENTITY.CREATECUBE` or short primitives like `CUBE()` and `SPHERE()`.
3. **Logic:** Drive movement each frame using `Entity.Move` and `Entity.Turn`.
4. **Render Pass:**
   - `RENDER.CLEAR` to reset the frame.
   - `RENDER.BEGIN3D(cam)` to start the 3D pass.
   - `ENTITY.DRAWALL` to render all active entities.
   - `RENDER.END3D` and `RENDER.FRAME` to finalize.

---

## 1. Global Entities (Short Primitives)

MoonBASIC provides short constructor aliases that return an `EntityRef` handle, enabling modern **Method Chaining**.

| Command | Arguments | Returns | Notes |
|---------|-----------|---------|-------|
| `CUBE(w, h, d)` | Float... | Handle | Creates a box entity. |
| `SPHERE(radius)` | Float | Handle | Creates a sphere entity. |
| `PLANE(size)` | Float | Handle | Creates an XZ plane tile. |
| `CAM()` | None | Handle | Alias for `CAMERA.CREATE`. |

### Fluent Method Chaining
Once you have an `EntityRef` (like from `CUBE()`), you can configure it fluently:
```basic
player = CUBE(1, 2, 1).POS(0, 1, 0).COLOR(255, 0, 0).SETTAG("Player")
```

---

## 2. Camera Helpers

| Command | Arguments | Returns | Notes |
|---------|-----------|---------|-------|
| `Camera.SetPos(c, x, y, z)` | Handle, Float...| Handle | Positions the camera. |
| `Camera.LookAt(c, x, y, z)` | Handle, Float...| None | Aims at a world point. |
| `Camera.Orbit(c, e, dist)` | Handle, Int... | None | Follow-orbit an entity. |
| `Camera.Zoom(c, amount)` | Handle, Float | None | Adjusts FOV in degrees. |

---

## 3. Entity Management (`Entity.*`)

These commands work with standard **integer entity IDs**.

| Command | Arguments | Returns | Notes |
|---------|-----------|---------|-------|
| `Entity.Move(e, f, r, u)` | Int, Float... | None | Local-space movement. |
| `Entity.Turn(e, p, y, r)` | Int, Float... | None | Incremental rotation (radians). |
| `Entity.SetPos(e, x, y, z)` | Int, Float... | None | Absolute world position. |
| `Entity.Parent(c, p)` | Int, Int | None | Attaches `c` to `p`. |
| `Entity.Free(e)` | Int | None | Unloads and destroys entity. |

---

## Full Example

A classic Blitz3D-style rotating cube demo using modern MoonBASIC fluent chaining.

```basic
WINDOW.OPEN(1280, 720, "Blitz3D Modern Style")
cam = CAM().POS(0, 5, -10).LOOK(0, 0, 0)
cube = CUBE(2, 2, 2).COLOR(0, 200, 255)

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    
    ; Rotate using the increment-based 'Turn' method
    cube.TURN(0, 1.5 * dt, 0)
    
    RENDER.CLEAR(15, 15, 25)
    RENDER.BEGIN3D(cam)
        ENTITY.DRAWALL()
        DRAW3D.GRID(50, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND

cube.FREE()
cam.FREE()
WINDOW.CLOSE()
```

## See also

- [BLITZ_COMMAND_INDEX.md](BLITZ_COMMAND_INDEX.md) — full command mapping
- [ENTITY.md](ENTITY.md) — core entity API reference
- [CAMERA.md](CAMERA.md) — advanced camera controls
- [INPUT.md](INPUT.md) — Keyboard, mouse, and gamepad actions
- [COLLISION.md](COLLISION.md) — Box-to-land and overlap tests
