# Joint3D Commands

3D Jolt Physics constraints connecting two `BODY3D` bodies. Supports fixed, hinge, slider, and cone constraints.

Requires **CGO + Jolt** (Windows/Linux fullruntime) and an active `PHYSICS3D` session.

## Core Workflow

1. Create and commit two `BODY3D` bodies.
2. Call the joint command — bodies must already be in the world.
3. The constraint is enforced automatically on each `PHYSICS3D.UPDATE`.
4. `JOINT3D.DELETE(joint)` to destroy.

---

## Joints

### `JOINT3D.FIXED(bodyA, bodyB)` 

Rigidly welds two bodies together. No relative motion allowed. Returns a **joint handle**.

---

### `JOINT3D.HINGE(bodyA, bodyB, ax, ay, az, pivotBx, pivotBy, pivotBz)` 

Creates a hinge constraint between `bodyA` and `bodyB`. The hinge axis is `(ax, ay, az)` in bodyA's local space. `(pivotBx, pivotBy, pivotBz)` is the pivot point in bodyB's local space. Returns a **joint handle**.

---

### `JOINT3D.SLIDER(bodyA, bodyB, ax, ay, az, pivotBx, pivotBy, pivotBz)` 

Creates a prismatic (slider) constraint — bodies can only translate along axis `(ax, ay, az)`. Returns a **joint handle**.

---

### `JOINT3D.CONE(bodyA, bodyB, ax, ay, az, halfAngle, pivotBx, pivotBy)` 

Creates a cone constraint limiting the angle between the two bodies' axes to `halfAngle` radians. Returns a **joint handle**.

---

## Lifetime

### `JOINT3D.DELETE(joint)` 

Destroys the joint constraint.

---

## Full Example

A hinged door swinging on a Y-axis hinge.

```basic
WINDOW.OPEN(960, 540, "Joint3D Demo")
WINDOW.SETFPS(60)

PHYSICS3D.START()
PHYSICS3D.SETGRAVITY(0, -10, 0)

; static door frame
frameDef = BODY3D.CREATE("STATIC")
BODY3D.ADDBOX(frameDef, 0.1, 2, 0.1)
frame = BODY3D.COMMIT(frameDef, -1, 2, 0)

; dynamic door panel
doorDef = BODY3D.CREATE("DYNAMIC")
BODY3D.ADDBOX(doorDef, 1, 2, 0.1)
BODY3D.SETMASS(doorDef, 5.0)
door = BODY3D.COMMIT(doorDef, 0, 2, 0)

; hinge on Y axis at the left edge
hinge = JOINT3D.HINGE(frame, door, 0, 1, 0, -1, 0, 0)

cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 4, -8)
CAMERA.SETTARGET(cam, 0, 2, 0)

WHILE NOT WINDOW.SHOULDCLOSE()
    IF INPUT.KEYPRESSED(KEY_SPACE) THEN
        BODY3D.APPLYIMPULSE(door, 0, 0, 8)
    END IF

    PHYSICS3D.UPDATE()

    RENDER.CLEAR(20, 25, 35)
    RENDER.BEGIN3D(cam)
        DRAW3D.GRID(20, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND

JOINT3D.DELETE(hinge)
BODY3D.FREE(door)
BODY3D.FREE(frame)
PHYSICS3D.STOP()
WINDOW.CLOSE()
```

---

## See also

- [BODY3D.md](BODY3D.md) — 3D rigid bodies
- [PHYSICS3D.md](PHYSICS3D.md) — world setup
- [PHYSICS_ADVANCED.md](PHYSICS_ADVANCED.md) — joints, CCD, advanced constraints
- [JOINT2D.md](JOINT2D.md) — 2D joints
