# Character Commands (`CHARACTER.*` / `Character.*`)

High-level **heap-backed character** API: **`CHARACTER.CREATE`** returns a **`Character`** handle (`CHARACTERREF.*`) for kinematic movement, independent of the **`PLAYER.*`** / **`CHAR.CREATE`** entity-id style (deprecated **`CHAR.MAKE`**) when you want a **facade** that feels like a single object.

**Desktop (Windows/Linux, `fullruntime`, CGO + Jolt):** **`CHARACTER.CREATE`** is **`(entity, radius#, height#)`** only — it allocates **Jolt `CharacterVirtual`** for that entity (see [`character_ref_cgo.go`](../../runtime/player/character_ref_cgo.go)). There is no standalone **`(x, y, z)`** spawn; use **`CHAR.CREATE` / `PLAYER.CREATE`** (deprecated **`CHAR.MAKE`**) on a positioned entity (e.g. **`MODEL.CREATECAPSULE`**) or **`CHARCONTROLLER.CREATE`** (deprecated **`CHARCONTROLLER.MAKE`**) for low-level handles.

For **`PLAYER.*`** gameplay (look targets, swim, nav), see [PLAYER.md](PLAYER.md). For low-level **`CHARCONTROLLER.*`** handles, see [CHARCONTROLLER.md](CHARCONTROLLER.md). For **`CHAR.*`** / KCC tuning, see [KCC.md](KCC.md).

**Documentation order:** [Platform priority](../DEVELOPER.md#platform-priority-windows-then-linux) — Windows-first where OSes differ.

---

## Creation (`CHARACTER.CREATE`)

| Mode | Signature |
|------|-----------|
| **Entity-bound** | **`CHARACTER.CREATE(entity, radius#, height#)`** — KCC is bound to that **visual entity**; **`radius`** and **`height`** define the capsule (same idea as **`CHAR.CREATE`** / deprecated **`CHAR.MAKE`**). |

Returns a **handle** to a heap **`Character`** object; use **`CHARACTERREF.*`** methods on that handle.

### Entity-bound — `Character.Create(entity, radius#, height#)`

- Same capsule semantics as **`CHAR.CREATE(entity, radius#, height#)`**: scripted physics on the entity is cleared and **Jolt CharacterVirtual** drives motion on **Linux and Windows** when built with **CGO + Jolt**.

---

## Core workflow

1. **`World.Setup()`** (Recommended) or **`PHYSICS3D.START()`** and set gravity.
2. **`hero = Character.Create(playerEnt, 0.4, 1.0)`** (entity must exist, e.g. **`MODEL.CREATECAPSULE`**).
3. (Optional) Tune physics: **`CharacterRef.SetPadding(hero, 0.02)`**, **`CharacterRef.SetFriction(hero, 0.9)`**.
4. Each frame: **`CHARACTERREF.UPDATE(hero, dt)`** (or your game loop’s **`UPDATEPHYSICS`** bundle), input → **`CHARACTERREF.MOVEWITHCAMERA`**, **`CHARACTERREF.JUMP`**, etc.

---

## `CHARACTERREF.*` (handle receiver)

See manifest entries under **`CHARACTERREF.*`** in [API_CONSISTENCY.md](../API_CONSISTENCY.md). Typical calls:

* **`CHARACTERREF.UPDATE(handle, dt#)`**: Advances simulation.
* **`CHARACTERREF.SETPOS(handle, x#, y#, z#)`**: Teleports the capsule (canonical); **`CHARACTERREF.SETPOSITION`** is a deprecated alias.
* **`CHARACTERREF.SETFRICTION(handle, friction#)`**: Sets sliding resistance (0..1).
* **`CHARACTERREF.SETBOUNCE(handle, bounciness#)`**: Sets restitution (0..1).
* **`CHARACTERREF.SETPADDING(handle, padding#)`**: Sets collision margin (default 0.02).
* **`CHARACTERREF.JUMP(handle, force#)`**: Applies upward impulse.
* **`CHARACTERREF.MOVEWITHCAMERA(handle, forward, right, ...)`**: Smart movement relative to view.
* **`CHARACTERREF.FREE(handle)`**: Releases resources (Jolt character via charcontroller).

---

## Examples

### Entity-bound spawn

```moonbasic
PHYSICS3D.START()
playerEnt = MODEL.CREATECAPSULE(0.4, 1.0)
playerEnt.Pos(0, 5, 0)

hero = Character.Create(playerEnt, 0.4, 1.0)

WHILE NOT WINDOW.SHOULDCLOSE()
    CHARACTERREF.UPDATE(hero, TIME.DELTA())
    RENDER.FRAME()
WEND
```

---

## See also

- [KCC.md](KCC.md) — **`CHAR.*`** / **`PLAYER.*`** gameplay layer  
- [PLAYER.md](PLAYER.md) — **`PLAYER.CREATE`**, queries, swim  
- [CHARCONTROLLER.md](CHARCONTROLLER.md) — **`CHARCONTROLLER.CREATE`** capsule API (deprecated **`CHARCONTROLLER.MAKE`**)  
- [PHYSICS3D.md](PHYSICS3D.md) — world step and picks  
- [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) — reference formatting  
