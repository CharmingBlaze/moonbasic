# 2D Physics Commands

2D rigid-body simulation using **Box2D**. Registry names use **`PHYSICS2D.*`** / **`BODY2D.*`** / **`JOINT2D.*`** prefixes. Full map: [moonbasic-command-set/physics-2d.md](moonbasic-command-set/physics-2d.md).

## Core Workflow

1. **`PHYSICS2D.START()`** — initialise the world once.
2. **`BODY2D.CREATE(type)`** → add shapes → **`BODY2D.COMMIT(def, x, y)`** — build each body.
3. Each frame: **`PHYSICS2D.STEP()`** — advance simulation.
4. Read **`BODY2D.X`** / **`BODY2D.Y`** / **`BODY2D.ROT`** to sync visuals.
5. **`PHYSICS2D.STOP()`** — tear down when done.

### Method chaining 

All **`BODY2D.*`** mutating builtins return the body handle on success: `body.setPos(x,y).applyForce(fx,fy)`.

---

## World Management

### `PHYSICS2D.START([gx, gy])` 

Initialises the 2D physics world. Default gravity `(0, 500)` if omitted (positive Y is down).

---

### `PHYSICS2D.STOP()` 

Shuts down the simulation and frees all internal buffers.

---

### `PHYSICS2D.STEP()` 

Advances the simulation one step. Call once per frame.

---

### `PHYSICS2D.SETGRAVITY(x, y)` 

Sets the global gravity vector.

---

## Body Creation

### `BODY2D.CREATE(type)` 

Creates a body *definition* (not yet in the world). `type`: `"static"`, `"dynamic"`, or `"kinematic"`. Alias: `BODY2D.MAKE` (deprecated).

---

### `BODY2D.ADDRECT(def, w, h)` 

Adds a rectangle collision shape to the definition.

- *Handle shortcut*: `def.addRect(w, h)`

---

### `BODY2D.ADDCIRCLE(def, radius)` 

Adds a circle collision shape to the definition.

- *Handle shortcut*: `def.addCircle(radius)`

---

### `BODY2D.ADDPOLYGON(def, ...)` 

Adds a convex polygon shape. See `commands.json` for vertex arity overloads.

- *Handle shortcut*: `def.addPolygon(...)`

---

### `BODY2D.COMMIT(def, x, y)` 

Finalises the body definition and inserts it into the world at `(x, y)`. Returns a **body handle**. The definition handle is consumed.

- *Handle shortcut*: `def.commit(x, y)`

---

### `BODY2D.SETMASS(body, mass)` 

Sets the mass of a dynamic body. Call after `COMMIT`.

---

### `BODY2D.SETRESTITUTION(body, r)` 

Sets the restitution (bounciness) coefficient (0..1).

---

### `BODY2D.SETFRICTION(body, f)` 

Sets the friction coefficient.

---

## Body Transform

### `BODY2D.SETPOS(body, x, y)` 

Teleports the body to a new position. Alias: `BODY2D.SETPOSITION` (deprecated).

- *Handle shortcut*: `body.setPos(x, y)`

---

### `BODY2D.X(body)` / `BODY2D.Y(body)` 

Returns the current X or Y coordinate of the body's centre.

- *Handle shortcut*: `body.x()` / `body.y()`

---

### `BODY2D.ROT(body)` 

Returns the current rotation in radians.

- *Handle shortcut*: `body.rot()`

---

### `BODY2D.GETPOS(body)` 

Returns position as two floats `x, y`.

---

### `BODY2D.GETVELOCITY(body)` / `BODY2D.GETMASS(body)` / `BODY2D.GETRESTITUTION(body)` 

Getters for velocity vector, mass, and restitution respectively.

---

## Velocity & Forces

### `BODY2D.SETVEL(body, vx, vy)` / `BODY2D.SETVELOCITY(body, vx, vy)` 

Sets the linear velocity directly.

- *Handle shortcut*: `body.setVel(vx, vy)`

---

### `BODY2D.APPLYFORCE(body, fx, fy)` 

Applies a continuous force (mass-scaled acceleration). Alias: `BODY2D.ADDFORCE`.

- *Handle shortcut*: `body.applyForce(fx, fy)`

---

### `BODY2D.APPLYIMPULSE(body, ix, iy)` 

Applies an instant velocity-change impulse. Alias: `BODY2D.ADDIMPULSE`.

- *Handle shortcut*: `body.applyImpulse(ix, iy)`

---

## Body State

### `BODY2D.FREE(body)` 

Removes the body from the simulation and frees its heap slot.

- *Handle shortcut*: `body.free()`

---

## Per-Body Collision Queries

Call these after **`PHYSICS2D.STEP`** each frame.

### `BODY2D.COLLIDED(body)` 

Returns `TRUE` if this body overlapped another body this step.

- *Handle shortcut*: `body.collided()`

---

### `BODY2D.COLLISIONOTHER(body)` 

Returns the **other** body handle from the last collision this step.

- *Handle shortcut*: `body.collisionOther()`

---

### `BODY2D.COLLISIONPOINT(body)` 

Returns the world contact point as a `[x, y]` array handle.

- *Handle shortcut*: `body.collisionPoint()`

---

### `BODY2D.COLLISIONNORMAL(body)` 

Returns the contact surface normal as a `[nx, ny]` array handle.

- *Handle shortcut*: `body.collisionNormal()`

---

## Full Example: Falling Box

```basic
WINDOW.OPEN(800, 600, "2D Physics")
WINDOW.SETFPS(60)

PHYSICS2D.START()
PHYSICS2D.SETGRAVITY(0, 500) ; positive Y is down in 2D

; Static floor
floorDef  = BODY2D.CREATE("STATIC")
BODY2D.ADDRECT(floorDef, 800, 50)
floorBody = BODY2D.COMMIT(floorDef, 400, 575)

; Dynamic box
boxDef  = BODY2D.CREATE("DYNAMIC")
BODY2D.ADDRECT(boxDef, 40, 40)
boxBody = BODY2D.COMMIT(boxDef, 400, 100)

WHILE NOT WINDOW.SHOULDCLOSE()
    PHYSICS2D.STEP()

    RENDER.CLEAR(10, 20, 30)
    CAMERA2D.BEGIN()
        bx = INT(BODY2D.X(boxBody))
        by = INT(BODY2D.Y(boxBody))
        fx = INT(BODY2D.X(floorBody))
        fy = INT(BODY2D.Y(floorBody))

        DRAW.RECTANGLE(fx - 400, fy - 25, 800, 50, 100, 100, 100, 255)
        DRAW.RECTANGLE(bx - 20,  by - 20,  40,  40, 200,  50,  50, 255)

        ; collision flash
        IF BODY2D.COLLIDED(boxBody)
            DRAW.TEXT("HIT!", bx, by - 30, 20, 255, 255, 0, 255)
        END IF
    CAMERA2D.END()
    RENDER.FRAME()
WEND

PHYSICS2D.STOP()
WINDOW.CLOSE()
```

---

## Extended Command Reference

| Command | Description |
|--------|-------------|
| `PHYSICS2D.SETSTEP(dt)` | Override fixed 2D physics timestep. |
| `PHYSICS2D.SETITERATIONS(n)` | Set velocity/position solver iterations. |
| `PHYSICS2D.ONCOLLISION(callback)` | Register 2D collision event handler. |
| `PHYSICS2D.PROCESSCOLLISIONS()` | Flush and dispatch pending 2D collision events. |
| `PHYSICS2D.DEBUGDRAW(bool)` | Enable/disable Box2D debug wireframe overlay. |
| `PHYSICS2D.GETDEBUGSEGMENTS()` | Returns array of debug line segments `[x0,y0,x1,y1,...]`. |

---

## See also

- [PHYSICS3D.md](PHYSICS3D.md) — Jolt 3D simulation
- [PHYSICS_ADVANCED.md](PHYSICS_ADVANCED.md) — joints, constraints
- [COLLISION.md](COLLISION.md) — stateless geometry tests
- [DRAW2D.md](DRAW2D.md) — 2D drawing
