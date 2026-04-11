# Player interaction (`PLAYER.*`)

High-level helpers for **kinematic character control** (KCC) and **spatial queries** against level geometry. These wrap Jolt’s **`CharacterVirtual`** (same subsystem as **`CHARCONTROLLER.*`**) and entity/tag data from **`LEVEL.LOAD`**.

## Platform

| Feature | Linux + CGO + Jolt | Windows / no physics |
|--------|---------------------|----------------------|
| **`PLAYER.CREATE` / `MOVE` / `JUMP` / `ISGROUNDED` / `SYNCANIM`** | Supported | Error (use **`CHARCONTROLLER.*`** manually if you add Windows Jolt later) |
| **`PLAYER.GETLOOKTARGET` / `GETNEARBY` / `SETSTATE`** | Full | Works (look uses physics ray when **`PHYSICS3D.START`**; otherwise mesh AABB fallback) |

Start the world with **`PHYSICS3D.START()`** before **`PLAYER.CREATE`**.

---

## Kinematic character (Linux + Jolt)

| Command | Purpose |
|--------|---------|
| **`PLAYER.CREATE(entity)`** | Spawns a **capsule** character controller at the entity’s world position (default radius **0.4**, height **1.75**). Stores a mapping **entity → controller** for subsequent **`PLAYER.*`** calls. |
| **`PLAYER.MOVE(entity, velocityX, velocityZ)`** | World-space **horizontal** velocity in **units per second** (multiplied by **`TIME.DELTA`** internally). Uses **`CharacterMoveXZVelocity`** (slide/step via Jolt **`ExtendedUpdate`**). Syncs the **entity** transform to the capsule after the move. |
| **`PLAYER.JUMP(entity, impulseY)`** | Adds **impulseY** to upward linear velocity (same idea as **`CharacterJump`**). |
| **`PLAYER.ISGROUNDED(entity)`** → **bool** | **`true`** if the Jolt character reports ground support (**`IsSupported`**). |

Lower-level access without entity ids: **`CHARCONTROLLER.MAKE` / `MOVE` / …** ([CHARCONTROLLER.md](CHARCONTROLLER.md)).

---

## Interaction & detection

| Command | Purpose |
|--------|---------|
| **`PLAYER.GETLOOKTARGET(entity, maxDist)`** → **entity** | **Eye height ≈ 1.65** above feet. Casts a **physics ray** along the entity’s **world forward** (`PickCastEntityID`). If the first hit is the player or nothing, falls back to **`ENTITY.PICK`-style AABB** ray vs **static** entities. Returns **0** if none. |
| **`PLAYER.GETNEARBY(entity, radius, tag)`** → **float array** | Entities within **radius** whose **`ENTITY`** name **or** Blender **`tag`** extra matches **`tag`** (case-insensitive **`path.Match`** glob, e.g. **`Enemy*`**). Returns a **numeric array** of entity ids (same pattern as other “list of ids” APIs). |

---

## Triggers & animation

| Command | Purpose |
|--------|---------|
| **`PLAYER.ONTRIGGER(entity, callbackFunc)`** | **Not implemented** — the VM cannot be entered from Jolt sensors yet. Use **`LEVEL.BINDSCRIPT`** + **`LEVEL.MATCHSCRIPTBIND`**, **`EntityCollided`**, or **`PHYSICS3D`** collision hooks instead. A future **physics → BASIC** callback path would need a strict **main-thread / post-step** queue (no VM reentrancy inside Jolt) and is separate from the **wazero** WASM story. |
| **`PLAYER.SETSTATE(entity, state)`** | Stores an integer **state id** for gameplay logic (e.g. **0 = idle**, **1 = walk**, **2 = jump**). Constants **`STATE_*`** are not built-ins yet—use literals or your own **`CONST`**. |
| **`PLAYER.SYNCANIM(entity [, scale])`** | Sets **`ENTITY`** animation speed from **horizontal** linear velocity (× optional **scale**, default **1**). Requires **`PLAYER.CREATE`**. |
| **`PLAYER.SETSTEPHEIGHT(entity, height)`** | Records a desired max **step/stair** height for tooling and future Jolt tuning (runtime step height is not exposed in the current **`jolt-go`** character wrapper). |
| **`PLAYER.SETSLOPELIMIT(entity, maxSlopeDegrees)`** | **Rebuilds** the **`CharacterVirtual`** with **`MaxSlopeAngle`** = **maxSlopeDegrees** (must be between **0** and **90**). Preserves linear velocity. |
| **`PLAYER.GETVELOCITY(entity)`** → **vec3 handle** | **`CharacterVirtual`** linear velocity (**vx, vy, vz**). |
| **`PLAYER.TELEPORT(entity, x, y, z)`** | **`SetPosition`** + clears velocity + **`ExtendedUpdate`** + syncs the **entity** transform (snap teleport without smoothing). |
| **`PLAYER.SETGRAVITYSCALE(entity, scale)`** | Scales **gravity on Y** during **`CharacterMoveXZVelocity`** (**1** = default; values below **1** lighten gravity; above **1** strengthen it). |
| **`PLAYER.GETCROUCH(entity)`** / **`PLAYER.SETCROUCH(entity, bool)`** | Stored **crouch** flag for gameplay. **Capsule height** is not changed yet (Jolt wrapper limitation). |
| **`PLAYER.SWIM(entity, buoyancy, drag)`** | **Swim mode**: **buoyancy** (0–1) reduces downward gravity; **drag** damps horizontal velocity per second. Use **`(0, 0)`** to disable. |
| **`PLAYER.SETSTEPOFFSET(entity, height)`** | Alias of **`PLAYER.SETSTEPHEIGHT`** (reserved for future stair tuning). |
| **`PLAYER.GETSTANDNORMAL(entity)`** → **vec3 handle** | Ground/floor normal under the feet (**`GetGroundNormal`** or short downward ray). |
| **`PLAYER.PUSH(player, target, force)`** | Forward **horizontal** push on **target** via host **`ENTITY.ADDFORCE`**-style integration; scaled by **`PLAYER.SETMASS`**. |
| **`PLAYER.GRAB(player, target)`** | Each **`PLAYER.MOVE`**, repositions **target** in front of the player ( **`target 0`** releases). Not a Jolt **fixed constraint** yet. |
| **`PLAYER.SETMASS(entity, mass)`** | Stores **gameplay mass** (e.g. **`PLAYER.PUSH`**); Jolt **CharacterVirtual** mass is fixed at **`PLAYER.CREATE`**. |
| **`PLAYER.GETSURFACETYPE(entity)`** → **string** | Downward **Jolt** ray → hit entity → **`SurfaceMaterialHint`** from glTF **`material` / `footstep`** metadata or **Blender tag**; else **`Default`**. |
| **`PLAYER.SETFOVKICK` / `PLAYER.GETFOVKICK`** | Stores **extra FOV degrees** per entity; each frame do **`Camera.SetFOV(cam, base + Player.GetFovKick(hero))`** (or your own base). |
| **`PLAYER.ISMOVING(entity)`** → **bool** | **True** if horizontal **linear speed** is above ~**0.05** (for footsteps / sprint FX). |
| **`PLAYER.SnapToGround(entity, terrain, offset)`** | Sets **Y** from **`Terrain.GetHeight`** at the entity’s **XZ** plus **offset** (feet vs pivot). On **Linux + Jolt** after **`PLAYER.CREATE`**, also syncs the **CharacterVirtual** capsule. |
| **`PLAYER.ISSWIMMING(entity)`** → **bool** | **True** when the entity’s position lies inside a **`WATER`** volume column (between **bed** and the wavy surface). Use with **`PLAYER.SETGRAVITYSCALE`** for floatier movement. |

---

## Example (Linux)

```moonbasic
Physics3D.Start()
hero = Entity.Load("hero.iqm")
Player.Create(hero)

WHILE Window.Open()
    dt = Time.Delta()
    Player.Move(hero, Input.AxisX() * 5.0, Input.AxisY() * 5.0)
    IF Player.IsGrounded(hero) AND Input.KeyPressed(KEY_SPACE) THEN
        Player.Jump(hero, 6.0)
    ENDIF
    target = Player.GetLookTarget(hero, 3.0)
    IF target <> 0 AND Input.KeyPressed(KEY_E) THEN
        fn = LEVEL.MATCHSCRIPTBIND(EntityName(target))
        REM dispatch fn in your script...
    ENDIF
    Player.SyncAnim(hero, 0.12)
    Begin3D()
        Entity.Draw(hero)
    End3D()
WEND
```

Naming: use **`LEVEL.LOAD`** / **`Entity.Draw`** (or your project’s draw path), not **`Scene.Draw`**, so **`SCENE.*`** stays reserved for **mbscene** game scenes.

---

## See also

- [CHARCONTROLLER.md](CHARCONTROLLER.md) — capsule API and full sample  
- [LEVEL.md](LEVEL.md) — glTF, tags, **`LEVEL.BINDSCRIPT`**  
- [PHYSICS3D.md](PHYSICS3D.md) — Jolt world, **`PICK.*`**, rays  
