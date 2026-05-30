# Body2D Commands

2D rigid body creation and simulation using Box2D. Build body definitions, add shapes, commit to the world, then drive physics with `PHYSICS2D.STEP`.

## Core Workflow

1. `PHYSICS2D.START()` — initialise the 2D physics world.
2. `BODY2D.CREATE(type)` — begin a body definition (`"DYNAMIC"`, `"STATIC"`, or `"KINEMATIC"`).
3. `BODY2D.ADDRECT` / `BODY2D.ADDCIRCLE` / `BODY2D.ADDPOLYGON` — attach collision shapes.
4. `BODY2D.COMMIT(def, x, y)` — finalise and place the body; returns a **live body handle**.
5. Each frame: `PHYSICS2D.STEP()` → read `BODY2D.X` / `BODY2D.Y` → sync visuals.
6. `BODY2D.FREE(handle)` when done.

---

## Creation

### `BODY2D.CREATE(type)` 

Begins a new body definition. `type` is `"DYNAMIC"`, `"STATIC"`, or `"KINEMATIC"`. Returns a **definition handle** (not yet in the world — call `BODY2D.COMMIT` to finalise).

---

### `BODY2D.ADDRECT(def, halfW, halfH)` 

Attaches a rectangular fixture to a body definition. `halfW` and `halfH` are half-extents in world units (so a 2×1 box uses `1.0, 0.5`).

---

### `BODY2D.ADDCIRCLE(def, radius)` 

Attaches a circle fixture with the given `radius` to the body definition.

---

### `BODY2D.ADDPOLYGON(def, pointsHandle)` 

Attaches a convex polygon fixture. `pointsHandle` is a flat numeric array handle of alternating `x, y` pairs (max 8 vertices, must be convex).

---

### `BODY2D.COMMIT(def, x, y)` 

Finalises the body definition at world position `(x, y)` and returns the **live body handle**. The definition handle is consumed — do not reuse it.

---

## Position & Rotation

### `BODY2D.SETPOS(body, x, y)` 

Teleports the body to world coordinates `(x, y)`. Wakes the body if sleeping.

- *Handle shortcut*: `body.setPos(x, y)`

---

### `BODY2D.GETPOS(body)` 

Returns the current world position as a 2-element array handle `[x, y]`.

- *Handle shortcut*: `body.getPos()`

---

### `BODY2D.X(body)` 

Returns the body's X world coordinate as a float scalar.

---

### `BODY2D.Y(body)` 

Returns the body's Y world coordinate as a float scalar.

---

### `BODY2D.SETROT(body, angle)` 

Sets the body's rotation in **radians**.

- *Handle shortcut*: `body.setRot(angle)`

---

### `BODY2D.GETROT(body)` 

Returns the current rotation in **radians**.

- *Handle shortcut*: `body.getRot()`

---

### `BODY2D.ROT(body)` 

Alias of `BODY2D.GETROT`. Returns rotation in **radians**.

---

## Mass & Material

### `BODY2D.SETMASS(body, mass)` 

Overrides the body mass in kg. Only meaningful on `DYNAMIC` bodies.

---

### `BODY2D.GETMASS(body)` 

Returns the body mass in kg.

---

### `BODY2D.SETFRICTION(body, friction)` 

Sets the surface friction coefficient (0 = frictionless, 1 = high friction).

---

### `BODY2D.GETFRICTION(body)` 

Returns the current friction coefficient.

---

### `BODY2D.SETRESTITUTION(body, restitution)` 

Sets the bounciness (0 = no bounce, 1 = perfectly elastic).

---

### `BODY2D.GETRESTITUTION(body)` 

Returns the current restitution value.

---

## Velocity & Forces

### `BODY2D.SETLINEARVELOCITY(body, vx, vy)` 

Sets the linear velocity directly in world units per second.

- *Handle shortcut*: `body.setLinearVelocity(vx, vy)`

---

### `BODY2D.GETLINEARVELOCITY(body)` 

Returns linear velocity as a 2-element array handle `[vx, vy]`.

- *Handle shortcut*: `body.getLinearVelocity()`

---

### `BODY2D.SETANGULARVELOCITY(body, omega)` 

Sets the angular velocity in radians per second.

---

### `BODY2D.GETANGULARVELOCITY(body)` 

Returns the angular velocity in radians per second.

---

### `BODY2D.APPLYFORCE(body, fx, fy)` 

Applies a continuous force (Newtons) at the body centre this step.

- *Handle shortcut*: `body.applyForce(fx, fy)`

---

### `BODY2D.APPLYIMPULSE(body, ix, iy)` 

Applies an instantaneous impulse (kg·m/s) at the body centre.

- *Handle shortcut*: `body.applyImpulse(ix, iy)`

---

## Collision Queries

### `BODY2D.COLLIDED(body)` 

Returns `1` if the body had at least one contact during the last `PHYSICS2D.STEP`, `0` otherwise.

---

### `BODY2D.COLLISIONOTHER(body)` 

Returns the **handle** of the other body involved in the most recent collision.

---

### `BODY2D.COLLISIONPOINT(body)` 

Returns a 2-element array handle `[x, y]` of the world-space contact point.

---

### `BODY2D.COLLISIONNORMAL(body)` 

Returns a 2-element array handle `[nx, ny]` — the contact normal pointing away from the other body.

---

## Lifetime

### `BODY2D.FREE(body)` 

Destroys the body and removes it from the physics world.

- *Handle shortcut*: `body.free()`

---

## Full Example

A bouncing ball and a static floor.

```basic
WINDOW.OPEN(800, 600, "Body2D Demo")
WINDOW.SETFPS(60)

PHYSICS2D.START()
PHYSICS2D.SETGRAVITY(0, 500)

; Static floor
floorDef = BODY2D.CREATE("STATIC")
BODY2D.ADDRECT(floorDef, 400, 10)
floor = BODY2D.COMMIT(floorDef, 400, 580)

; Dynamic ball
ballDef = BODY2D.CREATE("DYNAMIC")
BODY2D.ADDCIRCLE(ballDef, 20)
BODY2D.SETRESTITUTION(ballDef, 0.8)
ball = BODY2D.COMMIT(ballDef, 400, 100)

WHILE NOT WINDOW.SHOULDCLOSE()
    PHYSICS2D.STEP()

    bx = BODY2D.X(ball)
    by = BODY2D.Y(ball)

    RENDER.CLEAR(20, 20, 40)
    DRAW.CIRCLE(INT(bx), INT(by), 20, 80, 160, 255, 255)
    DRAW.RECT(0, 570, 800, 20, 100, 100, 100, 255)
    RENDER.FRAME()
WEND

BODY2D.FREE(ball)
BODY2D.FREE(floor)
PHYSICS2D.STOP()
WINDOW.CLOSE()
```

---

## Extended Command Reference

| Command | Description |
|--------|-------------|
| `BODY2D.MAKE(...)` | Deprecated alias of `BODY2D.CREATE`. |
| `BODY2D.SETPOSITION(body, x, y)` | Teleport body to 2D position (wakes it). |

---

## See also

- [PHYSICS2D.md](PHYSICS2D.md) — world setup, stepping, gravity, rope
- [JOINT2D.md](JOINT2D.md) — distance, revolute, prismatic joints
- [COLLISION.md](COLLISION.md) — overlap and distance queries
