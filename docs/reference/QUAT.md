# Quat Commands

Quaternion math handles for 3D rotation: create from Euler angles or axis-angle, multiply, slerp, invert, convert to matrix or Euler. Complements [VEC_QUAT.md](VEC_QUAT.md).

## Core Workflow

1. `QUAT.FROMEULER(pitch, yaw, roll)` or `QUAT.IDENTITY()` — create a quaternion handle.
2. Combine rotations with `QUAT.MULTIPLY(a, b)`.
3. Rotate vectors with `VEC3.ROTATEBYQUAT(v, q)`.
4. Convert for use in transforms: `QUAT.TOMAT4(q)` or `QUAT.TOEULER(q)`.
5. `QUAT.FREE(q)` when done.

---

## Creation

### `QUAT.IDENTITY()` 

Returns the identity quaternion (no rotation). Handle.

---

### `QUAT.FROMEULER(pitch, yaw, roll)` 

Creates a quaternion from Euler angles in **degrees**. Returns a handle.

---

### `QUAT.FROMAXISANGLE(ax, ay, az, angleDeg)` 

Creates a quaternion for rotation of `angleDeg` around axis `(ax, ay, az)`. Returns a handle.

---

### `QUAT.FROMVEC3TOVEC3(fromHandle, toHandle)` 

Creates the shortest-arc quaternion rotating `from` to `to`. Returns a handle.

---

### `QUAT.FROMMAT4(mat4Handle)` 

Extracts a quaternion from a 4×4 matrix. Returns a handle.

---

## Operations

### `QUAT.MULTIPLY(a, b)` 

Returns a new quaternion equal to `a × b` (applies `b` rotation first).

---

### `QUAT.SLERP(a, b, t)` 

Spherically interpolates between `a` and `b` by `t` (0–1). Returns a new handle.

---

### `QUAT.NORMALIZE(q)` 

Returns a normalised unit quaternion.

---

### `QUAT.INVERT(q)` 

Returns the inverse of `q`.

---

### `QUAT.TRANSFORM(q, vec3Handle)` 

Rotates `vec3Handle` by quaternion `q`. Returns a new Vec3 handle. Alias of `VEC3.ROTATEBYQUAT`.

---

## Conversion

### `QUAT.TOEULER(q)` 

Returns `[pitch, yaw, roll]` in degrees as a Vec3 handle.

---

### `QUAT.TOMAT4(q)` 

Returns a 4×4 rotation matrix handle from the quaternion.

---

## Lifetime

### `QUAT.FREE(q)` 

Frees the quaternion handle.

---

## Full Example

Smooth turret tracking using slerp.

```basic
WINDOW.OPEN(960, 540, "Quat Demo")
WINDOW.SETFPS(60)

cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 5, -10)
CAMERA.SETTARGET(cam, 0, 0, 0)

turret = ENTITY.CREATECUBE(1.0)
ENTITY.SETPOS(turret, 0, 0, 0)

currentRot = QUAT.IDENTITY()
targetYaw  = 0.0

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()

    IF INPUT.KEYDOWN(KEY_LEFT)  THEN targetYaw = targetYaw - 90 * dt
    IF INPUT.KEYDOWN(KEY_RIGHT) THEN targetYaw = targetYaw + 90 * dt

    targetQ  = QUAT.FROMEULER(0, targetYaw, 0)
    currentRot = QUAT.SLERP(currentRot, targetQ, 5.0 * dt)
    QUAT.FREE(targetQ)

    euler = QUAT.TOEULER(currentRot)
    ENTITY.SETROT(turret, 0, VEC3.Y(euler), 0)
    VEC3.FREE(euler)

    ENTITY.UPDATE(dt)

    RENDER.CLEAR(20, 25, 35)
    RENDER.BEGIN3D(cam)
        ENTITY.DRAWALL()
        DRAW3D.GRID(10, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND

QUAT.FREE(currentRot)
ENTITY.FREE(turret)
CAMERA.FREE(cam)
WINDOW.CLOSE()
```

---

## See also

- [VEC_QUAT.md](VEC_QUAT.md) — combined vector and quaternion examples
- [VEC3.md](VEC3.md) — `VEC3.ROTATEBYQUAT`
- [TRANSFORM.md](TRANSFORM.md) — matrix transforms using `QUAT.TOMAT4`
