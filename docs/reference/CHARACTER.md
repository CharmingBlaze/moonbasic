# Character Commands (`CHARACTER.*` / `CHARACTERREF.*`)

High-level **heap-backed character controller** API. **`CHARACTER.CREATE`** binds a **Jolt `CharacterVirtual`** capsule to a visual entity and returns a **`CHARACTERREF`** handle for all subsequent motion and query calls.

**Platform:** Windows and Linux with **CGO + Jolt** (`fullruntime`). Same keys are registered as stubs on other builds. **Documentation order:** Windows-first ([DEVELOPER.md](../DEVELOPER.md#platform-priority-windows-then-linux)).

For **`PLAYER.*`** high-level gameplay, see [PLAYER.md](PLAYER.md). For the standalone capsule API, see [CHARCONTROLLER.md](CHARCONTROLLER.md). For KCC tuning, see [KCC.md](KCC.md).

## Core Workflow

1. **`PHYSICS3D.START()`** (or **`World.Setup()`**) and set gravity.
2. Create a visual entity, e.g. **`MODEL.CREATECAPSULE(radius, height)`**, position it.
3. **`hero = CHARACTER.CREATE(entity, radius, height)`** — binds Jolt capsule to the entity.
4. (Optional) tune: `hero.setPadding(0.02)`, `hero.setFriction(0.9)`, `hero.setGravityScale(1.0)`.
5. Each frame: **`CHARACTERREF.UPDATE(hero, dt)`**, then input → **`CHARACTERREF.MOVEWITHCAMERA`** / **`CHARACTERREF.JUMP`**.
6. **`hero.free()`** when done.

### Method chaining 

All `CHARACTERREF.*` mutating calls return the handle: `hero.setPos(x,y,z).jump(force)`.

---

## Creation

### `CHARACTER.CREATE(entity, radius, height)` 

Binds a Jolt `CharacterVirtual` capsule to `entity`. Clears any scripted physics on the entity. Returns a **`CHARACTERREF` handle**. Aliases: `CHAR.CREATE` (preferred entity-id style); deprecated `CHAR.MAKE`.

---

## Simulation

### `CHARACTERREF.UPDATE(hero, dt)` 

Advances the character simulation by `dt` seconds. Call once per frame before reading state.

- *Handle shortcut*: `hero.update(dt)`

---

### `CHARACTERREF.UPDATEMOVE(hero, ...)` 

Low-level `UpdateGroundVelocity` + move cycle. Used when you need finer control than `UPDATE`. See runtime source for argument details.

- *Handle shortcut*: `hero.updateMove(...)`

---

## Position & Velocity

### `CHARACTERREF.SETPOS(hero, x, y, z)` 

Teleports the capsule to world position. Alias: `CHARACTERREF.SETPOSITION` (same function).

- *Handle shortcut*: `hero.setPos(x, y, z)`

---

### `CHARACTERREF.GETPOSITION(hero)` 

Returns current world position as a `[x, y, z]` array handle.

- *Handle shortcut*: `hero.getPosition()`

---

### `CHARACTERREF.SETVELOCITY(hero, vx, vy, vz)` 

Sets world linear velocity directly. Alias: `CHARACTERREF.SETLINEARVELOCITY`.

- *Handle shortcut*: `hero.setVelocity(vx, vy, vz)`

---

### `CHARACTERREF.GETVELOCITY(hero)` 

Returns current linear velocity as a `[vx, vy, vz]` array handle.

- *Handle shortcut*: `hero.getVelocity()`

---

## Motion

### `CHARACTERREF.MOVE(hero, dx, dy, dz)` 

Applies a displacement this frame; collisions resolved via Jolt extended update.

- *Handle shortcut*: `hero.move(dx, dy, dz)`

---

### `CHARACTERREF.MOVEWITHCAMERA(hero, ...)` 

Smart movement relative to the active camera view. See `commands.json` for arity.

- *Handle shortcut*: `hero.moveWithCamera(...)`

---

### `CHARACTERREF.JUMP(hero, force)` 

Applies an upward impulse of `force`.

- *Handle shortcut*: `hero.jump(force)`

---

## Physics Tuning

### `CHARACTERREF.SETGRAVITY(hero, g)` 

Sets the effective gravity magnitude for this character.

- *Handle shortcut*: `hero.setGravity(g)`

---

### `CHARACTERREF.SETGRAVITYSCALE(hero, scale)` 

Scales world gravity for this character (e.g. `0.5` = half gravity).

- *Handle shortcut*: `hero.setGravityScale(scale)`

---

### `CHARACTERREF.SETSTICKDOWN(hero, dist)` 

Sets the snap-to-ground stick distance (keeps character grounded on shallow slopes/steps).

- *Handle shortcut*: `hero.setStickDown(dist)`

---

### `CHARACTERREF.SETFRICTION(hero, f)` 

Sets sliding resistance (0..1).

- *Handle shortcut*: `hero.setFriction(f)`

---

### `CHARACTERREF.SETBOUNCE(hero, b)` 

Sets restitution coefficient (0..1). Alias: `CHARACTERREF.SETBOUNCINESS`.

- *Handle shortcut*: `hero.setBounce(b)`

---

### `CHARACTERREF.SETPADDING(hero, p)` 

Sets the collision contact margin (default `0.02`).

- *Handle shortcut*: `hero.setPadding(p)`

---

### `CHARACTERREF.SETMAXSLOPE(hero, angle)` 

Sets the maximum walkable slope angle (radians).

- *Handle shortcut*: `hero.setMaxSlope(angle)`

---

### `CHARACTERREF.SETSTEPHEIGHT(hero, h)` 

Sets the maximum step height the character can climb.

- *Handle shortcut*: `hero.setStepHeight(h)`

---

### `CHARACTERREF.SETSETTING(hero, key, value)` 

Generic `CharacterVirtual` setting by string key. For advanced Jolt tuning not covered by named setters.

- *Handle shortcut*: `hero.setSetting(key, value)`

---

### `CHARACTERREF.SETCONTACTLISTENER(hero, fn)` 

Attaches a contact callback function called on each contact event.

- *Handle shortcut*: `hero.setContactListener(fn)`

---

## Ground & State Queries

### `CHARACTERREF.ISGROUNDED(hero)` 

Returns `TRUE` when the character is on flat/traversable ground.

- *Handle shortcut*: `hero.isGrounded()`

---

### `CHARACTERREF.GROUNDSTATE(hero)` 

Returns Jolt `EGroundState` integer: `0` OnGround, `1` OnSteepGround, `2` NotSupported, `3` InAir.

- *Handle shortcut*: `hero.groundState()`

---

### `CHARACTERREF.GETCEILING(hero)` 

Returns `TRUE` if the character is touching a ceiling this frame.

- *Handle shortcut*: `hero.getCeiling()`

---

### `CHARACTERREF.GETISSLIDING(hero)` 

Returns `TRUE` if the character is sliding on a steep slope.

- *Handle shortcut*: `hero.getIsSliding()`

---

### `CHARACTERREF.GETGROUNDVELOCITY(hero)` 

Returns the velocity of the ground surface under the character as a `[vx, vy, vz]` array handle.

- *Handle shortcut*: `hero.getGroundVelocity()`

---

### `CHARACTERREF.GETGROUNDNORMAL(hero)` 

Returns the contact normal under the capsule as a `[nx, ny, nz]` array handle.

- *Handle shortcut*: `hero.getGroundNormal()`

---

## Physics Property Getters

### `CHARACTERREF.GETGRAVITY(hero)` 

Returns current gravity scale. Alias: `CHARACTERREF.GETGRAVITYSCALE`.

---

### `CHARACTERREF.GETFRICTION(hero)` 

Returns current friction value.

---

### `CHARACTERREF.GETBOUNCINESS(hero)` 

Returns current restitution. Alias: `CHARACTERREF.GETBOUNCE`.

---

### `CHARACTERREF.GETPADDING(hero)` 

Returns current collision margin.

---

### `CHARACTERREF.GETMAXSLOPE(hero)` 

Returns current maximum walkable slope (radians).

---

### `CHARACTERREF.GETSTEPHEIGHT(hero)` 

Returns current maximum step height.

---

### `CHARACTERREF.GETSPEED(hero)` 

Returns current horizontal speed scalar.

---

## Contacts

### `CHARACTERREF.DRAINCONTACTS(hero)` 

Returns and clears the pending contact queue for this frame. Useful for custom contact-response logic.

- *Handle shortcut*: `hero.drainContacts()`

---

## Lifetime

### `CHARACTERREF.FREE(hero)` 

Releases the Jolt character and frees the heap slot. Safe to call after `PHYSICS3D.STOP`.

- *Handle shortcut*: `hero.free()`

---

## Full Example

```basic
WINDOW.OPEN(960, 540, "Character")
WINDOW.SETFPS(60)

PHYSICS3D.START()
PHYSICS3D.SETGRAVITY(0, -10, 0)

cam = CAMERA.CREATE()

; Floor
floorDef = BODY3D.CREATE("STATIC")
BODY3D.ADDBOX(floorDef, 50, 0.5, 50)
BODY3D.COMMIT(floorDef, 0, 0, 0)

; Player entity + character controller
playerEnt = MODEL.CREATECAPSULE(0.4, 1.0)
playerEnt.setPos(0, 5, 0)
hero = CHARACTER.CREATE(playerEnt, 0.4, 1.0)
hero.setFriction(0.9).setPadding(0.02)

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    CHARACTERREF.UPDATE(hero, dt)

    ; movement
    dx = 0 : dz = 0
    IF INPUT.KEYDOWN(KEY_W) THEN dz = -5 * dt
    IF INPUT.KEYDOWN(KEY_S) THEN dz =  5 * dt
    IF INPUT.KEYDOWN(KEY_A) THEN dx = -5 * dt
    IF INPUT.KEYDOWN(KEY_D) THEN dx =  5 * dt
    IF INPUT.KEYDOWN(KEY_SPACE) AND hero.isGrounded() THEN hero.jump(5)
    CHARACTERREF.MOVE(hero, dx, 0, dz)

    ; camera follow
    px, py, pz = ENTITY.GETPOS(playerEnt)
    cam.setPos(px, py + 8, pz + 12)
    cam.setTarget(px, py, pz)

    RENDER.CLEAR(20, 30, 40)
    RENDER.BEGIN3D(cam)
        ENTITY.DRAWALL()
        DRAW.GRID(100, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND

hero.free()
PHYSICS3D.STOP()
WINDOW.CLOSE()
```

---

## See also

- [CHARCONTROLLER.md](CHARCONTROLLER.md) — standalone capsule handle (`CHARCONTROLLER.CREATE`)
- [KCC.md](KCC.md) — `CHAR.*` / `PLAYER.*` gameplay layer
- [PLAYER.md](PLAYER.md) — `PLAYER.CREATE`, swim, nav
- [PHYSICS3D.md](PHYSICS3D.md) — world step, picks, entity bridge
