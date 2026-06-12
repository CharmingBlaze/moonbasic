# 3D collision, physics, and picking

> Hits in 3D: rigid bodies (Jolt), ray picks, and pure math overlap — when to use each.

**Namespaces:** `PHYSICS` · `PHYSICS3D` · `BODY3D` · `ENTITY` · `PICK` · `COLLISION` · **Status:** Shipped (3D physics: **Windows / Linux** full runtime)

**Commands:** [COMMAND_REGISTRY.md#physics](COMMAND_REGISTRY.md#physics) · [reference/PHYSICS3D.md](../../reference/PHYSICS3D.md) · [reference/PICK.md](../../reference/PICK.md)

---

## Table of contents

- [At a glance](#at-a-glance)
- [When to use 3D collision](#when-to-use-3d-collision)
- [Choose the right tool](#choose-the-right-tool)
- [Physics world workflow](#physics-world-workflow)
- [Bodies on entities](#bodies-on-entities)
- [Ray picking (`PICK.*`)](#ray-picking-pick)
- [Math overlap (`COLLISION.*` 3D)](#math-overlap-collision-3d)
- [Full example](#full-example)
- [Common mistakes](#common-mistakes)
- [See also](#see-also)

---

## At a glance

| Tool | Answers the question | Needs physics world? |
|------|----------------------|----------------------|
| **`PICK.SCREENCAST`** | What did I click in 3D? | Uses scene / physics ray |
| **`PHYSICS3D` + `ENTITY.ADDPHYSICS`** | Blocks fall, stacks, explosions | Yes |
| **`COLLISION.SPHEREOVERLAP3D`** | Are two spheres touching? (math) | No |
| **`BODY.*` aliases** | Quick static/dynamic box on entity | Yes |

**Platform:** Desktop **full runtime** with Jolt. Check scripts anywhere: `moonbasic --check`.

---

## When to use 3D collision

**Use physics when:**

- Objects should fall, slide, and collide without hand-written resolution.
- Stacks, dominoes, projectiles with bounce.

**Use picking when:**

- Mouse selects units, FPS shooting, tool placement.

**Use math `COLLISION.*` when:**

- Cheap distance checks (pickup radius) without waking the physics engine.

**Skip 3D physics when:**

- Pure 2D game — use [COLLISION-2D.md](COLLISION-2D.md).

---

## Choose the right tool

| I want to… | Use | Why |
|------------|-----|-----|
| Click entity in 3D | `PICK.SCREENCAST(cam)` | Returns entity + hit point |
| Crate stack simulation | `PHYSICS.START` + `ENTITY.ADDPHYSICS` | Jolt resolves contacts |
| “Within 2m of hero?” | `COLLISION` distance or `MATH.DISTANCE` | No bodies needed |
| Character walks slopes | `CHAR.*` / `CHARCONTROLLER.*` | KCC vs raw dynamic body |
| Shoot ray from gun | `PICK.CAST` or `PHYSICS3D.RAYCAST` | Direction + max distance |

---

## Physics world workflow

**Why order matters:** The world must exist before bodies; simulation must **step** before you read positions.

1. `PHYSICS.START()` — **Why:** Allocates Jolt world (alias of `PHYSICS3D.START` on desktop).
2. `PHYSICS.SETGRAVITY(0, -9.8, 0)` — **Why:** Defines “down” for dynamics.
3. Create **static** floor entity + `ADDPHYSICS` + `SETSTATIC(true)`.
4. Create **dynamic** props + `BODY.ADDDYNAMICBOX` or `ADDPHYSICS`.
5. Each frame: gameplay forces → `PHYSICS.STEP()` → draw at entity transforms.
6. `PHYSICS.STOP()` on exit.

```basic
PHYSICS.START()
PHYSICS.SETGRAVITY(0, -9.81, 0)

floor = ENTITY.CREATECUBE(20, 1, 20)
ENTITY.ADDPHYSICS(floor)
ENTITY.SETSTATIC(floor, true)

ball = ENTITY.CREATESPHERE(0.5)
ball.pos(0, 5, 0)
BODY.ADDSPHERE(ball, 0.5)

WHILE NOT APP.SHOULDCLOSE()
    PHYSICS.STEP()
    ; entity transforms updated from physics
    RENDER.BEGIN(cam)
    SCENE.DRAW()
    RENDER.END()
    RENDER.FRAME()
WEND
PHYSICS.STOP()
```

---

## Bodies on entities

**Why attach physics to entities:** One ID for draw + simulation; `PHYSICS.STEP` syncs transform back to the entity.

| Command | Why |
|---------|-----|
| `ENTITY.ADDPHYSICS(ent)` | Default dynamic body from entity shape |
| `ENTITY.SETSTATIC(ent, true)` | Infinite mass floor/wall |
| `BODY.SETMASS` / `SETFRICTION` / `SETBOUNCE` | Tune feel |
| `BODY.APPLYIMPULSE` | Jump, explosion kick |

Checklist aliases `BODY.ADDSTATICBOX` etc. map to the same paths — see [05-PHYSICS.md](../05-PHYSICS.md).

---

## Ray picking (`PICK.*`)

**Why:** Converts mouse screen position into a 3D ray through the camera — answers “what entity is under the cursor?”

**Workflow:**

```basic
PICK.SCREENCAST(cam)
IF PICK.HIT() THEN
    hitEnt = PICK.ENTITY()
    hx = PICK.X()
    hy = PICK.Y()
    hz = PICK.Z()
    dist = PICK.DIST()
ENDIF
```

| Command | Role |
|---------|------|
| `PICK.SCREENCAST(cam)` | Mouse ray from camera |
| `PICK.CAST(ox,oy,oz,dx,dy,dz)` | Custom ray (AI, lasers) |
| `PICK.HIT()` | Any hit this frame? |
| `PICK.ENTITY()` | Which entity |

Call **after** you know mouse position; typically once per frame on click or hover.

---

## Math overlap (`COLLISION.*` 3D)

**Why:** O(1) tests without creating bodies — good for pickups, aggro radius.

Examples (see [reference/COLLISION.md](../../reference/COLLISION.md)):

- `COLLISION.SPHEREOVERLAP3D(centerA, radiusA, centerB, radiusB)`
- `COLLISION.BOXOVERLAP3D` — AABB tests with `VEC3` handles

Use **`VEC3.CREATE`** for centers; compare distance to threshold for simple “magnet” pickups.

---

## Full example

Click to impulse a sphere ([`examples/sphere_drop`](../../../examples/sphere_drop/main.mb) is a larger sample).

```basic
APP.OPEN(800, 600, "3D pick + physics")
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
ball.pos(0, 4, 0)
BODY.ADDSPHERE(ball, 0.5)

WHILE NOT APP.SHOULDCLOSE()
    PHYSICS.STEP()
  PICK.SCREENCAST(cam)
    IF INPUT.MOUSEHIT(MOUSE_LEFT) AND PICK.HIT() THEN
        BODY.APPLYIMPULSE(ball, 0, 8, 0)
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

## Common mistakes

| Mistake | Fix |
|---------|-----|
| No `PHYSICS.STEP()` | Bodies never move |
| Dynamic floor | Use `SETSTATIC(true)` on terrain |
| Pick before camera ready | Set active camera first |
| Mix manual `SETPOS` every frame on dynamic body | Fight physics — use forces or kinematic mode |
| Expect physics on compiler-only zip | Need full runtime |

---

## See also

- [ENTITY-SYSTEM.md](ENTITY-SYSTEM.md)
- [05-PHYSICS.md](../05-PHYSICS.md)
- [reference/CHARACTER_PHYSICS.md](../../reference/CHARACTER_PHYSICS.md) — walking characters
