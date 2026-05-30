# Porting from BlitzBASIC 3D

moonBASIC ships a **Blitz3D compatibility layer** (`blitzengine`) so existing Blitz habits map cleanly to modern APIs. This guide shows the same program three ways: **original Blitz3D**, **compat-style moonBASIC**, and **modern moonBASIC**.

For command-level aliases see [BLITZ3D.md](reference/BLITZ3D.md). For entity-first style see [ENTITY.md](reference/ENTITY.md).

---

## Minimal 3D scene

### BlitzBASIC 3D

```basic
Graphics3D 800,600,32,2
SetBuffer BackBuffer()

camera = CreateCamera()
light = CreateLight()

cube = CreateCube()
PositionEntity cube, 0, 0, 5

While Not KeyHit(1)
    TurnEntity cube, 0.1, 0.2, 0.3
    RenderWorld
    Flip
Wend

End
```

### moonBASIC — compat style

Uses familiar names (`CreateCamera`, `RenderWorld`, `KeyHit`) wired to raylib + the entity stack:

```basic
WINDOW.OPEN(800, 600, "Blitz port")
CAMERA.CREATE()
LIGHT.CREATE()

cube = ENTITY.CREATEBOX(1, 1, 1)
ENTITY.POSITION(cube, 0, 0, 5)

WHILE NOT KEYHIT(KEY_ESCAPE)
    ENTITY.TURN(cube, 0.1, 0.2, 0.3)
    RENDER.CLEAR()
    ENTITY.DRAWALL()
    WINDOW.FLIP()
WEND

END
```

### moonBASIC — modern style

Same game with namespaced helpers and interpolation-friendly structure:

```basic
WINDOW.OPEN(800, 600, "Modern port")
cam = CAMERA.CREATE()
LIGHT.CREATE()

cube = ENTITY.CREATEBOX(1, 1, 1)
ENTITY.SET(cube, 0, 0, 5)

WHILE WINDOW.RUNNING()
    IF INPUT.KEYDOWN(KEY_ESCAPE) THEN EXIT

    yaw = yaw + 0.2 * TIME.DELTA()
    ENTITY.SETROT(cube, 0, yaw, 0)

    RENDER.CLEAR()
    ENTITY.DRAWALL()
    WINDOW.FLIP()
WEND

END
```

---

## Common mapping table

| Blitz3D | Compat moonBASIC | Modern moonBASIC |
|---------|------------------|------------------|
| `Graphics3D w,h,d,fs` | `WINDOW.OPEN(w,h,title)` | same |
| `CreateCamera()` | `CAMERA.CREATE()` | `cam = CAMERA.CREATE()` |
| `CreateLight()` | `LIGHT.CREATE()` | same |
| `CreateCube()` | `ENTITY.CREATEBOX(1,1,1)` | same |
| `PositionEntity e,x,y,z` | `ENTITY.POSITION(e,x,y,z)` | `ENTITY.SET(e,x,y,z)` |
| `TurnEntity e,p,y,r` | `ENTITY.TURN(e,p,y,r)` | `ENTITY.SETROT(e,…)` or tween |
| `KeyHit(1)` | `KEYHIT(KEY_ESCAPE)` | `INPUT.KEYHIT(KEY_ESCAPE)` |
| `RenderWorld` | `ENTITY.DRAWALL()` after `RENDER.CLEAR()` | same |
| `Flip` | `WINDOW.FLIP()` | same |
| `MoveEntity` / velocity | `ENTITY.MOVE*` helpers | `BODY3D.*` + physics |
| `Collision` callbacks | string function name | `@MyHandler` or string — see [LANGUAGE.md](LANGUAGE.md) |

---

## Physics and collision

Blitz used **`Collisions`** and entity **`Type`** flags. In moonBASIC:

- **3D rigid bodies:** `BODY3D.MAKE`, `PHYSICS3D.STEP`, `PHYSICS3D.ONCOLLISION(a, b, @Handler)`  
- **2D:** `PHYSICS2D.*` — see examples under `examples/physics/`  
- **Entity queries:** `ENTITY.*` spatial helpers — [ENTITY.md](reference/ENTITY.md)

---

## What is not 1:1

| Blitz concept | moonBASIC note |
|---------------|----------------|
| `Bank` / `Poke` / `Peek` | Use `DIM` arrays or typed `TYPE` fields |
| `Inkey$` polling strings | `INPUT.*` + `KEY_*` constants |
| `EntityParent` hierarchy | Flat handles; use fields on `TYPE` or entity props |
| `LoadMesh` paths relative to CWD | Use **`ASSET.PATH`** / paths relative to your `.mb` file — see [ASSET.md](reference/ASSET.md) |
| `Function` pointers | **`@FuncName`** references (string names still work) |

---

## Suggested port workflow

1. Get the Blitz program running under **compat names** (`ENTITY.*`, `KEYHIT`, `CAMERA.*`).  
2. Replace magic numbers with **`ENUM`** and **`GAMEPAD_*` / `KEY_*`** constants.  
3. Swap string collision callbacks for **`@Handler`** where helpful.  
4. Introduce **`TIME.DELTA()`** for frame-independent motion.  
5. Split shared code into **`INCLUDE "utils.mb"`** or **`IMPORT "mylib"`** when the project grows.

---

## Further reading

- [BLITZ3D.md](reference/BLITZ3D.md) — alias catalog  
- [GETTING_STARTED.md](GETTING_STARTED.md) — install and first run  
- [EXAMPLES.md](EXAMPLES.md) — runnable ports and demos  
- [ROADMAP.md](ROADMAP.md) — what's next for the language and ecosystem  
