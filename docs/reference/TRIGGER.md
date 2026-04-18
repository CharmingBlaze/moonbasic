# Trigger Commands

Non-solid sensor bodies (Jolt triggers) that fire events when other bodies enter or exit their volume. Use for pickups, checkpoints, damage zones, and area detection.

Requires **CGO + Jolt** (Windows/Linux fullruntime).

## Core Workflow

1. Create a shape with `SHAPE.CREATEBOX` / `SHAPE.CREATESPHERE` etc.
2. `TRIGGER.CREATE(shapeHandle)` — wrap it as a non-solid sensor.
3. Or `TRIGGER.CREATEFROMENTITY(entityId)` — use an existing entity's shape.
4. Or `TRIGGER.CREATEZONE(x, y, z, hw, hh, hd, tag)` — quick named zone.
5. The trigger fires hit events readable via the event system when entities enter/exit.

---

## Creation

### `TRIGGER.CREATE(shapeHandle)` 

Creates a non-solid trigger body from a `SHAPE.*` handle. The body is placed at the shape's default origin. Returns a **trigger handle**.

---

### `TRIGGER.CREATEFROMENTITY(entityId)` 

Creates a trigger sensor using an existing entity's collision shape and position. The entity itself remains a visual only.

---

### `TRIGGER.CREATEZONE(x, y, z, halfW, halfH, halfD, tag)` 

Creates a named box-shaped trigger zone at `(x, y, z)` with half-extents. `tag` is a string event name fired when an entity enters. Returns a **trigger handle**.

---

## Full Example

A pickup zone that fires when the player enters it.

```basic
WINDOW.OPEN(960, 540, "Trigger Demo")
WINDOW.SETFPS(60)

PHYSICS3D.START()
PHYSICS3D.SETGRAVITY(0, -10, 0)

cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 8, -12)
CAMERA.SETTARGET(cam, 0, 0, 0)

; coin pickup zone
coinZone = TRIGGER.CREATEZONE(0, 1, 0, 1, 1, 1, "coin_pickup")

; player character
ctrl = CHARCONTROLLER.CREATE(0.4, 1.8, 0, 5, 0)

collected = FALSE
msg = "Walk into the zone!"

EVENT.ON("coin_pickup", FUNCTION()
    collected = TRUE
    msg = "Coin collected!"
END FUNCTION)

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    dx = 0
    dz = 0
    IF INPUT.KEYDOWN(KEY_D) THEN dx =  4 * dt
    IF INPUT.KEYDOWN(KEY_A) THEN dx = -4 * dt
    IF INPUT.KEYDOWN(KEY_W) THEN dz = -4 * dt
    IF INPUT.KEYDOWN(KEY_S) THEN dz =  4 * dt
    CHARCONTROLLER.MOVE(ctrl, dx, 0, dz)

    PHYSICS3D.UPDATE()

    RENDER.CLEAR(20, 25, 35)
    RENDER.BEGIN3D(cam)
        IF NOT collected THEN
            DRAW3D.CUBE(0, 1, 0, 2, 2, 2, 255, 200, 60, 120)
        END IF
        DRAW3D.SPHERE(CHARCONTROLLER.X(ctrl), CHARCONTROLLER.Y(ctrl), CHARCONTROLLER.Z(ctrl), 0.4, 80, 160, 255, 255)
        DRAW3D.GRID(20, 1.0)
    RENDER.END3D()
    DRAW.TEXT(msg, 10, 10, 20, 255, 255, 255, 255)
    RENDER.FRAME()
WEND

CHARCONTROLLER.FREE(ctrl)
PHYSICS3D.STOP()
WINDOW.CLOSE()
```

---

## Extended Command Reference

| Command | Description |
|--------|-------------|
| `TRIGGER.MAKE(shape)` | Deprecated alias of `TRIGGER.CREATE`. |
| `TRIGGER.MAKEFROMENTITY(entity)` | Create a trigger zone matching entity bounds. |
| `TRIGGER.MAKEZONE(x,y,z, hw,hh,hd)` | Create an AABB trigger zone at position with half-extents. |

---

## See also

- [SHAPE.md](SHAPE.md) — shapes for `TRIGGER.CREATE`
- [EVENT.md](EVENT.md) — event system for trigger callbacks
- [PHYSICS3D.md](PHYSICS3D.md) — world setup
- [CHARCONTROLLER.md](CHARCONTROLLER.md) — player controller
