# Gameplay Helpers

A collection of beginner-friendly building blocks for common game patterns: entity distance/movement, proximity triggers, and camera-relative controls.

## Core Workflow

1. Use **`ENTITY.WITHINRADIUS`** or **`ENTITY.DISTANCE`** to detect proximity.
2. Use **`ENTITY.MOVETOWARD`** or **`ENTITY.MOVEWITHCAMERA`** to drive movement.
3. Read **`ENTITY.GETPOS`** for world-space positions when needed.

---

## Entity verbs

### `ENTITY.DISTANCE(a, b)` / `ENTITY.DIST` / `ENTITY.DISTANCETO(a, b)` 

Returns the **3D** world distance between two entities (numeric **entity#** or **EntityRef** handle). **`ENTITY.DISTANCETO`** is the same distance; **`ENTITY.DIST`** is a short alias.

---

### Horizontal (XZ) distance 

There is no separate **XZ-only** built-in. Use **`ENTITY.GETPOS`** on both entities and compare on **X/Z**, or use **`ENTITY.WITHINRADIUS`** when a **3D** sphere test is acceptable.

---

### `ENTITY.WITHINRADIUS(a, b, radius)` 

Returns **`TRUE`** if entity **`b`** is within **3D** **`radius`** of entity **`a`** (center-to-center sphere test).

---

### `ENTITY.MOVETOWARD(id, target, speed)` / `ENTITY.MOVETOWARD(id, x, z, speed)` 

Moves an entity toward a **target entity** or a world **(x, z)** on the ground plane at **`speed`**.

---

### `ENTITY.TURNTOWARD(id, targetX, targetZ, turnSpeed)` 

Smoothly rotates an entity to face a world **(x, z)** position at **`turnSpeed`**.

---

### `ENTITY.LOOKAT(id, x, z)` 

Instantly makes an entity face a world **(x, z)** position.

---

## World / terrain

### `TERRAIN.GETHEIGHT(terrain, x, z)` 

Returns the terrain surface height at world **(x, z)** for a loaded heightmap terrain **handle** (see [TERRAIN.md](TERRAIN.md)).

---

### `ENTITY.RAYCAST(ox, oy, oz, dx, dy, dz, maxDist)` 

Physics ray cast (Jolt path when linked); returns the first hit **entity#** or **`0`**. Same family as **`PHYSICS3D.RAYCAST`** / picking — see [ENTITY.md](ENTITY.md) and [PHYSICS3D.md](PHYSICS3D.md).

---

## Character and camera

### `ENTITY.MOVEWITHCAMERA(id, cam, forward, strafe, speed)` 

Moves an entity relative to the camera's orientation (WASD style).

---

### `CAMERA.ORBITENTITY(cam, entity, yaw, pitch, dist)` 

Sets up a third-person orbit camera around an entity (see runtime **`CAMERA.ORBITENTITY`** — **yaw**, **pitch**, **dist** in radians / world units).

---

## Examples

### Proximity trigger 

```basic
px, py, pz = ENTITY.GETPOS(player)
IF ENTITY.WITHINRADIUS(player, door, 2.0) THEN
    ENTITY.TURNTOWARD(door, px, pz, 5.0)
    PRINT "Door opening..."
ENDIF
```

---

### Camera-relative movement 

```basic
f = INPUT.AXIS(KEY_S, KEY_W)
s = INPUT.AXIS(KEY_A, KEY_D)
ENTITY.MOVEWITHCAMERA(player, mainCam, f, s, 10.0)
```

---

## Full Example

Player entity that follows camera direction and opens a door on proximity.

```basic
WINDOW.OPEN(960, 540, "Gameplay Helpers")
WINDOW.SETFPS(60)

cam    = CAMERA.CREATE()
player = ENTITY.CREATECUBE(1.0)
door   = ENTITY.CREATECUBE(0.2)
ENTITY.SETPOS(player, 0, 0.5, 0)
ENTITY.SETPOS(door,   5, 0.5, 0)

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()

    ; camera-relative WASD
    f = INPUT.AXIS(KEY_S, KEY_W)
    s = INPUT.AXIS(KEY_A, KEY_D)
    ENTITY.MOVEWITHCAMERA(player, cam, f, s, 6.0)

    ; proximity check
    IF ENTITY.WITHINRADIUS(player, door, 2.0)
        PRINT "Near door!"
    END IF

    ENTITY.UPDATE(dt)
    px, py, pz = ENTITY.GETPOS(player)
    CAMERA.ORBIT(cam, player, 8.0)

    RENDER.CLEAR(20, 30, 40)
    RENDER.BEGIN3D(cam)
        ENTITY.DRAWALL()
        DRAW3D.GRID(20, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND

ENTITY.FREE(player)
ENTITY.FREE(door)
CAMERA.FREE(cam)
WINDOW.CLOSE()
```

---

## See also

- [ENTITY.md](ENTITY.md) — full entity API
- [INPUT.md](INPUT.md) — `INPUT.AXIS` and key constants
- [CAMERA.md](CAMERA.md) — orbit and follow cameras
