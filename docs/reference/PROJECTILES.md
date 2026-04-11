# Projectiles, pooling, and memory

This document ties together **entity lifecycle**, **`Entity.Shoot()`**, **`Pool.*()`**, and what the runtime **does not** hide for you.

## No hidden Go-side entity pool

The **`mbentity`** module does **not** maintain a central pool that recycles **`ent` + Jolt bodies** for bullets. Each shot uses normal allocation: **`Entity.Copy()`** (or `CREATE`) → optional physics build → **`Entity.DestroyAfter()`** / **`Entity.Free()`**, which end in **`purgeEntityByID`** and unload Raylib/Jolt resources in one place.

- **Go heap:** New `ent` / `entExt` churn is visible to the **GC**.
- **GPU / C / Jolt:** Freed only on explicit teardown — see [MEMORY.md](../MEMORY.md).

That is enough for most indie-scale games. For **bullet hell** density, optimize in script (below) or add a dedicated engine pool later (new work).

---

## `Entity.Shoot(prefab, speed, lifetime [, shape])`

Spawns a projectile based on a prefab. Returns a new **entity handle**.

1. **`Entity.Copy(prefab)`** — reloads mesh from disk path when needed; see [ENTITY.md](ENTITY.md).
2. **`Entity.ClearPhysBuffer(id)`** on the clone — removes a duplicated **`physBufIndex`** from the copy.
3. Align **position** and **rotation** to the **shooter**.
4. **`Entity.AddPhysics(id)`** — builds a **new** Jolt dynamic body (shape from optional **`shape`**, default **`SPHERE`**) with **continuous collision detection (CCD)** enabled. **Requires Linux + CGO Jolt**.
5. **`Entity.SetLinearVel(id, vx, vy, vz)`** along shooter **pitch + yaw** facing.
6. **`Entity.DestroyAfter(id, lifetime)`** — **mandatory** positive lifetime so stray bullets are purged.

**Shooter overlap:** Bullets may still hit the shooter if both are dynamic and overlapping. Short-term mitigations: spawn slightly forward, use thicker colliders, or cull hits in gameplay code.

---

## `Body2D.Shoot(shooter, speed, lifetime [, radius])`

Spawns a **small dynamic circle** at the shooter’s front, **`SetBullet(true)`** (Box2D continuous collision), sets linear velocity along the shooter’s **angle**, and schedules **`Body2D.Free()`** after **`lifetime`** seconds.

---

## VM-side `Pool.*()` (different layer)

**`Pool.*()`** pools **VM heap handles** produced by a **factory function** — see [POOL.md](POOL.md). It is **not** an automatic entity+Jolt recycler.

Typical pattern:

- `Pool.Make()` → `Pool.SetFactory()` / `Pool.SetReset()` → `Pool.Prewarm()` (optional) → `Pool.Get()` / `Pool.Return()`.
- Factory might return a **`Model`** handle, an **`Entity`** handle, or another pooled object — you must still match **create/free** rules for **entities** ([MEMORY.md](../MEMORY.md)).

**High-density pattern (script-side):** `Pool.Make("bullets", capacity)` → `Pool.SetFactory(pool, "FactoryFn")` → each frame `handle = Pool.Get(pool)` → reposition, `Entity.SetLinearVel()`, `Entity.Visible(id, TRUE)` → when done, **`Pool.Return(pool, handle)`** instead of `Entity.Free()`.

---

## Reducing churn (practical checklist)

| Technique | Role |
|-----------|------|
| **`Entity.DestroyAfter()` / `Entity.Free()`** | Correct teardown; avoids leaks. |
| **`Entity.Copy()`** for bullets | Avoids reloading the same `.glb` every shot when the prefab has a **load path**. |
| **`Pool.*()`** | Reuse **heap** objects you control (factory returns handles). |
| **`Entity.Shoot()`** | One call for clone → physics → velocity → timed purge. |

---

## Power-user getters (today)

| Need | Command | Notes |
|------|---------|--------|
| Position | `Entity.X()` / `Entity.Y()` / `Entity.Z()` | Shorthand macros. |
| Rotation | `Entity.Pitch()` / `Entity.Yaw()` / `Entity.Roll()` | |
| Linear velocity | `Entity.GetLinearVel()` | Reads **Jolt** when physics is active. |
| Cap speed | `Entity.LimitSpeed(id, max)` | Clamps vector length. |

---

## See also

- [GAMEHELPERS.md](GAMEHELPERS.md) — bridge API overview  
- [POOL.md](POOL.md) — `Pool.*()` contract  
- [MEMORY.md](../MEMORY.md) — three-layer memory model  
