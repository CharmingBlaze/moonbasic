# Beginner “Full Stack” Gameplay Helpers

A collection of high-level, friendly commands that bridge input, world interactions, and entity logic for rapid game prototyping.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. **Input:** Poll mouse or keyboard state using `INPUT.*` or bridge helpers.
2. **Spatial:** Convert screen clicks to world positions using `WORLD.MOUSEFLOOR3D`.
3. **Entity:** Drive behavior using high-level routines like `ENTITY.NAVTO` or `ENTITY.PATROL`.
4. **RPG Logic:** Manage life and combat with `ENTITY.SETHEALTH` and `ENTITY.DAMAGE`.
5. **Juice:** Apply polish with `CAMERA.SHAKE` or `WORLD.FLASH`.

---

## 1. Input and Mouse (2D / 3D Bridge)

| Command | Arguments | Returns | Notes |
|---------|-----------|---------|-------|
| `INPUT.MOUSEX()` | None | Float | Window client X in pixels. |
| `INPUT.MOUSEY()` | None | Float | Window client Y in pixels. |
| `INPUT.MOUSEDELTAX()`| None | Float | Delta since last frame. |
| `INPUT.LOCKMOUSE(toggle)`| Boolean | None | Locks/unlocks cursor for FPS mouselook. |
| `WORLD.MOUSE2D(cam)` | Handle | Array | World `[x, y]` under cursor. |
| `WORLD.MOUSEFLOOR3D(cam, y)`| Handle, Float | Array | Ray vs plane at `y` → `[x, z]`. |
| `WORLD.MOUSETOENTITY(cam)`| Handle | Integer | Entity ID under cursor (**CGO/Jolt**). |
| `PHYSICS3D.MOUSEHIT(cam)`| Handle | Array | World `[x, y, z]` of physics hit. |

---

## 2. Navigation and Movement

These commands drive **scripted** entities. Use `ENTITY.STOP(e)` to cancel any active navigation or patrol.

| Command | Arguments | Returns | Notes |
|---------|-----------|---------|-------|
| `ENTITY.MOVEFORWARD(e, s)`| Integer, Float | Handle | Moves along facing direction. |
| `ENTITY.NAVTO(e, x, z, s)`| Integer, Float... | Handle | Pathfinds to destination. |
| `ENTITY.WALKTO(e, x, z, s)`| Integer, Float... | Handle | `NAVTO` with arrival logic. |
| `ENTITY.PATROL(e, ax, az, bx, bz, s)`| Integer, Float... | Handle | Ping-pongs between two points. |
| `ENTITY.STOP(e)` | Integer | Handle | Cancels all active movement. |

---

## 3. Combat and Tags (RPG-Style)

| Command | Arguments | Returns | Notes |
|---------|-----------|---------|-------|
| `ENTITY.SETHEALTH(e, v)`| Integer, Float | Handle | Sets current health. |
| `ENTITY.DAMAGE(e, v)` | Integer, Float | Handle | Reduces health by `v`. |
| `ENTITY.ISALIVE(e)` | Integer | Boolean | `TRUE` if health > 0. |
| `ENTITY.SETTAG(e, tag)`| Integer, String| Handle | Attaches a string label. |
| `ENTITY.FINDNEARESTWITHTAG(e, tag)`| Integer, String| Integer | Nearest entity with label. |

---

## 4. World and Juice

| Command | Arguments | Returns | Notes |
|---------|-----------|---------|-------|
| `WORLD.SETGRAVITY(x, y, z)`| Float, Float, Float| None | Sets global physics gravity. |
| `WORLD.SETTIMESCALE(s)`| Float | None | Slow-motion or speed-up factor. |
| `WORLD.EXPLOSION(x, y, z, f, r)`| Float... | None | Radial physics push (**CGO**). |
| `WORLD.FLASH(color, d)`| Handle, Float | None | Full-screen color flash. |
| `CAMERA.SHAKE(cam, i, d)`| Handle, Float, Float| None | Screen shake effect. |

---

## Full Example

A simple "Click to Move" script with health management and a camera shake effect on damage.

```basic
WINDOW.OPEN(1280, 720, "Full Stack Demo")
cam = CAMERA.CREATE()
player = ENTITY.CREATECUBE(1.0).SETHEALTH(100).SETTAG("Player")

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    
    ; 1. Input: Click to navigate
    IF INPUT.MOUSEHIT(1)
        target = WORLD.MOUSEFLOOR3D(cam, 0)
        IF target <> NIL
            player.NAVTO(target(0), target(1), 5.0)
        END IF
    END IF

    ; 2. Combat: Self-damage for testing
    IF INPUT.KEYHIT(KEY_SPACE)
        player.DAMAGE(10)
        CAMERA.SHAKE(cam, 0.5, 0.2)
        PRINT "Health: " + STR(player.GETHEALTH())
    END IF

    ; 3. Update and Render
    ENTITY.UPDATE(dt)
    CAMERA.FOLLOW(cam, player, 10.0)
    
    RENDER.CLEAR(20, 20, 30)
    RENDER.BEGIN3D(cam)
        ENTITY.DRAWALL()
        DRAW3D.GRID(50, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND
```

## See also

- [GAMEPLAY_HELPERS.md](GAMEPLAY_HELPERS.md) — triggers and proximity
- [PROJECTILES.md](PROJECTILES.md) — shooting and pooling
- [ENTITY.md](ENTITY.md) — core entity documentation
