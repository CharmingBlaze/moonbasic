# Blitz 2025: Modern Ergonomics

A vision for modern game development in MoonBASIC, combining the rapid prototyping speed of Blitz3D with modern engine features like Jolt Physics, JSON state, and Fluent APIs.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. **Modern Types:** Use `MAP` and `LIST` handles instead of legacy `Type` structures.
2. **Physics First:** Use `ENTITY.ADDPHYSICS` (or `.PHYSICS()`) to let Jolt handle motion and collision.
3. **Data Driven:** Load levels and config from JSON using the `JSON.*` namespace.
4. **Fluent API:** Chain your entity configuration for readable, compact code.

---

## 1. Modern Data Structures

Forget legacy globals. Use handles for dynamic state.

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| `MAP()` | None | Handle | Creates a key-value dictionary. |
| `LIST()`| None | Handle | Creates a dynamic array. |
| `JSON.LOAD(path)` | String | Handle | Parses a JSON file into a Map/List. |

---

## 2. Jolt Physics Integration

The "Blitz 2025" way is to let the physics engine solve the world.

| Method | Arguments | Returns | Description |
|--------|-----------|---------|-------------|
| `.PHYSICS(mode)`| Int | Handle | Adds a Jolt body to an entity. |
| `.IMPULSE(x,y,z)`| Float... | Handle | Applies a physics impulse. |
| `.GRAVITY(g)` | Float | Handle | Sets per-entity gravity scale. |

---

## 3. Advanced Scene Control

| Command | Arguments | Returns | Description |
|---------|-----------|---------|-------------|
| `SCENE.LOAD(path)`| String | Handle | Loads a GLTF/GLB scene with entities. |
| `FOG.COLOR(r,g,b)`| Int... | None | Sets global distance fog. |
| `SKYBOX(path)` | String | Handle | Sets the environment cubemap. |

---

## Full Example: The 2025 Character Controller

A modern character controller using Jolt physics and method chaining.

```basic
WINDOW.OPEN(1280, 720, "Blitz 2025")
cam = CAM().POS(0, 5, -10).LOOK(0, 0, 0)

; Create a player with Jolt physics (Mode 2 = Dynamic)
player = SPHERE(0.5).POS(0, 10, 0).COLOR(0, 255, 100).PHYSICS(2)

; Create a floor
floor = CUBE(20, 1, 20).POS(0, -0.5, 0).COLOR(50, 50, 50).PHYSICS(0)

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    
    ; Physics is solved automatically in the background
    ; We just apply forces based on input
    IF KEYDOWN(KEY_W) THEN player.IMPULSE(0, 0, 500 * dt)
    
    RENDER.CLEAR(10, 10, 15)
    RENDER.BEGIN3D(cam)
        ENTITY.DRAWALL()
    RENDER.END3D()
    RENDER.FRAME()
WEND

WINDOW.CLOSE()
```

## See also

- [PHYSICS3D.md](PHYSICS3D.md) — Deep dive into Jolt integration
- [JSON.md](JSON.md) — Working with modern data formats
- [NETWORKING.md](NETWORKING.md) — Building multiplayer "2025" games
