# Gameplay Helpers

A collection of beginner-friendly building blocks for common game patterns.

---

## Entity verbs

### `ENTITY.DISTANCE(a, b)` / `ENTITY.DIST` / `ENTITY.DISTANCETO(a, b)`

Returns the **3D** world distance between two entities (numeric **entity#** or **EntityRef** handle). **`ENTITY.DISTANCETO`** is the same distance; **`ENTITY.DIST`** is a short alias.

### Horizontal (XZ) distance

There is no separate **XZ-only** built-in. Use **`ENTITY.GETPOS`** on both entities and compare on **X/Z**, or use **`ENTITY.WITHINRADIUS`** when a **3D** sphere test is acceptable.

---

### `ENTITY.WITHINRADIUS(a, b, radius)`

Returns **`TRUE`** if entity **`b`** is within **3D** **`radius`** of entity **`a`** (center-to-center sphere test).

---

### `ENTITY.MOVETOWARD(id, target, speed)` / `ENTITY.MOVETOWARD(id, x, z, speed)`

Moves an entity toward a **target entity** or a world **(x, z)** on the ground plane at **`speed`**.

### `ENTITY.TURNTOWARD(id, targetX, targetZ, turnSpeed)`

Smoothly rotates an entity to face a world **(x, z)** position at **`turnSpeed`**.

### `ENTITY.LOOKAT(id, x, z)`

Instantly makes an entity face a world **(x, z)** position.

---

## World / terrain

### `TERRAIN.GETHEIGHT(terrain, x, z)`

Returns the terrain surface height at world **(x, z)** for a loaded heightmap terrain **handle** (see [TERRAIN.md](TERRAIN.md)).

### `ENTITY.RAYCAST(ox, oy, oz, dx, dy, dz, maxDist)`

Physics ray cast (Jolt path when linked); returns the first hit **entity#** or **`0`**. Same family as **`PHYSICS3D.RAYCAST`** / picking — see [ENTITY.md](ENTITY.md) and [PHYSICS3D.md](PHYSICS3D.md).

---

## Character and camera

### `ENTITY.MOVEWITHCAMERA(id, cam, forward, strafe, speed)`

Moves an entity relative to the camera's orientation (WASD style).

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

### Camera-relative movement

```basic
f = INPUT.AXIS(KEY_S, KEY_W)
s = INPUT.AXIS(KEY_A, KEY_D)
ENTITY.MOVEWITHCAMERA(player, mainCam, f, s, 10.0)
```
