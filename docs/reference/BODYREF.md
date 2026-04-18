# BodyRef / Kinematic / Static Commands

Handle-method API for moving kinematic and static Jolt bodies after creation. Covers `BODYREF.*`, `KINEMATIC.*`, `KINEMATICREF.*`, and `STATIC.*`.

## Core Workflow

1. Create a shape: `shape = SHAPE.CREATEBOX(w, h, d)` (see [SHAPE.md](SHAPE.md)).
2. Create a body: `body = STATIC.CREATE(shape)` or `body = KINEMATIC.CREATE(shape)`.
3. Set position: `BODYREF.SETPOS(body, x, y, z)`.
4. For kinematic bodies each frame: `KINEMATICREF.SETVELOCITY(body, vx, vy, vz)` → `KINEMATICREF.UPDATE(body)`.
5. Free: `BODYREF.FREE(body)`, then `SHAPEREF.FREE(shape)`.

---

## Overview

| Type | Created via | Moved via |
|---|---|---|
| **Static** | `STATIC.CREATE(shapeHandle)` | `BODYREF.SETPOS` (teleport only) |
| **Kinematic** | `KINEMATIC.CREATE(shapeHandle)` | `KINEMATICREF.SETVELOCITY` + `UPDATE` |
| **Dynamic** | `BODY3D.CREATE` | `BODY3D.*` forces |

---

## STATIC Commands

### `STATIC.CREATE(shapeHandle)` 

Creates a static (immovable) physics body from a `SHAPE.*` handle. Returns a **bodyref handle**. The body participates in collision but never moves under physics forces.

---

## KINEMATIC Commands

### `KINEMATIC.CREATE(shapeHandle)` 

Creates a kinematic body from a `SHAPE.*` handle. Returns a **bodyref handle** usable with `BODYREF.*` and `KINEMATICREF.*`. Kinematic bodies move via code and push dynamic bodies.

---

## BODYREF Commands

All `BODYREF.*` commands operate on handles returned by `STATIC.CREATE`, `KINEMATIC.CREATE`, or `TRIGGER.CREATE`.

### `BODYREF.SETPOS(ref, x, y, z)` 

Teleports the body to world position `(x, y, z)`. For kinematic bodies prefer `KINEMATICREF.SETVELOCITY`.

---

### `BODYREF.SETROTATION(ref, pitch, yaw, roll)` 

Sets body orientation in **degrees**.

---

### `BODYREF.SETLAYER(ref, layer)` 

Sets the Jolt collision layer (0–31).

---

### `BODYREF.ENABLECOLLISION(ref, enabled)` 

Enables or disables the body's participation in collision detection.

---

### `BODYREF.GETPOSITION(ref)` 

Returns world position as a Vec3 handle.

---

### `BODYREF.GETROTATION(ref)` 

Returns orientation as a Vec3 handle.

---

### `BODYREF.GETVELOCITY(ref)` 

Returns velocity as a Vec3 handle.

---

### `BODYREF.SETVELOCITY(ref, vx, vy, vz)` 

Sets the body's velocity directly.

---

### `BODYREF.FREE(ref)` 

Destroys the body and frees the handle.

---

## KINEMATICREF Commands

### `KINEMATICREF.SETVELOCITY(ref, vx, vy, vz)` 

Sets the kinematic body's target velocity. The body moves at this velocity and resolves collisions with dynamic bodies on `UPDATE`.

---

### `KINEMATICREF.UPDATE(ref)` 

Steps the kinematic body forward, applying the set velocity and resolving contacts.

---

## Full Example

A moving platform that lifts dynamic crates.

```basic
WINDOW.OPEN(960, 540, "Kinematic Platform")
WINDOW.SETFPS(60)

PHYSICS3D.START()
PHYSICS3D.SETGRAVITY(0, -10, 0)

; static floor
floorShape = SHAPE.CREATEBOX(15, 0.5, 15)
floor      = STATIC.CREATE(floorShape)
BODYREF.SETPOS(floor, 0, -0.5, 0)

; kinematic platform
platShape = SHAPE.CREATEBOX(3, 0.25, 3)
platform  = KINEMATIC.CREATE(platShape)
BODYREF.SETPOS(platform, 0, 1, 0)
KINEMATICREF.SETVELOCITY(platform, 0, 1.5, 0)

; dynamic crate on top
crateDef = BODY3D.CREATE("DYNAMIC")
BODY3D.ADDBOX(crateDef, 0.5, 0.5, 0.5)
crate = BODY3D.COMMIT(crateDef, 0, 3, 0)

cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 8, -14)
CAMERA.SETTARGET(cam, 0, 4, 0)

direction = 1.0
py = 1.0

WHILE NOT WINDOW.SHOULDCLOSE()
    py = py + 1.5 * direction * TIME.DELTA()
    IF py > 6.0 THEN direction = -1.0
    IF py < 0.5 THEN direction =  1.0
    KINEMATICREF.SETVELOCITY(platform, 0, 1.5 * direction, 0)
    KINEMATICREF.UPDATE(platform)
    PHYSICS3D.UPDATE()

    RENDER.CLEAR(20, 25, 35)
    RENDER.BEGIN3D(cam)
        DRAW3D.GRID(20, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND

BODYREF.FREE(platform)
BODYREF.FREE(floor)
BODY3D.FREE(crate)
PHYSICS3D.STOP()
WINDOW.CLOSE()
```

---

## Extended Command Reference

| Command | Description |
|--------|-------------|
| `BODYREF.SETPOSITION(body, x,y,z)` | Alias of `BODYREF.SETPOS` — teleport static/kinematic body to world position. |

---

## See also

- [SHAPE.md](SHAPE.md) — shape handles for `STATIC.CREATE` / `KINEMATIC.CREATE`
- [BODY3D.md](BODY3D.md) — dynamic bodies
- [TRIGGER.md](TRIGGER.md) — sensor bodies (also use `BODYREF.*`)
- [PHYSICS3D.md](PHYSICS3D.md) — world setup
