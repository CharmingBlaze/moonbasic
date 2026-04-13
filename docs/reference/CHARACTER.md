# Character Commands (`CHARACTER.*` / `Character.*`)

High-level **heap-backed character** API: **`CHARACTER.CREATE`** returns a **`Character`** handle (`CHARACTERREF.*`) for kinematic movement, independent of the **`PLAYER.*`** / **`CHAR.MAKE`** entity-id style when you want a **facade** that feels like a single object.

**Platform note (read first):** **Polymorphic** **`CHARACTER.CREATE`** — **standalone** **`(x#, y#, z#)`** — is implemented on the **host KCC** module (**`api_host`**, e.g. **Windows** `fullruntime`). **Linux + CGO + Jolt** currently wires **`CHARACTER.CREATE`** only as **`(entity, radius#, height#)`** (see [`character_ref_linux.go`](../../runtime/player/character_ref_linux.go)); for a Jolt capsule at a world position without this split, use **`CHAR.MAKE` / `PLAYER.CREATE`** on a positioned entity or **`CHARCONTROLLER.MAKE`**.

For **`PLAYER.*`** gameplay (look targets, swim, nav), see [PLAYER.md](PLAYER.md). For low-level **`CHARCONTROLLER.*`** handles, see [CHARCONTROLLER.md](CHARCONTROLLER.md). For **`CHAR.*`** / KCC tuning, see [KCC.md](KCC.md).

**Documentation order:** [Platform priority](../DEVELOPER.md#platform-priority-windows-then-linux) — Windows-first where OSes differ.

---

## Polymorphic creation (`CHARACTER.CREATE`)

The engine chooses the mode from the **first argument**:

| Mode | Signature | When it applies |
|------|-----------|-------------------|
| **Standalone** | **`CHARACTER.CREATE(x#, y#, z#)`** | First argument is **not** an entity handle **and** does not resolve to a positive entity id — three world coordinates spawn a **physics-only** character at **`(x, y, z)`** (host KCC). |
| **Entity-bound** | **`CHARACTER.CREATE(entity, radius#, height#)`** | First argument is an **EntityRef** or **entity id** — KCC is bound to that **visual entity**; **`radius`** and **`height`** define the capsule (same idea as **`CHAR.MAKE`**). |

Both forms return a **handle** to a heap **`Character`** object; use **`CHARACTERREF.*`** methods on that handle.

### Standalone — `Character.Create(x#, y#, z#)`

- Spawns a kinematic character **without** requiring a **`MODEL.*`** / **EntityRef**.
- The runtime assigns a **virtual id**: a **negative** **`int64`** taken from a descending counter (initialized at **`-1000`**, then **`-1001`**, …). That id keys **host KCC** state so the **solver** can run **without** a scene entity row, while scripts still pass **numeric** ids into helpers that expect “entity-like” indices where wired.
- **VM note:** Values live in the normal **register / value** model (64-bit integer slots for ids); virtual ids are **not** heap handles themselves — the **returned Character handle** is the **`CHARACTERREF`** target for **`CHARACTERREF.UPDATE`**, **`CHARACTERREF.JUMP`**, etc.

### Entity-bound — `Character.Create(entity, radius#, height#)`

- Same capsule semantics as **`CHAR.MAKE(entity, radius#, height#)`** on the host: scripted physics on the entity is cleared and **KCC** drives motion.
- On **Linux + Jolt**, this path allocates **Jolt `CharacterVirtual`** and syncs the entity transform from the capsule.

---

## Core workflow (host / standalone-friendly)

1. **`World.Setup()`** (Recommended) or **`PHYSICS3D.START()`** and set gravity.
2. **`hero = Character.Create(0, 5, 0)`** *or* **`hero = Character.Create(playerEnt, 0.4, 1.0)`**.
3. (Optional) Tune physics: **`CharacterRef.SetPadding(hero, 0.02)`**, **`CharacterRef.SetFriction(hero, 0.9)`**.
4. Each frame: **`CHARACTERREF.UPDATE(hero, dt)`** (or your game loop’s **`UPDATEPHYSICS`** bundle), input → **`CHARACTERREF.MOVEWITHCAMERA`**, **`CHARACTERREF.JUMP`**, etc.
4. **Standalone:** sync any **visual** (e.g. separate **`MODEL`**) by reading **`CHARACTERREF.GETPOSITION`** and moving the mesh, or attach gameplay to the virtual id where the engine exposes it.

---

## `CHARACTERREF.*` (handle receiver)

See manifest entries under **`CHARACTERREF.*`** in [API_CONSISTENCY.md](../API_CONSISTENCY.md). Typical calls:

* **`CHARACTERREF.UPDATE(handle, dt#)`**: Advances simulation.
* **`CHARACTERREF.SETFRICTION(handle, friction#)`**: Sets sliding resistance (0..1).
* **`CHARACTERREF.SETBOUNCE(handle, bounciness#)`**: Sets restitution (0..1).
* **`CHARACTERREF.SETPADDING(handle, padding#)`**: Sets collision margin (default 0.02).
* **`CHARACTERREF.JUMP(handle, force#)`**: Applies upward impulse.
* **`CHARACTERREF.MOVEWITHCAMERA(handle, forward, right, ...)`**: Smart movement relative to view.
* **`CHARACTERREF.FREE(handle)`**: Releases resources (Jolt/Host).

---

## Examples

### Standalone spawn (host KCC)

```moonbasic
PHYSICS3D.START()
WORLD.GRAVITY(0, -40, 0)

hero = Character.Create(0, 5, 0)

WHILE NOT Window.ShouldClose()
    dt = TIME.DELTA()
    CHARACTERREF.UPDATE(hero, dt)
    ; ... input, movement on hero ...
    RENDER.FRAME()
WEND
```

### Entity-bound spawn

```moonbasic
PHYSICS3D.START()
playerEnt = MODEL.CREATECAPSULE(0.4, 1.0)
playerEnt.Pos(0, 5, 0)

hero = Character.Create(playerEnt, 0.4, 1.0)

WHILE NOT Window.ShouldClose()
    CHARACTERREF.UPDATE(hero, TIME.DELTA())
    RENDER.FRAME()
WEND
```

---

## See also

- [KCC.md](KCC.md) — **`CHAR.*`** / **`PLAYER.*`** gameplay layer  
- [PLAYER.md](PLAYER.md) — **`PLAYER.CREATE`**, queries, swim  
- [CHARCONTROLLER.md](CHARCONTROLLER.md) — **`CHARCONTROLLER.MAKE`** capsule API  
- [PHYSICS3D.md](PHYSICS3D.md) — world step and picks  
- [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) — reference formatting  
