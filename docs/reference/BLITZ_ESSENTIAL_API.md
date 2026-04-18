# Essential Blitz-Style API

A curated list of the most important Blitz-style commands and their modern MoonBASIC equivalents. For a full mapping, see the [Command Index](BLITZ_COMMAND_INDEX.md).

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. **Short Constructors:** Use `CUBE()`, `SPHERE()`, and `CAM()` for instant object creation.
2. **Method Chaining:** Configure your objects in a single line using the fluent handle methods.
3. **Logic & Render:** Drive the simulation with `ENTITY.UPDATE` and `ENTITY.DRAWALL`.

---

## 1. Primitives & Factories

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| `CUBE(size)` | Float | Handle | Creates a box entity. |
| `SPHERE(radius)`| Float | Handle | Creates a sphere entity. |
| `CAM()` | None | Handle | Creates a 3D camera. |
| `TEX(path)` | String | Handle | Loads a texture. |

---

## 2. Fluent Entity Methods

Most commands return the primary handle to support **Method Chaining**.

| Method | Arguments | Returns | Description |
|--------|-----------|---------|-------------|
| `.POS(x, y, z)` | Float... | Handle | Sets world position. |
| `.ROT(p, y, r)` | Float... | Handle | Sets rotation (radians). |
| `.SCALE(x, y, z)`| Float... | Handle | Sets non-uniform scale. |
| `.COLOR(r, g, b)`| Int... | Handle | Sets RGB tint. |
| `.TURN(p, y, r)` | Float... | Handle | Incremental rotation. |

---

## 3. High-Level Logic

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| `KEYHIT(key)` | Int | Bool | True on initial press. |
| `KEYDOWN(key)` | Int | Bool | True while held. |
| `RND(min, max)` | Float... | Float | Random range value. |
| `LERP(a, b, t)` | Float... | Float | Linear interpolation. |

---

## Full Example

```basic
WINDOW.OPEN(1280, 720, "Essential Blitz")
cam = CAM().POS(0, 5, -10).LOOK(0, 0, 0)
ball = SPHERE(1.0).POS(0, 2, 0).COLOR(255, 100, 0)

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    
    ; Logic
    ball.TURN(0, 2.0 * dt, 0)
    
    RENDER.CLEAR(20, 20, 30)
    RENDER.BEGIN3D(cam)
        ENTITY.DRAWALL()
    RENDER.END3D()
    RENDER.FRAME()
WEND

WINDOW.CLOSE()
```

## See also

- [BLITZ3D.md](BLITZ3D.md) — Main Blitz-to-MoonBASIC guide
- [BLITZ_COMMAND_INDEX.md](BLITZ_COMMAND_INDEX.md) — Full command mapping
- [PROGRAMMING.md](../PROGRAMMING.md) — Core engine patterns
