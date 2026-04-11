# Gameplay Helpers

A collection of beginner-friendly building blocks for common game patterns.

---

## Entity Verbs

### `Entity.Distance(a, b)`
Returns the 3D world distance between two entities.

### `Entity.DistanceXZ(a, b)`
Returns the distance between two entities on the XZ plane (ignoring height).

---

### `Entity.Within(a, b, radius)`
Returns `TRUE` if entity `b` is within a 3D `radius` of entity `a`.

### `Entity.WithinXZ(a, b, radius)`
Returns `TRUE` if entity `b` is within a 2D `radius` of entity `a` on the XZ plane.

---

### `Entity.MoveToward(id, target, speed)`
Moves an entity toward a target position or another entity at a specific speed.

### `Entity.TurnToward(id, target, speed)`
Smoothly rotates an entity to face a target point or entity.

### `Entity.LookAt(id, x, y, z)`
Instantly makes an entity face a world position.

---

## World Helpers

### `World.GetGroundY(x, z)`
Returns the terrain height at the specified world coordinates.

### `World.Raycast(ox, oy, oz, dx, dy, dz, max)`
Returns the ID of the first entity hit by a ray, or `0` if none.

---

## Character & Camera

### `Entity.MoveWithCamera(id, cam, f, s, speed)`
Moves an entity relative to the camera's orientation (WASD style).

### `Camera.Orbit(cam, id, dist)`
Sets up a third-person orbit camera around an entity.

---

## Examples

### Proximity Trigger
```basic
IF Entity.WithinXZ(player, door, 2.0) THEN
    Entity.TurnToward(door, player, 5.0)
    PRINT "Door opening..."
ENDIF
```

### Camera-Relative Movement
```basic
f = Input.Axis(KEY_S, KEY_W)
s = Input.Axis(KEY_A, KEY_D)
Entity.MoveWithCamera(player, mainCam, f, s, 10.0)
```

