# Beginner “full stack” gameplay helpers

This page maps **friendly names** (Input / World / Math / Entity dot methods) to the **registry keys** moonBASIC actually runs. Everything here is implemented in the Go runtime unless noted.

**Related:** [GAMEPLAY_HELPERS.md](GAMEPLAY_HELPERS.md), [PROJECTILES_AND_POOLING.md](PROJECTILES_AND_POOLING.md), [MEMORY.md](../MEMORY.md), [POOL.md](POOL.md).

---

## 1. Input and mouse (2D / 3D bridge)

| You want | Registry command | Notes |
|----------|------------------|--------|
| Mouse X / Y in pixels | `INPUT.MOUSEX()`, `INPUT.MOUSEY()` | Window client pixels |
| Delta since last frame | `INPUT.MOUSEDELTAX()`, `INPUT.MOUSEDELTAY()` | FPS mouselook |
| Lock / unlock cursor | `INPUT.LOCKMOUSE(TRUE)` / `FALSE` | Calls Raylib `DisableCursor` / `EnableCursor` |
| 2D world under cursor | `WORLD.MOUSE2D(camera2D)` | → float array `[wx, wy]` |
| Ground under cursor (3D) | `WORLD.MOUSEFLOOR3D(cam, floorY#)` or **`WORLD.MOUSETOFLOOR`** (alias) | Ray vs plane **y = floorY** → `[wx, wz]` or **NIL** |
| Entity under cursor (Jolt) | **`WORLD.MOUSETOENTITY(cam)`** or `CAMERA.RAYCASTMOUSE(cam)` | → **entity#** or **0** (needs **Linux + CGO + Jolt**; Windows stub errors) |
| Hit **point** under cursor | `PHYSICS3D.MOUSEHIT(cam)` | → `[x,y,z]` or **NIL** (same platform note) |

**Minimal 3D click-to-ground pattern:** after `CAMERA.BEGIN(cam)`, `pt = WORLD.MOUSEFLOOR3D(cam, 0)` — if not NIL, read index 0 and 1 from the float array for `wx` / `wz`, then drive `ENTITY.NAVTO` or your own logic.

---

## 2. Navigation and movement (non-physics movers)

These target **scripted** entities (not `ENTITY.ADDPHYSICS` Jolt bodies). Use **`ENTITY.STOP`** to cancel nav / patrol / magnet and zero velocity.

| Intent | Command | Returns `EntityRef` for chaining? |
|--------|---------|-----------------------------------|
| Move along facing | `ENTITY.MOVEFORWARD(entity, speed#)` | Yes (handle) |
| Click-to-move | `ENTITY.NAVTO(entity, tx#, tz#, speed#)` | Optional 5th **arrival** distance; yes |
| Walk with default arrival | `ENTITY.WALKTO(entity, tx#, tz#, speed# [, arrival#])` | Yes |
| Ping-pong two XZ points | `ENTITY.PATROL(entity, ax#, az#, bx#, bz#, speed#)` | Yes (6-arg form). Legacy **3-arg** `PATROL` still sets internal AI patrol |
| Keep distance band | `ENTITY.KEEPDISTANCE(entity, target, min#, max#, speed#)` | See [GAMEPLAY_HELPERS.md](GAMEPLAY_HELPERS.md) |
| Thrust (Jolt, Linux) | `ENTITY.THRUSTFORWARD` / `THRUSTUP` / `BRAKE` | Dot: `Thrust` → `THRUSTFORWARD` |
| Stop everything | `ENTITY.STOP(entity)` | Yes |

---

## 3. Combat and tags (RPG-style)

| Intent | Command |
|--------|---------|
| HP | `ENTITY.SETHEALTH(entity, max#)` or **`(entity, current#, max#)`** |
| Damage | `ENTITY.DAMAGE(entity, amount#)` |
| Alive | `ENTITY.ISALIVE(entity)` → bool |
| Tag | `ENTITY.SETTAG(entity, tag$)` — returns handle |
| Nearest tagged | `ENTITY.FINDNEARESTWITHTAG(entity, tag$)` → **entity#** |
| Loot on death | `ENTITY.ONDEATHDROP(entity, prefabEntity, chance%)` |

Also: `ENTITY.GETCLOSESTWITHTAG(entity#, radius#, tag$)` when you need a **radius** limit.

---

## 4. World, camera, juice

| Intent | Command |
|--------|---------|
| Global gravity | `PHYSICS3D.SETGRAVITY(gx, gy, gz)` or **`WORLD.SETGRAVITY`** (alias) |
| Slow-mo | `GAME.SETTIMESCALE(factor)` or **`WORLD.SETTIMESCALE`** |
| Radial push | `PHYSICS.EXPLOSION(x, y, z, force, radius)` or **`WORLD.EXPLOSION`** |
| Full-screen flash | `WORLD.FLASH(color, duration)` — see [WINDOW.md](WINDOW.md) |
| Camera shake | `CAMERA.SHAKE(camera, intensity, duration)` |
| Look at entity | `CAMERA.LOOKATENTITY(camera, entity#)` — dot on camera: **`LookAtEntity`** |

---

## 5. Math (no heavy calculus)

| Intent | Command |
|--------|---------|
| Random % | `MATH.CHANCE(percent#)` → bool |
| Random float range | `MATH.RANGE(min#, max#)` (alias of `MATH.RNDF`) |
| Snap to grid | `MATH.SNAP(value#, grid#)` |
| Wrap range | `MATH.WRAP(value#, min#, max#)` |
| Smooth toward | `MATH.SMOOTH(current#, target#, speed#)` — uses frame **dt** internally |

Also: `MATH.LERP`, `MATH.APPROACH`, `MATH.CLAMP` — see [MATH.md](MATH.md).

---

## 6. Screen / UI bridge (3D → 2D)

| Intent | Command |
|--------|---------|
| World → screen array | `ENTITY.WORLDTOSCREEN(entity)` → `[sx, sy]` |
| Component | `ENTITY.TOSCREENX(entity)`, `ENTITY.TOSCREENY(entity)` |
| On screen | `ENTITY.ISONSCREEN(entity)` |
| HP bar | `ENTITY.DRAWHEALTHBAR(entity [, offsetY#])` — requires `ENTITY.SETHEALTH` first; call during your **2D / UI** pass |

---

## 7. Lifecycle and spawners

| Intent | Command |
|--------|---------|
| Timed despawn | `ENTITY.DESTROYAFTER(entity, sec#)` or **`ENTITY.SETLIFETIMER`** (alias) — returns **handle** |
| Parent | `ENTITY.PARENT(child, parent [, global])` — returns handle |
| Unparent | `ENTITY.PARENTCLEAR` / `UNPARENT` or **`ENTITY.DETACH`** |
| Spawner | `SPAWNER.CREATE(prefab, interval#)` at origin XZ **0**, or **`(prefab, interval, x#, z#)`** |

---

## 8. Chaining (fluent style)

Many **mutating** entity builtins now return the **`EntityRef` handle** so you can write:

```moonbasic
player.SetHealth(100, 100).SetTag("Hero").Pos(0, 1, 0)
```

Supported paths include **`ENTITY.SETPOSITION` / `POS`**, **`SETHEALTH`**, **`SETTAG`**, **`NAVTO`**, **`DESTROYAFTER` / `SETLIFETIMER`**, **`PARENT`**, **`PARENTCLEAR` / `DETACH`**, **`SETSTATE`**, and several helpers registered in `runtime/mbentity`.

---

## 9. Example snippet

See **[examples/snippets/beginner_full_stack.mb](../../examples/snippets/beginner_full_stack.mb)** for a commented walkthrough.

---

## 10. Platform notes

- **Jolt** (3D physics, mouse pick, `WORLD.MOUSETOENTITY`): **Linux + CGO** builds. Windows fullruntime uses stubs for native Jolt — use the same **API names**, but expect errors or no-ops where the stub applies.
- **`ENTITY.DRAWHEALTHBAR`** uses immediate-mode 2D drawing; call when your render pass is appropriate for screen-space UI.
