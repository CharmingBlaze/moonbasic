# CharacterRef Commands

Jolt `CharacterVirtual`-backed kinematic character controller — higher fidelity than `CHARCONTROLLER`, designed for desktop (Windows/Linux + CGO). Returns a **ref handle** from `CHARACTER.CREATE` for entity-bound characters.

See [CHARACTER.md](CHARACTER.md) for creation. `CHAR.*` are aliases for most commands listed here.

## Core Workflow

1. `CHARACTER.CREATE(entity)` → `CHARACTERREF` handle (see [CHARACTER.md](CHARACTER.md)).
2. Configure: `CHARACTERREF.SETGRAVITY`, `SETMAXSLOPE`, `SETSTEPHEIGHT`, `SETSNAPDISTANCE`.
3. Each frame: `CHARACTERREF.MOVE(ref, dx, dy, dz)` or `CHARACTERREF.MOVEWITHCAMERA`.
4. `CHARACTERREF.UPDATE(ref)` — resolves the KCC step.
5. Read state: `CHARACTERREF.ISGROUNDED`, `GETVELOCITY`, `GETGROUNDSTATE`, etc.
6. `CHARACTERREF.FREE(ref)` when done.

---

## Movement

### `CHARACTERREF.MOVE(ref, dx, dy, dz)` 

Sets the intended movement velocity for this frame. `(dx, dy, dz)` in world units/s.

- *Handle shortcut*: `ref.move(dx, dy, dz)`

---

### `CHARACTERREF.MOVEWITHCAMERA(ref, speed, yaw, camHandle)` 

Camera-relative movement: `speed` in m/s, `yaw` from camera. Handles direction mapping internally.

---

### `CHARACTERREF.ADDVELOCITY(ref, vx, vy, vz)` 

Adds to the current velocity (e.g. for jump or knockback).

---

### `CHARACTERREF.SETLINEARVELOCITY(ref, vx, vy, vz)` / `CHARACTERREF.SETVELOCITY(ref, vx, vy, vz)` 

Overrides the velocity directly. `SETVELOCITY` is an alias.

---

### `CHARACTERREF.JUMP(ref, impulse)` 

Applies an upward jump impulse.

---

### `CHARACTERREF.UPDATE(ref)` 

Resolves the KCC step — call once per frame after setting velocity.

---

### `CHARACTERREF.UPDATEMOVE(ref)` 

Combined velocity-apply + update in one call.

---

## Position

### `CHARACTERREF.SETPOS(ref, x, y, z)` / `CHARACTERREF.SETPOSITION(ref, x, y, z)` 

Teleports the character to world position.

---

### `CHARACTERREF.GETPOSITION(ref)` 

Returns `[x, y, z]` position array.

---

### `CHARACTERREF.GETROT(ref)` 

Returns approximate `[pitch, yaw, roll]` in radians from velocity direction.

---

## Ground & Slope State

### `CHARACTERREF.ISGROUNDED(ref)` 

Returns `TRUE` when standing on a supported surface.

---

### `CHARACTERREF.GETGROUNDSTATE(ref)` 

Returns Jolt `EGroundState`: `0`=OnGround, `1`=OnSteepGround, `2`=NotSupported, `3`=InAir.

---

### `CHARACTERREF.ONSLOPE(ref)` 

Returns `TRUE` if standing on a slope steeper than `GETMAXSLOPE`.

---

### `CHARACTERREF.ONWALL(ref)` 

Returns `TRUE` if pressing against a wall.

---

### `CHARACTERREF.GETISSLIDING(ref)` 

Returns `TRUE` when sliding on steep ground.

---

### `CHARACTERREF.GETCEILING(ref)` 

Returns `TRUE` if head/ceiling contact was detected.

---

### `CHARACTERREF.GETSLOPEANGLE(ref)` 

Returns the current slope angle in degrees.

---

## Velocity Queries

### `CHARACTERREF.GETVELOCITY(ref)` 

Returns current velocity as a 3-element array handle `[vx, vy, vz]`.

---

### `CHARACTERREF.GETSPEED(ref)` 

Returns the horizontal speed scalar.

---

### `CHARACTERREF.ISMOVING(ref)` 

Returns `TRUE` if the character has non-zero horizontal velocity.

---

### `CHARACTERREF.GETGROUNDVELOCITY(ref)` 

Returns the ground/platform velocity `[vx, vy, vz]`.

---

## Configuration

### `CHARACTERREF.SETGRAVITY(ref, g)` 

Sets gravity scale (`1.0` = world gravity, `0.0` = weightless).

---

### `CHARACTERREF.SETGRAVITYSCALE(ref, scale)` 

Same as `SETGRAVITY`.

---

### `CHARACTERREF.GETGRAVITY(ref)` 

Returns the gravity scale.

---

### `CHARACTERREF.SETMAXSLOPE(ref, degrees)` 

Maximum walkable slope angle. Steeper = sliding.

---

### `CHARACTERREF.GETMAXSLOPE(ref)` 

Returns max slope in degrees.

---

### `CHARACTERREF.SETSTEPHEIGHT(ref, height)` 

Stair step-up height in world units.

---

### `CHARACTERREF.GETSTEPHEIGHT(ref)` 

Returns step height.

---

### `CHARACTERREF.SETSNAPDISTANCE(ref, dist)` 

Snap-to-floor distance when walking down ramps.

---

### `CHARACTERREF.GETSNAPDISTANCE(ref)` 

Returns snap distance.

---

### `CHARACTERREF.SETFRICTION(ref, f)` 

Sets surface friction.

---

### `CHARACTERREF.GETFRICTION(ref)` 

Returns friction.

---

### `CHARACTERREF.SETPADDING(ref, p)` 

Sets KCC skin/padding width.

---

### `CHARACTERREF.GETPADDING(ref)` 

Returns padding.

---

### `CHARACTERREF.SETJUMPBUFFER(ref, seconds)` 

Coyote-time jump buffer (allows jump shortly after leaving ground).

---

### `CHARACTERREF.GETJUMPBUFFER(ref)` 

Returns jump buffer duration.

---

### `CHARACTERREF.SETAIRCONTROL(ref, factor)` 

Sets how much the player can steer mid-air (0–1).

---

### `CHARACTERREF.GETAIRCONTROL(ref)` 

Returns air control factor.

---

### `CHARACTERREF.SETGROUNDCONTROL(ref, factor)` 

Scales horizontal input while grounded.

---

### `CHARACTERREF.GETGROUNDCONTROL(ref)` 

Returns ground control factor.

---

### `CHARACTERREF.SETBOUNCE(ref, restitution)` 

Sets collision restitution (bounciness).

---

### `CHARACTERREF.SETSTICKDOWN(ref, enabled)` 

Enables floor-snap on steep descents.

---

### `CHARACTERREF.SETSETTING(ref, key, value)` 

Generic setting setter for future expansion.

---

## Contacts

### `CHARACTERREF.DRAINCONTACTS(ref)` 

Returns an array of contact events accumulated since last drain.

---

### `CHARACTERREF.SETCONTACTLISTENER(ref, handlerName)` 

Registers a callback function called on KCC contact events.

---

## Lifetime

### `CHARACTERREF.FREE(ref)` 

Destroys the kinematic character controller.

---

## Full Example

```basic
WINDOW.OPEN(960, 540, "CharacterRef Demo")
WINDOW.SETFPS(60)

PHYSICS3D.START()
PHYSICS3D.SETGRAVITY(0, -20, 0)

e    = ENTITY.CREATECUBE(1.0)
ref  = CHARACTER.CREATE(e)
CHARACTERREF.SETSTEPHEIGHT(ref, 0.4)
CHARACTERREF.SETMAXSLOPE(ref, 45)
CHARACTERREF.SETGRAVITY(ref, 1.0)

cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 6, -10)
CAMERA.SETTARGET(cam, 0, 0, 0)

WHILE NOT WINDOW.SHOULDCLOSE()
    dt  = TIME.DELTA()
    spd = 5.0
    dx  = 0.0 : dz = 0.0
    IF INPUT.KEYDOWN(KEY_D) THEN dx =  spd
    IF INPUT.KEYDOWN(KEY_A) THEN dx = -spd
    IF INPUT.KEYDOWN(KEY_W) THEN dz = -spd
    IF INPUT.KEYDOWN(KEY_S) THEN dz =  spd

    CHARACTERREF.MOVE(ref, dx, 0, dz)
    IF INPUT.KEYPRESSED(KEY_SPACE) AND CHARACTERREF.ISGROUNDED(ref) THEN
        CHARACTERREF.JUMP(ref, 8.0)
    END IF
    CHARACTERREF.UPDATE(ref)

    PHYSICS3D.UPDATE()
    ENTITY.UPDATE(dt)

    RENDER.CLEAR(20, 25, 35)
    RENDER.BEGIN3D(cam)
        ENTITY.DRAWALL()
        DRAW3D.GRID(20, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND

CHARACTERREF.FREE(ref)
ENTITY.FREE(e)
PHYSICS3D.STOP()
WINDOW.CLOSE()
```

---

## See also

- [CHARACTER.md](CHARACTER.md) — `CHARACTER.CREATE` to get a ref handle
- [CHARCONTROLLER.md](CHARCONTROLLER.md) — simpler capsule controller
- [PLAYER.md](PLAYER.md) — high-level player (`CHAR.*` aliases)
- [PHYSICS3D.md](PHYSICS3D.md) — world setup
