# Physics systems: PHYSICS, BODY, COLLISION, PICK

> 3D physics world, rigid bodies, simple collision helpers, and screen/world picking.

**Platform:** **Windows** and **Linux** full runtime include Jolt 3D physics. Other platforms may return stub errors — use **`moonbasic --check`** to validate scripts without running physics.

**All commands:** [COMMAND_REGISTRY.md#physics](COMMAND_REGISTRY.md#physics)

**Deep guides:** [guides/COLLISION-2D.md](guides/COLLISION-2D.md) · [guides/COLLISION-3D.md](guides/COLLISION-3D.md) · [guides/PHYSICS-2D-PLATFORMER.md](guides/PHYSICS-2D-PLATFORMER.md) · [guides/CHARACTER-3D-WALKING.md](guides/CHARACTER-3D-WALKING.md)

**See also:** [01-CORE](01-CORE.md) · [reference/PHYSICS3D.md](../reference/PHYSICS3D.md) · [reference/PICK.md](../reference/PICK.md)

---

## Table of contents

- [PHYSICS system](#physics-system)
- [BODY system](#body-system)
- [COLLISION system](#collision-system)
- [PICK system](#pick-system)
- [Full example](#full-example)
- [Memory notes](#memory-notes)
- [See also](#see-also)

---

## PHYSICS system

World-level simulation control.

### Core workflow

1. `PHYSICS.START()` — initialize the 3D physics world.
2. `PHYSICS.SETGRAVITY(x, y, z)` — default `(0, -9.81, 0)`.
3. Each frame: `PHYSICS.STEP()` or `PHYSICS.STEP(dt)` after gameplay logic.
4. `PHYSICS.STOP()` on shutdown.

**Aliases:** checklist `PHYSICS.CREATEWORLD` → `START`

---

### `PHYSICS.START()` / `PHYSICS.STOP()`

Starts or stops the physics session.

**Returns:** nothing

**Example:**

```basic
PHYSICS.START()
PHYSICS.SETGRAVITY(0, -9.8, 0)
```

---

### `PHYSICS.SETGRAVITY(x, y, z)`

Sets world gravity vector.

**Example:**

```basic
PHYSICS.SETGRAVITY(0, -20, 0)
```

---

### `PHYSICS.STEP([dt])`

Advances simulation one tick. Call once per frame in the common case.

**Example:**

```basic
PHYSICS.STEP()
```

---

### `PHYSICS.ENABLE()` / `PHYSICS.DISABLE()`

Toggle simulation without destroying bodies.

---

## BODY system

Checklist **`BODY.*`** names are **aliases** for `ENTITY.ADDPHYSICS` plus property helpers — attach physics to entities.

### Core workflow

1. Create an entity (`ENTITY.CREATECUBE`, …).
2. `ENTITY.ADDPHYSICS(entity [, shape])` — dynamic body.
3. Or use **`BODY.ADDSTATICBOX`** / **`ADDDYNAMICBOX`** / **`ADDSPHERE`** / **`ADDCAPSULE`** aliases.
4. Set mass, friction, bounce; apply forces.
5. Physics step syncs entity transforms automatically.

---

### Body alias commands

| Command | Description |
|---------|-------------|
| `BODY.ADDSTATICBOX(ent, w, h, d)` | Static box collider |
| `BODY.ADDDYNAMICBOX(ent, w, h, d)` | Dynamic box |
| `BODY.ADDSPHERE(ent, radius)` | Sphere collider |
| `BODY.ADDCAPSULE(ent, radius, height)` | Capsule collider |
| `BODY.SETMASS(ent, mass)` | Mass in kg |
| `BODY.SETFRICTION(ent, value)` | Friction coefficient |
| `BODY.SETBOUNCE(ent, value)` | Restitution |
| `BODY.APPLYFORCE(ent, fx, fy, fz)` | Continuous force |
| `BODY.APPLYIMPULSE(ent, ix, iy, iz)` | Instant impulse |

**Canonical:** `ENTITY.ADDPHYSICS(entity)` — see [reference/ENTITY.md](../reference/ENTITY.md).

**Example:**

```basic
PHYSICS.START()
floor = ENTITY.CREATECUBE(10, 1, 10)
ENTITY.ADDPHYSICS(floor)
ENTITY.SETSTATIC(floor, true)

crate = ENTITY.CREATECUBE(1, 1, 1)
crate.pos(0, 3, 0)
BODY.ADDDYNAMICBOX(crate, 1, 1, 1)
BODY.SETMASS(crate, 5)
```

---

## COLLISION system

**Status:** Partial — full beginner slide/collision-rules API from the checklist is not all shipped. Use **physics bodies** for dynamic games, or geometry overlap helpers in [reference/COLLISION.md](../reference/COLLISION.md).

| Approach | When to use |
|----------|-------------|
| `PHYSICS` + `BODY` / `ENTITY.ADDPHYSICS` | 3D rigid bodies, stacks, projectiles |
| `PHYSICS2D` | 2D platformers — [reference/PHYSICS2D.md](../reference/PHYSICS2D.md) |
| `ENTITY` overlap queries | Simple distance / AABB tests without full physics |

Checklist names like `COLLISION.RULE` and `COLLISION.SLIDE` are planned; document gaps in [FINAL_POLISH_SYSTEMS.md](../FINAL_POLISH_SYSTEMS.md).

---

## PICK system

Ray casting for mouse clicks and gameplay traces.

### Core workflow

1. Configure ray: `PICK.SCREENCAST(camera)` or `PICK.CAST(ox, oy, oz, dx, dy, dz)`.
2. Read hit: `PICK.HIT()`, `PICK.ENTITY()`, `PICK.X/Y/Z`, `PICK.DIST`.
3. Optional: `PICK.FROMCAMERA`, layer mask, max distance setters.

**Aliases:** checklist `PICK.MOUSE` → `PICK.SCREENCAST` / `PICK.FROMCAMERA`

---

### `PICK.SCREENCAST(camera)`

Casts from the camera through the mouse cursor into the scene.

| Argument | Type | Description |
|----------|------|-------------|
| camera | handle | Active camera |

**Returns:** nothing (fills pick state)

**Example:**

```basic
PICK.SCREENCAST(cam)
IF PICK.HIT() THEN
    hitEnt = PICK.ENTITY()
    px = PICK.X()
ENDIF
```

---

### `PICK.CAST(ox, oy, oz, dx, dy, dz)`

World-space ray from origin and direction.

**Example:**

```basic
PICK.CAST(0, 0, 0, 0, 0, 1)
```

---

### Hit queries

| Command | Returns |
|---------|---------|
| `PICK.HIT()` | `bool` — something was hit |
| `PICK.ENTITY()` | `handle` — hit entity |
| `PICK.X()` / `Y()` / `Z()` | `float` — hit point |
| `PICK.NX()` / `NY()` / `NZ()` | `float` — hit normal |
| `PICK.DIST()` | `float` — distance along ray |

**Aliases:** checklist `PICK.DISTANCE` → `PICK.DIST`

---

## Full example

```basic
APP.OPEN(800, 600, "Physics + Pick")
APP.SETFPS(60)

PHYSICS.START()
PHYSICS.SETGRAVITY(0, -9.8, 0)

cam = CAMERA.CREATE()
CAMERA.SETACTIVE(cam)
CAMERA.SETPOS(cam, 0, 4, -10)
CAMERA.LOOKAT(cam, 0, 0, 0)

ground = ENTITY.CREATECUBE(20, 1, 20)
ENTITY.ADDPHYSICS(ground)
ENTITY.SETSTATIC(ground, true)

ball = ENTITY.CREATESPHERE(0.5)
ball.pos(0, 5, 0)
BODY.ADDSPHERE(ball, 0.5)

WHILE NOT APP.SHOULDCLOSE()
    PHYSICS.STEP()
    PICK.SCREENCAST(cam)
    IF INPUT.MOUSEHIT(MOUSE_LEFT) AND PICK.HIT() THEN
        BODY.APPLYIMPULSE(ball, 0, 5, 0)
    ENDIF

    RENDER.CLEAR(18, 20, 28)
    RENDER.BEGIN(cam)
    SCENE.DRAW()
    RENDER.END()
    RENDER.FRAME()
WEND

PHYSICS.STOP()
APP.CLOSE()
```

---

## Memory notes

- `PHYSICS.STOP` before exit; static bodies freed with `ENTITY.FREE`.
- Pick state is per-frame — read hit fields immediately after `SCREENCAST` / `CAST`.

---

## See also

- [examples/sphere_drop](../examples/sphere_drop/main.mb) — Jolt demo
- [reference/CHARACTER.md](../reference/CHARACTER.md) — kinematic characters (KCC)
- [04-INPUT](04-INPUT.md) — mouse for picking
