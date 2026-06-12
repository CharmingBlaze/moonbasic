# 3D walking characters — kinematic controllers

> Move a humanoid through a 3D level with slopes, steps, and jumps — **without** tumbling like a loose physics crate.

**Namespaces:** `CHAR` · `CHARACTER` · `CHARCONTROLLER` · **Status:** Shipped · **Platform:** Windows / Linux **full runtime** (Jolt KCC)

**Commands:** [COMMAND_REGISTRY.md#physics](../COMMAND_REGISTRY.md#physics) · [reference/CHARACTER_PHYSICS.md](../../reference/CHARACTER_PHYSICS.md)

**Compare:** Raw dynamic body → [COLLISION-3D.md](COLLISION-3D.md)

---

## Table of contents

- [At a glance](#at-a-glance)
- [When to use KCC vs rigid body](#when-to-use-kcc-vs-rigid-body)
- [Core workflow](#core-workflow)
- [Key commands](#key-commands)
- [Full example sketch](#full-example-sketch)
- [Common mistakes](#common-mistakes)
- [See also](#see-also)

---

## At a glance

| Idea | Detail |
|------|--------|
| **Problem** | Dynamic `BODY3D` spheres wobble and tip — bad for FPS/third-person heroes |
| **Solution** | **Kinematic Character Controller (KCC)** — `CHAR.*` on desktop |
| **You need first** | `PHYSICS3D.START`, entity with model or capsule |
| **Sample** | [`examples/mario64/modern_blitz_hop_kcc.mb`](../../../examples/mario64/modern_blitz_hop_kcc.mb) |

**Why KCC exists:** Gameplay wants **responsive** walk, **snap to ground**, **step up curbs**, and **jump** — Jolt’s character virtual controller solves that better than raw rigid bodies.

---

## When to use KCC vs rigid body

| Use `CHAR.*` (KCC) | Use `ENTITY.ADDPHYSICS` (dynamic) |
|--------------------|-----------------------------------|
| Player, NPC walkers | Crates, balls, debris |
| Slopes + stairs | Stacks, explosions |
| Camera-relative move | Physics puzzles |

---

## Core workflow

1. `PHYSICS3D.START()` — **Why:** KCC lives inside the physics world.
2. Create **entity** (model or cube placeholder).
3. `CHAR.CREATE(entity [, radius, height])` — **Why:** Attaches capsule controller to entity.
4. Optional: `CHAR.SETSTEP`, `CHAR.SETSLOPE` — **Why:** Stairs and hill limits.
5. Each frame:
   - Read input → `CHAR.MOVE` or `CHAR.MOVEWITHCAMERA`
   - Jump if `CHAR.ISGROUNDED()`
   - `PHYSICS3D.STEP()` or unified physics step
6. Draw entity at updated transform.

```basic
PHYSICS3D.START()
hero = ENTITY.LOAD("assets/hero.glb")
CHAR.CREATE(hero, 0.4, 1.8)
CHAR.SETSTEP(hero, 0.35)
CHAR.SETSLOPE(hero, 45)

WHILE NOT APP.SHOULDCLOSE()
    ; camera-relative walk
    CHAR.MOVEWITHCAMERA(hero, INPUT.MOVEDIR(), 6 * APP.DELTA())
    IF CHAR.ISGROUNDED(hero) AND ACTION.HIT("Jump") THEN CHAR.JUMP(hero, 6)
    PHYSICS3D.STEP()
    RENDER.BEGIN(cam)
    SCENE.DRAW()
    RENDER.END()
    RENDER.FRAME()
WEND
```

(Exact jump API: `CHAR.JUMP` — verify arity in COMMAND_REGISTRY.)

---

## Key commands

| Command | Why |
|---------|-----|
| `CHAR.CREATE(ent, r, h)` | Capsule controller on entity |
| `CHAR.MOVE(ent, dx, dy, dz)` | World-space motion |
| `CHAR.MOVEWITHCAMERA(ent, dir, speed)` | WASD relative to camera |
| `CHAR.JUMP(ent, power)` | Vertical impulse when grounded |
| `CHAR.ISGROUNDED(ent)` | Can jump? |
| `CHAR.SETSTEP(ent, h)` | Max stair height |
| `CHAR.SETSLOPE(ent, degrees)` | Max walkable slope |

**Aliases:** `PLAYER.CREATE`, `CHARACTER.CREATE` — same handlers.

---

## Full example sketch

Pair with orbit camera ([CAMERA-AND-INPUT.md](CAMERA-AND-INPUT.md)) and `ACTION.MAPKEY` for bindings.

Run: **`moonrun`** · Check: **`moonbasic --check`**

---

## Common mistakes

| Mistake | Fix |
|---------|-----|
| KCC without `PHYSICS3D.START` | Start physics first |
| Also applying `ENTITY.MOVE` every frame | Let CHAR drive position |
| Jump without `ISGROUNDED` | Double-jump unless designed |
| Compiler-only install | Need full runtime + Jolt |

---

## See also

- [COLLISION-3D.md](COLLISION-3D.md)
- [ENTITY-SYSTEM.md](ENTITY-SYSTEM.md)
- [examples/mario64/README.md](../../../examples/mario64/README.md)
