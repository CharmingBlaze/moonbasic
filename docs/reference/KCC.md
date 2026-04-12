# Kinematic Character Controller (`CHAR.*` / `PLAYER.*`)

This page is the **gameplay-first** guide to MoonBASIC’s **Kinematic Character Controller (KCC)** — Jolt **`CharacterVirtual`** behind **`PLAYER.CREATE`**, **`CHAR.MAKE`**, and related commands. It bridges **“what I want the hero to do”** and **stable 3D navigation** (wall slide, stairs, floor stick) without you hand-writing collision response.

For the **low-level capsule API** (`CharController.*` handles), see [CHARCONTROLLER.md](CHARCONTROLLER.md). For the full **`PLAYER.*`** surface (swim, push, surface type, …), see [PLAYER.md](PLAYER.md). For **heap `Character` handles** — including **polymorphic** **`CHARACTER.CREATE`** (**`Character.Create(x, y, z)`** vs **`Character.Create(entity, r, h)`**) and **virtual ids** on the host solver — see [CHARACTER.md](CHARACTER.md).

## Platform

Project policy: document **Windows** first, **Linux** second ([DEVELOPER.md](../DEVELOPER.md#platform-priority-windows-then-linux)).

| | Windows (`fullruntime`, CGO + Raylib) | Linux + CGO + Jolt (`fullruntime`) |
|--|----------------------------------------|-------------------------------------|
| **`CHAR.*` / `NAV.*` / `PLAYER.NAV*` (KCC)** | **Host KCC** (analytic floor + static AABB; no Jolt). Same commands; stairs/slope parity is approximate. | Full **Jolt CharacterVirtual** |
| **`ENT.*`**, **`WORLD.TOSCREEN`**, **`WORLD.HITSTOP`**, **`ENT.SHOOT`**, **`ENT.FADE`** | Yes — **entity** and **time** helpers work wherever **`mbentity`** + Raylib run. | Yes |
| **`WORLD.MOUSEFLOOR` / `WORLD.MOUSEPICK`** | **Stub** returns errors without native Jolt (see [PHYSICS3D.md](PHYSICS3D.md)). | Needs Jolt picks |

Start the world with **`PHYSICS3D.START()`** (and set gravity) **before** **`CHAR.MAKE` / `PLAYER.CREATE`**.

**Entity handles:** Pass the **EntityRef** from **`MODEL.CREATECAPSULE`** (or **`CUBE`** / **`SPHERE`**) into **`CHAR.MAKE`**, **`CHAR.MOVEWITHCAMERA`**, **`CHAR.JUMP`**, etc. The runtime resolves the handle to an internal entity id; using a wrong integer (e.g. the heap slot) breaks KCC lookup.

---

## 1. Setup and tuning (“pro” feel)

| Command | Role | Beginner | Advanced |
|--------|------|----------|----------|
| **`CHAR.MAKE(entity)`** or **`CHAR.MAKE(entity, radius#, height#)`** | Create **`CharacterVirtual`** at the entity’s position and map **entity → controller**. Clears scripted gravity/velocity so the capsule drives motion. | `CHAR.MAKE(hero)` | Match mesh: `CHAR.MAKE(hero, 0.4, 1.0)` |
| **`CHAR.SETSTEP(entity, height#)`** | Max **step up** (stairs / curbs), via Jolt extended update — always **`(entity, height)`**, not a lone height. | `CHAR.SETSTEP(hero, 0.3)` | Tune per level art |
| **`CHAR.SETSLOPE(entity, degrees#)`** | **`PLAYER.SETSLOPELIMIT`** — rebuilds capsule with **`MaxSlopeAngle`**. | `CHAR.SETSLOPE(hero, 45)` | Lower to block “walking up walls” |
| **`CHAR.SETPADDING(entity, padding#)`** | Skin around capsule (**&gt; 0**); reduces snagging on messy geometry. | Often omit (runtime default) | `CHAR.SETPADDING(hero, 0.02)` |

Aliases: **`PLAYER.CREATE`** = **`CHAR.MAKE`**; **`PLAYER.SETSTEPOFFSET`** = **`CHAR.SETSTEP`**; **`PLAYER.SETSLOPELIMIT`** = **`CHAR.SETSLOPE`**; **`PLAYER.SETPADDING`** = **`CHAR.SETPADDING`**.

**Important:** Do **not** put the hero on **`ENTITY.PHYSICS`** as a **dynamic** body if you are using **`CHAR.MAKE`** — the KCC owns movement and collision for that entity. Keep **static** level meshes as usual.

### Capsule size and pivot (primitive or glTF hero)

- **`MODEL.CREATECAPSULE(radius#, height#)`** draws a **Jolt-style** capsule: pivot at the **center** of the shape; total height is **`height#`** (same convention as **`CHAR.MAKE(…, radius#, height#)`**).
- Use the **same** `radius` and `height` in **`CHAR.MAKE(hero, radius, height)`** as in **`MODEL.CREATECAPSULE`**, or feet vs floor will not match the mesh. For an imported **`MODEL.LOAD` / glTF** hero, pick **`radius` / `height`** that match your collision need; the host KCC on Windows uses **height/2** from the pivot down to the feet (center-pivot capsules), not the radius, for ground contact.
- Arbitrary meshes still **render** as authored; KCC uses the **numeric capsule** you pass — it does not auto-read mesh bounds yet.

---

## 2. Movement (no basis-vector math)

| Command | Role |
|--------|------|
| **`CHAR.MOVE(entity, dirX#, dirZ#, speed#)`** | **Intent direction × speed** in **world XZ** (e.g. **-1…1** from input). Slides on walls. |
| **`CHAR.MOVEWITHCAMERA`** / **`CHAR.MOVEWITHCAM`** | Same — camera-relative **WASD** on **XZ** (`CameraXZWalkBasis`). |
| **`NAV.GOTO`** | Alias of **`PLAYER.NAVTO`** — click-to-move; default **arrival** radius is **~0.2** world units so the hero **stops cleanly** (no jitter at the exact point). Optional **`arrivalXZ`** / **`brakeDist`** match **`PLAYER.NAVTO`**. |
| **`NAV.UPDATE`** | Alias of **`PLAYER.NAVUPDATE`** — call each frame while navigating. |
| **`NAV.CHASE(entity, target#, gap#, speed#)`** | **KCC only:** move toward **target** entity until within **gap** (world units), then hold. |
| **`NAV.PATROL(entity, ax#, az#, bx#, bz#, speed#)`** | **KCC only:** ping-pong between world **XZ** points **A** and **B** (same idea as **`ENTITY.PATROL`**, but for **`CHAR.MAKE`** entities). |
| **`PLAYER.NAVTO(entity, tx#, tz#, speed# [, arrivalXZ# [, brakeDist#]])`** | Same as **`NAV.GOTO`** — explicit **`PLAYER.*`** name. |
| **`PLAYER.NAVUPDATE(entity)`** | Steps navigation toward the active target (goto / chase / patrol) with **soft braking** when **`brakeDist`** is set. |
| **`CHAR.JUMP(entity, impulseY#)`** | Vertical impulse (snappy hop; not “physics toy” bounciness). |
| **`CHAR.STICK(entity, dist#)`** | Alias of **`PLAYER.SETSTICKFLOOR`** — **stick-to-floor** max step **down** (“glue” so stairs don’t feel like flying). |

Lower-level: **`PLAYER.MOVE(entity, vx#, vz#)`** is **world velocity** (units/sec), not **`dir * speed`**.

---

## 3. World awareness (mouse and queries)

| Command | Role |
|--------|------|
| **`WORLD.MOUSEFLOOR(camera, floorY#)`** | Alias of **`WORLD.MOUSEFLOOR3D`** — mouse ray vs plane **y = floorY** → **`[wx, wz]`** handle or **NIL**. |
| **`WORLD.MOUSEPICK(camera)`** | Alias of **`WORLD.MOUSETOENTITY`** — **entity id** under the cursor (ray into the physics world). |
| **`WORLD.TOSCREEN(wx#, wy#, wz#)`** | Active 3D camera (**after `CAMERA.BEGIN`**) → **`[sx, sy]`** pixel handle. |
| **`WORLD.TOSCREEN(entity#)`** | Same, using the entity’s **world position** (handy for HUD / health bars). |
| **`WORLD.HITSTOP(duration#)`** | **Gameplay freeze** for **wall-clock** seconds (uses **`HitStopEndAt`** + **`TIME.DELTA`/`DT` → 0**); impact-frame “crunch”. |
| **`CHAR.ISGROUNDED(entity)`** | **`TRUE`** if Jolt reports ground support. |
| **`CHAR.ISGROUNDED(entity, coyoteSec#)`** | Optional **coyote time**: still **`TRUE`** for **coyoteSec** seconds after the last grounded frame (wall-clock), so **`IF CHAR.ISGROUNDED(hero, 0.12)`** feels forgiving. **Linux+Jolt KCC only** (see stub table above). |
| **`CHAR.DIST(a, b)`** / **`ENTITY.DIST` / `ENT.DIST`** | **3D distance** between two entities (same implementation as **`ENTITY.DISTANCE`**). |

---

## 4. RPG helpers (health, tags, teams)

| Command | Role |
|--------|------|
| **`ENT.SET_HP` / `ENTITY.SETHEALTH`** | Init **HP** (and optional max). |
| **`ENT.DAMAGE` / `ENTITY.DAMAGE`** | Apply damage; triggers a **~0.1s red tint** on **`r,g,b`** (restored automatically — requires **`ENTITY.UPDATE`** / normal frame loop). |
| **`ENT.ISALIVE` / `ENTITY.ISALIVE`** | **FALSE** when HP ≤ 0. |
| **`ENT.SET_TEAM` / `ENT.SETTEAM`** | Store a **team id** for gameplay (friendly fire, AI factions). |
| **`ENT.SET_HP` / `ENT.SETHP`** | Aliases of **`ENTITY.SETHEALTH`**. |
| **`ENT.GET_NEAREST` / `ENT.GETNEAREST` / `PLAYER.GETNEARBY`** | Entities in **radius** with a **name/tag** match → **float array** of ids. |
| **`ENT.ONDEATH(entity, prefab)`** | **100%** death drop — **`prefab`** may be an **entity id** or a **name string** resolved via **`ENTITY.SETNAME`** / **`byName`** (same as **`ENTITY.FIND`**). Full control: **`ENTITY.ONDEATHDROP`**. |

---

## 5. Polish (juice)

| Command | Role |
|--------|------|
| **`ENT.WOBBLE`** | Alias of **`ENTITY.ADDWOBBLE`** — floating pickup motion (**height**, **speed**). |
| **`ENT.TWEEN`** | Alias of **`ENTITY.ANIMATETOWARD`** — smooth move to **world (x, y, z)** over **duration**. |
| **`ENT.FADE(entity, alpha#, duration#)`** | Convenience fade to **alpha** over **duration** (implemented via **`ENTITY.FADE`**). |
| **`ENT.SHOOT(shooter, prefab, speed#)`** | **`prefab`** = entity **or** registered **name** string. **`ENTITY.COPY`**, place forward, set **host velocity** — runs on **Windows + Linux** with **`mbentity`**; add Jolt **`BODY3D`** yourself if you need CCD. |
| **`WORLD.SHAKE` / `WORLD.SCREENSHAKE`** | Screen impact on the **active camera**. |

For property tweens (**alpha**, **yaw**, …) use **`ENTITY.TWEEN`** (different overload — see [ENTITY.md](ENTITY.md)).

---

## Minimal click-to-move sketch

```moonbasic
PHYSICS3D.START()
hero = Entity.Load("hero.iqm")
CHAR.MAKE(hero, 0.4, 1.75)

pt = WORLD.MOUSEFLOOR(cam, 0.0)
IF pt <> NIL THEN
    wx = Array.Get(pt, 0)
    wz = Array.Get(pt, 1)
    PLAYER.NAVTO(hero, wx, wz, 6.0, 0.2, 0.75)
ENDIF

WHILE ...
    PLAYER.NAVUPDATE(hero)
    UPDATEPHYSICS()
WEND
```

(Adjust **`Array.Get`** to your project’s float-array access pattern.)

---

## Example

See **`examples/mario64/modern_blitz_hop_kcc.mb`** — orbit camera + **`CHAR.MOVEWITHCAMERA`** + **`CHAR.JUMP`**.

## Implementation notes (engine behavior)

- **`NAV.GOTO` / `PLAYER.NAVTO`**: default **arrival** is **0.2** world units; when inside that radius, the runtime applies **zero horizontal velocity** on the **`CharacterVirtual`** so the capsule does not **overshoot and jitter** at the click point. **Soft braking** still uses **`brakeDist`** (quadratic ease).
- **`CHAR.STICK` / `PLAYER.SETSTICKFLOOR`**: maps to Jolt **ExtendedUpdate** **`StickToFloorStepDown`** (see `SetCharacterStickToFloorDown` in the Linux charcontroller).
- **`ENT.DAMAGE`**: **0.1s** material tint (red) then restore — no separate shader; tint is on entity **RGB** fields.

## See also

- [PLAYER.md](PLAYER.md) — full **`PLAYER.*`** list  
- [PHYSICS3D.md](PHYSICS3D.md) — world, **`PICK.*`**, layers  
- [PHYSICS_ERGONOMICS.md](../PHYSICS_ERGONOMICS.md) — when to use KCC vs dynamic **`ENTITY.PHYSICS`**
