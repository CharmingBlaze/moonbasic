# Joint2D Commands

2D physics joints connecting two `BODY2D` bodies. Supports distance, revolute (hinge), and prismatic (slider) constraints.

Requires `PHYSICS2D.START()` before use.

## Core Workflow

1. Create two `BODY2D` bodies and `BODY2D.COMMIT` them.
2. Call the appropriate joint command to connect them.
3. The joint is active automatically on each `PHYSICS2D.STEP`.
4. `JOINT2D.FREE(joint)` to destroy.

---

## Joints

### `JOINT2D.DISTANCE(bodyA, bodyB, anchorAx, anchorAy, anchorBx, anchorBy)` 

Creates a distance joint keeping the two anchor points a fixed length apart. Anchors are in local body space. Returns a **joint handle**.

---

### `JOINT2D.REVOLUTE(bodyA, bodyB, anchorX, anchorY)` 

Creates a revolute (hinge/pin) joint at world-space anchor `(anchorX, anchorY)`. Bodies can rotate freely around the pin. Returns a **joint handle**.

---

### `JOINT2D.PRISMATIC(bodyA, bodyB, anchorAx, anchorAy, anchorBx, anchorBy)` 

Creates a prismatic (slider) joint: bodies can only translate along the axis defined by the two anchor points, no rotation. Returns a **joint handle**.

---

## Lifetime

### `JOINT2D.FREE(joint)` 

Destroys the joint.

---

## Full Example

A pendulum: dynamic ball pinned to a static anchor.

```basic
WINDOW.OPEN(800, 600, "Joint2D Demo")
WINDOW.SETFPS(60)

PHYSICS2D.START()
PHYSICS2D.SETGRAVITY(0, 800)

; static anchor
anchorDef = BODY2D.CREATE("STATIC")
BODY2D.ADDRECT(anchorDef, 5, 5)
anchor = BODY2D.COMMIT(anchorDef, 400, 100)

; dynamic pendulum bob
bobDef = BODY2D.CREATE("DYNAMIC")
BODY2D.ADDCIRCLE(bobDef, 18)
BODY2D.SETMASS(bobDef, 1.0)
bob = BODY2D.COMMIT(bobDef, 550, 100)

; pin joint at anchor centre
pin = JOINT2D.REVOLUTE(anchor, bob, 400, 100)

WHILE NOT WINDOW.SHOULDCLOSE()
    PHYSICS2D.STEP()

    bx = BODY2D.X(bob)
    by = BODY2D.Y(bob)

    RENDER.CLEAR(15, 15, 30)
    DRAW.LINE(400, 100, INT(bx), INT(by), 120, 120, 160, 255)
    DRAW.CIRCLE(400, 100, 6, 200, 200, 200, 255)
    DRAW.CIRCLE(INT(bx), INT(by), 18, 80, 160, 255, 255)
    RENDER.FRAME()
WEND

JOINT2D.FREE(pin)
BODY2D.FREE(bob)
BODY2D.FREE(anchor)
PHYSICS2D.STOP()
WINDOW.CLOSE()
```

---

## See also

- [BODY2D.md](BODY2D.md) — 2D rigid bodies
- [PHYSICS2D.md](PHYSICS2D.md) — 2D world setup
- [JOINT3D.md](JOINT3D.md) — 3D joints
